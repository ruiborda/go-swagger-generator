package openapi_spec

// ResponseEntity describes a single response from an API Operation, including design-time, static links to operations based on the response.
type Response struct {
	Description string                `json:"description" yaml:"description"` // REQUIRED
	Headers     map[string]*HeaderRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]*LinkRef   `json:"links,omitempty" yaml:"links,omitempty"`
}
