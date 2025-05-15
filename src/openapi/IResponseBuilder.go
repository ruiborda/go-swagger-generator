package openapi

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type Response interface {
	Description(description string) Response
	Schema(s openapi_spec.SchemaEntity) Response
	SchemaFromDTO(dto interface{}) Response
	SchemaRef(ref string) Response
	Header(name string, config func(Header)) Response
	Example(mimeType string, exampleValue interface{}) Response
}
