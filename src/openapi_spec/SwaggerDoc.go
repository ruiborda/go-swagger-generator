package openapi_spec

// SwaggerDocEntity is the root document object for the OpenAPI specification.
// Field names are updated to match OAS3.
type SwaggerDocEntity struct {
	Openapi      string                       `json:"openapi" yaml:"openapi"` // REQUIRED. E.g. "3.0.0"
	Info         Info                         `json:"info" yaml:"info"`       // REQUIRED
	Servers      []*Server                    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        map[string]*PathItem         `json:"paths" yaml:"paths"`     // REQUIRED
	Components   *Components                  `json:"components,omitempty" yaml:"components,omitempty"`
	Security     []SecurityRequirement        `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         []*Tag                       `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation       `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}
