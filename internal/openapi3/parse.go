package openapi3

import (
	"fmt"
	"github.com/go-dummy/dummy/internal/faker"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"

	"github.com/go-dummy/dummy/internal/apischema"
)

func Parse(path string) (apischema.API, error) {
	file, err := read(path)
	if err != nil {
		return apischema.API{}, err
	}

	var openapi OpenAPI

	if err := yaml.Unmarshal(file, &openapi); err != nil {
		return apischema.API{}, err
	}

	f := faker.NewFaker()

	b := &builder{openapi: openapi, faker: f}
	return b.Build()
}

func read(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return readFromURL(path)
	}

	return readFromFile(path)
}

func readFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func readFromFile(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type builder struct {
	openapi    OpenAPI
	operations []apischema.Operation
	faker      faker.Faker
}

func (b *builder) Build() (apischema.API, error) {
	for path, method := range b.openapi.Paths {
		if err := b.Add(path, http.MethodGet, method.Get); err != nil {
			return apischema.API{}, err
		}

		if err := b.Add(path, http.MethodPost, method.Post); err != nil {
			return apischema.API{}, err
		}

		if err := b.Add(path, http.MethodPut, method.Put); err != nil {
			return apischema.API{}, err
		}

		if err := b.Add(path, http.MethodPatch, method.Patch); err != nil {
			return apischema.API{}, err
		}

		if err := b.Add(path, http.MethodDelete, method.Delete); err != nil {
			return apischema.API{}, err
		}
	}

	return apischema.API{Operations: b.operations}, nil
}

func (b *builder) Add(path, method string, o *Operation) error {
	if o != nil {
		p := RemoveTrailingSlash(path)

		operation, err := b.Set(p, method, o)
		if err != nil {
			return err
		}

		b.operations = append(b.operations, operation)
	}

	return nil
}

func (b *builder) Set(path, method string, o *Operation) (apischema.Operation, error) {
	operation := apischema.Operation{
		Method: method,
		Path:   path,
	}

	for code, resp := range o.Responses {
		statusCode, err := strconv.Atoi(code)
		if err != nil {
			return apischema.Operation{}, err
		}

		content, found := resp.Content["application/json"]
		if !found {
			operation.Responses = append(operation.Responses, apischema.Response{
				StatusCode: statusCode,
			})
			continue
		}

		example := ExampleToResponse(content.Example)

		examples := make(map[string]any, len(content.Examples)+1)
		if len(content.Examples) > 0 {
			for key, example := range content.Examples {
				examples[key] = ExampleToResponse(example.Value)
			}
			examples[""] = ExampleToResponse(content.Examples[content.Examples.GetKeys()[0]].Value)
		}

		schema, err := b.convertSchema(content.Schema)
		if err != nil {
			return apischema.Operation{}, err
		}

		operation.Responses = append(operation.Responses, apischema.Response{
			StatusCode: statusCode,
			MediaType:  "application/json",
			Schema:     schema,
			Example:    example,
			Examples:   examples,
		})
	}

	return operation, nil
}

func (b *builder) convertSchema(s Schema) (apischema.Schema, error) {
	if s.Reference != "" {
		schema, err := b.openapi.LookupByReference(s.Reference)
		if err != nil {
			return nil, fmt.Errorf("resolve reference: %w", err)
		}

		s = schema
	}

	if s.Faker != "" {
		return apischema.FakerSchema{Example: b.faker.ByName(s.Faker)}, nil
	}

	switch s.Type {
	case "boolean":
		val, _ := s.Example.(bool)
		return apischema.BooleanSchema{Example: val}, nil
	case "integer":
		val, _ := s.Example.(int64)
		return apischema.IntSchema{Example: val}, nil
	case "number":
		val, _ := s.Example.(float64)
		return apischema.FloatSchema{Example: val}, nil
	case "string":
		val, _ := s.Example.(string)
		return apischema.StringSchema{Example: val}, nil
	case "array":
		itemsSchema, err := b.convertSchema(*s.Items)
		if err != nil {
			return nil, err
		}
		return apischema.ArraySchema{Type: itemsSchema, Example: parseArrayExample(s.Example)}, nil
	case "object":
		obj := apischema.ObjectSchema{Properties: make(map[string]apischema.Schema, len(s.Properties))}

		for key, prop := range s.Properties {
			propSchema, err := b.convertSchema(*prop)
			if err != nil {
				return nil, err
			}
			obj.Properties[key] = propSchema
		}

		obj.Example = parseObjectExample(s.Example)

		return obj, nil
	default:
		panic(fmt.Sprintf("unknown type %q", s.Type))
	}
}

func parseArrayExample(data any) []any {
	if data == nil {
		return []any{}
	}

	if data, ok := data.([]any); ok {
		res := make([]any, len(data))
		for k, v := range data {
			res[k] = parseExample(v.(map[any]any))
		}

		return res
	}

	panic(fmt.Sprintf("unpredicted type for example %T", data))
}

func parseObjectExample(data any) map[string]any {
	if data == nil {
		return map[string]any{}
	}

	if data, ok := data.(map[any]any); ok {
		return parseExample(data)
	}

	panic(fmt.Sprintf("unpredicted type for example %T", data))
}

// RemoveTrailingSlash returns path without trailing slash
func RemoveTrailingSlash(path string) string {
	if len(path) > 0 && path[len(path)-1] == '/' {
		return path[0 : len(path)-1]
	}

	return path
}
