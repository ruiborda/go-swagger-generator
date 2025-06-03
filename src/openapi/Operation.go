package openapi

import (
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
)

type Operation interface {
	Summary(summary string) Operation
	Description(description string) Operation
	OperationID(id string) Operation
	Tag(tag string) Operation
	Tags(tags ...string) Operation

	Parameter(name, in string, config func(Parameter)) Operation
	QueryParameter(name string, config func(Parameter)) Operation
	PathParameter(name string, config func(Parameter)) Operation
	HeaderParameter(name string, config func(Parameter)) Operation
	CookieParameter(name string, config func(Parameter)) Operation

	RequestBody(config func(RequestBody)) Operation

	Response(statusCode int, config func(Response)) Operation
	DefaultResponse(config func(Response)) Operation

	Security(schemeName string, scopes ...string) Operation
	Deprecated(deprecated bool) Operation
	ExternalDocumentation(url string, description string) Operation
	Server(url string, config func(Server)) Operation

	Path() PathItem
}

type RequestBody interface {
	Description(description string) RequestBody
	Content(mimeType mime.MimeType, config func(MediaType)) RequestBody
	Required(required bool) RequestBody
	Ref(ref string) RequestBody // For referencing a request body in components
}

type MediaType interface {
	Schema(config func(Schema)) MediaType
	SchemaFromDTO(dto interface{}) MediaType
	SchemaRef(ref string) MediaType
	Example(value interface{}) MediaType
	Examples(name string, config func(Example)) MediaType
	Encoding(propertyName string, config func(Encoding)) MediaType
}

type Encoding interface {
	ContentType(contentType string) Encoding
	Header(name string, config func(Header)) Encoding
	Style(style string) Encoding
	Explode(explode bool) Encoding
	AllowReserved(allowReserved bool) Encoding
}
