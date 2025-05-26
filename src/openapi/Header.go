package openapi

type Header interface {
	Description(description string) Header
	Required(required bool) Header
	Deprecated(deprecated bool) Header
	AllowEmptyValue(allow bool) Header

	Schema(config func(Schema)) Header
	SchemaFromDTO(dto interface{}) Header
	SchemaRef(ref string) Header

	Example(value interface{}) Header
	Examples(name string, config func(Example)) Header

	// Content(mimeType string, config func(MediaType)) Header // For complex headers
	Ref(ref string) Header // For referencing a header in components
}
