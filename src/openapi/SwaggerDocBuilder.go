package openapi

import (
	entity "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
)

type SwaggerDocBuilder interface {
	OpenAPIVersion(version string) SwaggerDocBuilder // Renamed from SwaggerVersion
	Info(config func(info Info)) SwaggerDocBuilder
	Server(url string, config func(Server)) SwaggerDocBuilder
	Servers(servers ...entity.Server) SwaggerDocBuilder
	Tag(name string, config func(tag Tag)) SwaggerDocBuilder
	Path(pathPattern string) PathItem

	// Component related methods
	ComponentSchema(name string, schema entity.Schema) SwaggerDocBuilder
	ComponentSchemaRef(name string, ref string) SwaggerDocBuilder
	ComponentResponse(name string, config func(Response)) SwaggerDocBuilder
	ComponentParameter(name string, config func(Parameter)) SwaggerDocBuilder
	ComponentExample(name string, config func(Example)) SwaggerDocBuilder
	ComponentRequestBody(name string, config func(RequestBody)) SwaggerDocBuilder
	ComponentHeader(name string, config func(Header)) SwaggerDocBuilder
	ComponentSecurityScheme(name string, config func(SecurityScheme)) SwaggerDocBuilder
	ComponentLink(name string, config func(Link)) SwaggerDocBuilder
	ComponentCallback(name string, config func(Callback)) SwaggerDocBuilder

	SchemaFromDTO(dto interface{}) (string, error) // Renamed from DefinitionFromDTO
	ExternalDocumentation(url string, description string) SwaggerDocBuilder
	GlobalSecurity(requirement entity.SecurityRequirement) SwaggerDocBuilder
	Build() entity.SwaggerDocEntity
}
