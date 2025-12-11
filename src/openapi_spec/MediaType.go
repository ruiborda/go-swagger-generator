package openapi_spec

// MediaType provides schema and examples for the media type identified by its key.
type MediaType struct {
	Schema   *SchemaRef             `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*ExampleRef `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding   `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// Encoding is a declaration of how specific properties are encoded in a multipart request body.
type Encoding struct {
	ContentType   string                `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       map[string]*HeaderRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string                `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool                  `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool                  `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}
