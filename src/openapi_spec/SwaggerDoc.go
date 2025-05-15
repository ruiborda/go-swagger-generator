package openapi_spec

type SwaggerDocEntity struct {
	Swagger             string                          `json:"swagger"`
	Info                InfoEntity                      `json:"info"`
	Host                string                          `json:"host,omitempty"`
	BasePath            string                          `json:"basePath,omitempty"`
	Tags                []TagEntity                     `json:"tags,omitempty"`
	Schemes             []string                        `json:"schemes,omitempty"`
	Paths               map[string]PathItemEntity       `json:"paths"`
	SecurityDefinitions map[string]SecuritySchemeEntity `json:"securityDefinitions,omitempty"`
	Definitions         map[string]SchemaEntity         `json:"definitions,omitempty"`
	ExternalDocs        *ExternalDocumentationEntity    `json:"externalDocs,omitempty"`
}
