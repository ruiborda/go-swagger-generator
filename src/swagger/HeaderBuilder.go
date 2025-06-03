package swagger

import (
	openapi "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
)

type HeaderBuilder struct {
	header     *entity.Header
	docBuilder *SwaggerDocBuilder
}

func (b *HeaderBuilder) Description(description string) openapi.Header {
	b.header.Description = description
	return b
}

func (b *HeaderBuilder) Required(required bool) openapi.Header {
	b.header.Required = required
	return b
}

func (b *HeaderBuilder) Deprecated(deprecated bool) openapi.Header {
	b.header.Deprecated = deprecated
	return b
}

func (b *HeaderBuilder) AllowEmptyValue(allow bool) openapi.Header {
	b.header.AllowEmptyValue = allow
	return b
}

func (b *HeaderBuilder) initSchema() {
	if b.header.Schema == nil {
		b.header.Schema = &entity.Schema{}
	}
}

func (b *HeaderBuilder) Schema(config func(openapi.Schema)) openapi.Header {
	b.initSchema()
	schemaBuilder := &SchemaBuilder{schema: b.header.Schema, docBuilder: b.docBuilder}
	config(schemaBuilder)
	return b
}

func (b *HeaderBuilder) SchemaFromDTO(dto interface{}) openapi.Header {
	dtoName, err := b.docBuilder.SchemaFromDTO(dto)
	if err != nil {
		// log error
		return b
	}
	b.initSchema()
	b.header.Schema.Ref = "#/components/schemas/" + dtoName
	return b
}

func (b *HeaderBuilder) SchemaRef(ref string) openapi.Header {
	b.initSchema()
	b.header.Schema.Ref = ref
	return b
}

func (b *HeaderBuilder) Example(value interface{}) openapi.Header {
	b.header.Example = value
	return b
}

func (b *HeaderBuilder) Examples(name string, config func(openapi.Example)) openapi.Header {
	if b.header.Examples == nil {
		b.header.Examples = make(map[string]*entity.ExampleRef)
	}
	ex := entity.Example{}
	exBuilder := &ExampleBuilder{example: &ex}
	config(exBuilder)
	b.header.Examples[name] = &entity.ExampleRef{Example: &ex}
	return b
}

func (b *HeaderBuilder) Ref(ref string) openapi.Header {
	// Handled by higher level component (e.g. ResponseBuilder Header method with ref)
	return b
}
