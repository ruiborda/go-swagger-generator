package swagger

import (
	"fmt"
	openapi "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
)

type ParameterBuilder struct {
	param      *entity.Parameter
	docBuilder *SwaggerDocBuilder
}

func (b *ParameterBuilder) Description(description string) openapi.Parameter {
	b.param.Description = description
	return b
}

func (b *ParameterBuilder) Required(required bool) openapi.Parameter {
	b.param.Required = required
	return b
}

func (b *ParameterBuilder) Deprecated(deprecated bool) openapi.Parameter {
	b.param.Deprecated = deprecated
	return b
}

func (b *ParameterBuilder) AllowEmptyValue(allow bool) openapi.Parameter {
	b.param.AllowEmptyValue = allow
	return b
}

func (b *ParameterBuilder) initSchema() {
	if b.param.Schema == nil {
		b.param.Schema = &entity.Schema{}
	}
}

func (b *ParameterBuilder) Schema(config func(openapi.Schema)) openapi.Parameter {
	b.initSchema()
	schemaBuilder := &SchemaBuilder{schema: b.param.Schema, docBuilder: b.docBuilder}
	config(schemaBuilder)
	return b
}

func (b *ParameterBuilder) SchemaFromDTO(dto interface{}) openapi.Parameter {
	dtoName, err := b.docBuilder.SchemaFromDTO(dto)
	if err != nil {
		fmt.Printf("Error adding DTO definition for parameter schema: %v\n", err)
		return b
	}
	b.initSchema()
	b.param.Schema.Ref = "#/components/schemas/" + dtoName
	return b
}

func (b *ParameterBuilder) SchemaRef(ref string) openapi.Parameter {
	b.initSchema()
	b.param.Schema.Ref = ref
	return b
}

func (b *ParameterBuilder) Style(style string) openapi.Parameter {
	b.param.Style = style
	return b
}

func (b *ParameterBuilder) Explode(explode bool) openapi.Parameter {
	b.param.Explode = explode
	return b
}

func (b *ParameterBuilder) Example(value interface{}) openapi.Parameter {
	b.param.Example = value
	return b
}

func (b *ParameterBuilder) Examples(name string, config func(openapi.Example)) openapi.Parameter {
	if b.param.Examples == nil {
		b.param.Examples = make(map[string]*entity.ExampleRef)
	}
	ex := entity.Example{}
	exBuilder := &ExampleBuilder{example: &ex}
	config(exBuilder)
	b.param.Examples[name] = &entity.ExampleRef{Example: &ex}
	return b
}

func (b *ParameterBuilder) Ref(ref string) openapi.Parameter {
	return b
}
