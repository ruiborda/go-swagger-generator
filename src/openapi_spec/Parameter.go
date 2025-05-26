package openapi_spec

// Parameter defines a parameter in OAS3.
// Type, Format, Items, etc., are now part of the Schema object.
// CollectionFormat is replaced by 'style' and 'explode'.
type Parameter struct {
	Name            string      `json:"name" yaml:"name"`
	In              string      `json:"in" yaml:"in"` // "query", "header", "path", "cookie"
	Description     string      `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool        `json:"required,omitempty" yaml:"required,omitempty"` // Note: Must be true if 'in' is "path".
	Deprecated      bool        `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`

	Style   string `json:"style,omitempty" yaml:"style,omitempty"`     // How the parameter is serialized.
	Explode bool   `json:"explode,omitempty" yaml:"explode,omitempty"` // Used with 'style'.

	Schema   *Schema                `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*ExampleRef `json:"examples,omitempty" yaml:"examples,omitempty"`

	// Content is an alternative to schema for more complex serialization scenarios.
	// Content map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}
