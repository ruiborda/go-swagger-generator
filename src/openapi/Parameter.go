package openapi

type Parameter interface {
	Description(description string) Parameter
	Required(required bool) Parameter
	Deprecated(deprecated bool) Parameter // Added for OAS3
	AllowEmptyValue(allow bool) Parameter

	Schema(config func(Schema)) Parameter // Schema is primary way to define type/format
	SchemaFromDTO(dto interface{}) Parameter
	SchemaRef(ref string) Parameter

	Style(style string) Parameter   // Added for OAS3 serialization
	Explode(explode bool) Parameter // Added for OAS3 serialization

	Example(value interface{}) Parameter
	Examples(name string, config func(Example)) Parameter

	// Content(mimeType string, config func(MediaType)) Parameter // For complex parameters, less common
	Ref(ref string) Parameter // For referencing a parameter in components
}
