package swagger

import (
	"fmt"
	openapi2 "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity2 "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type ResponseBuilder struct {
	response   *entity2.ResponseEntity
	docBuilder *SwaggerDocBuilder
}

func (b *ResponseBuilder) Description(description string) openapi2.Response {
	b.response.Description = description
	return b
}
func (b *ResponseBuilder) Schema(s entity2.SchemaEntity) openapi2.Response {
	b.response.Schema = &s
	return b
}
func (b *ResponseBuilder) SchemaFromDTO(dto interface{}) openapi2.Response {
	dtoName, err := b.docBuilder.DefinitionFromDTO(dto)
	if err != nil {
		fmt.Printf("Error adding DTO definition for response schema: %v\n", err)
		return b
	}
	b.response.Schema = &entity2.SchemaEntity{Ref: "#/definitions/" + dtoName}
	return b
}
func (b *ResponseBuilder) SchemaRef(ref string) openapi2.Response {
	b.response.Schema = &entity2.SchemaEntity{Ref: ref}
	return b
}
func (b *ResponseBuilder) Header(name string, config func(builder openapi2.Header)) openapi2.Response {
	if b.response.Headers == nil {
		b.response.Headers = make(map[string]entity2.HeaderEntity)
	}
	header := entity2.HeaderEntity{}
	// Pass docBuilder if HeaderBuilder needs it (e.g., for SchemaFromDTO in items)
	headerCfg := &HeaderBuilder{header: &header /* docBuilder: b.docBuilder */}
	config(headerCfg)
	b.response.Headers[name] = header
	return b
}
func (b *ResponseBuilder) Example(mimeType string, exampleValue interface{}) openapi2.Response {
	if b.response.Examples == nil {
		b.response.Examples = make(map[string]interface{})
	}
	b.response.Examples[mimeType] = exampleValue
	return b
}
