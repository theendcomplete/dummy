package openapi3

import (
	"fmt"
	"strings"
)

// OpenAPI Object
// See specification https://swagger.io/specification/#openapi-object
type OpenAPI struct {
	Info       Info       `json:"info" yaml:"info"`
	Paths      Paths      `json:"paths" yaml:"paths"`
	Components Components `json:"components,omitempty" yaml:"components,omitempty"`
}

func (api OpenAPI) LookupByReference(ref string) (Schema, error) {
	schema := api.Components.Schemas[schemaKey(ref)]
	if schema == nil {
		return Schema{}, fmt.Errorf("unknown schema %q", schema.Reference)
	}

	return *schema, nil
}

func schemaKey(ref string) string {
	const prefix = "#/components/schemas/"
	return strings.TrimPrefix(ref, prefix)
}
