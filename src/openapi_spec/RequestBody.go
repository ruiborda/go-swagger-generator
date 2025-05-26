package openapi_spec

// RequestBody describes a single request body.
type RequestBody struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]*MediaType `json:"content" yaml:"content"` // REQUIRED. Key is media type.
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty"`
}
