package openapi

import (
	entity2 "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type SwaggerDoc interface {
	SwaggerVersion(version string) SwaggerDoc
	Info(config func(Info)) SwaggerDoc
	Host(host string) SwaggerDoc
	BasePath(basePath string) SwaggerDoc
	Tag(name string, config func(Tag)) SwaggerDoc
	Scheme(scheme string) SwaggerDoc
	Schemes(schemes ...string) SwaggerDoc
	Path(pathPattern string) PathItem
	SecurityDefinition(name string, config func(SecurityScheme)) SwaggerDoc
	Definition(name string, schema entity2.SchemaEntity) SwaggerDoc
	DefinitionFromDTO(dto interface{}) (string, error)
	ExternalDocumentation(url string, description string) SwaggerDoc
	Build() entity2.SwaggerDocEntity
}
