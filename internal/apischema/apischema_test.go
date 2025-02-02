package apischema_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-dummy/dummy/internal/apischema"
)

func TestSchema_ExampleValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		schema apischema.Schema
		want   any
	}{
		{
			name:   "boolean: default",
			schema: apischema.BooleanSchema{},
			want:   false,
		},
		{
			name:   "boolean: with example",
			schema: apischema.BooleanSchema{Example: true},
			want:   true,
		},
		{
			name:   "int: default",
			schema: apischema.IntSchema{},
			want:   int64(0),
		},
		{
			name:   "int: with example",
			schema: apischema.IntSchema{Example: 42},
			want:   int64(42),
		},
		{
			name:   "float: default",
			schema: apischema.FloatSchema{},
			want:   0.0,
		},
		{
			name:   "float: with example",
			schema: apischema.FloatSchema{Example: 4.2},
			want:   4.2,
		},
		{
			name:   "string: default",
			schema: apischema.StringSchema{},
			want:   "",
		},
		{
			name:   "string: with example",
			schema: apischema.StringSchema{Example: "John"},
			want:   "John",
		},
		{
			name:   "array: default",
			schema: apischema.ArraySchema{Type: apischema.StringSchema{}},
			want:   []any{""},
		},
		{
			name:   "array: with int example",
			schema: apischema.ArraySchema{Example: []any{4, 2}},
			want:   []any{4, 2},
		},
		{
			name:   "array: with string example",
			schema: apischema.ArraySchema{Example: []any{"4", "2"}},
			want:   []any{"4", "2"},
		},
		{
			name:   "object: default",
			schema: apischema.ObjectSchema{},
			want:   map[string]any{},
		},
		{
			name:   "object: with example",
			schema: apischema.ObjectSchema{Example: map[string]any{"a": "4", "b": "2"}},
			want:   map[string]any{"a": "4", "b": "2"},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.schema.ExampleValue()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestResponse_ExampleValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		response apischema.Response
		key      string
		want     any
	}{
		{
			name:     "boolean examples",
			response: apischema.Response{Schema: apischema.BooleanSchema{}, Examples: map[string]any{"first": false, "second": true}},
			key:      "first",
			want:     false,
		},
		{
			name:     "int examples",
			response: apischema.Response{Schema: apischema.IntSchema{}, Examples: map[string]any{"first": int64(4), "second": int64(2)}},
			key:      "first",
			want:     int64(4),
		},
		{
			name:     "float examples",
			response: apischema.Response{Schema: apischema.FloatSchema{}, Examples: map[string]any{"first": 4.0, "second": 2.0}},
			key:      "first",
			want:     4.0,
		},
		{
			name:     "string examples",
			response: apischema.Response{Schema: apischema.StringSchema{}, Examples: map[string]any{"first": "abc", "second": "xyz"}},
			key:      "first",
			want:     "abc",
		},
		{
			name:     "array examples",
			response: apischema.Response{Schema: apischema.ArraySchema{}, Examples: map[string]any{"first": []int64{4, 2}, "second": []int64{0, 0}}},
			key:      "first",
			want:     []int64{4, 2},
		},
		{
			name:     "object examples",
			response: apischema.Response{Schema: apischema.ObjectSchema{}, Examples: map[string]any{"first": map[string]any{"first": "abc"}, "second": map[string]any{"first": "xyz"}}},
			key:      "first",
			want:     map[string]any{"first": "abc"},
		},
		{
			name: "use schema example",
			response: apischema.Response{Examples: map[string]any{"first": map[string]any{"first": "abc"}, "second": map[string]any{"first": "xyz"}},
				Schema: apischema.ObjectSchema{
					Example: map[string]any{"first": "schema variant"},
				}},
			key:  "third",
			want: map[string]any{"first": "schema variant"},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.response.ExampleValue(tc.key)

			require.Equal(t, tc.want, got)
		})
	}
}
