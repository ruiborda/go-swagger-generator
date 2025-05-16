package openapi

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
)

type Operation interface {
	Summary(summary string) Operation
	Description(description string) Operation
	OperationID(id string) Operation
	Tag(tag string) Operation
	Tags(tags ...string) Operation
	Consume(mimeType string) Operation
	Consumes(mimeTypes ...string) Operation
	Produce(mimeType mime.MimeType) Operation
	Produces(mimeTypes ...mime.MimeType) Operation
	Parameter(name, in string, config func(Parameter)) Operation
	QueryParameter(name string, config func(Parameter)) Operation
	PathParameter(name string, config func(Parameter)) Operation
	HeaderParameter(name string, config func(Parameter)) Operation
	FormParameter(name string, config func(Parameter)) Operation
	BodyParameter(config func(Parameter)) Operation
	Response(statusCode int, config func(Response)) Operation
	Security(schemeName string, scopes ...string) Operation
	Deprecated(deprecated bool) Operation
	ExternalDocumentation(url string, description string) Operation
	Path() PathItem
}
