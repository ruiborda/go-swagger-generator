package openapi_spec

// Header follows the structure of a Parameter but is used in different contexts.
// For OAS3, a Header Object is similar to a Parameter Object without 'name' and 'in'.
type Header struct {
	Description     string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                   `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                   `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Schema          *Schema                `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]*ExampleRef `json:"examples,omitempty" yaml:"examples,omitempty"`
	// Content map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"` // For complex headers
}
