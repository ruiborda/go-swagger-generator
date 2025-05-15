package swagger

import (
	"fmt"
	openapi2 "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity2 "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type ParameterBuilder struct {
	param      *entity2.ParameterEntity
	docBuilder *SwaggerDocBuilder
}

func (b *ParameterBuilder) Description(description string) openapi2.Parameter {
	b.param.Description = description
	return b
}
func (b *ParameterBuilder) Required(required bool) openapi2.Parameter {
	b.param.Required = required
	return b
}
func (b *ParameterBuilder) Schema(s entity2.SchemaEntity) openapi2.Parameter {
	b.param.Schema = &s
	return b
}
func (b *ParameterBuilder) SchemaFromDTO(dto interface{}) openapi2.Parameter {
	dtoName, err := b.docBuilder.DefinitionFromDTO(dto)
	if err != nil {
		// Consider logging or returning error in a real application
		fmt.Printf("Error adding DTO definition for parameter schema: %v\n", err)
		return b
	}
	b.param.Schema = &entity2.SchemaEntity{Ref: "#/definitions/" + dtoName}
	return b
}
func (b *ParameterBuilder) Type(paramType string) openapi2.Parameter {
	b.param.Type = paramType
	return b
}
func (b *ParameterBuilder) Format(format string) openapi2.Parameter {
	b.param.Format = format
	return b
}
func (b *ParameterBuilder) AllowEmptyValue(allow bool) openapi2.Parameter {
	b.param.AllowEmptyValue = allow
	return b
}
func (b *ParameterBuilder) Items(config func(builder openapi2.Schema)) openapi2.Parameter {
	itemsSchema := &entity2.SchemaEntity{}
	schemaBuilder := &SchemaBuilder{schema: itemsSchema, docBuilder: b.docBuilder}
	config(schemaBuilder)
	b.param.Items = itemsSchema
	return b
}
func (b *ParameterBuilder) CollectionFormat(format string) openapi2.Parameter {
	b.param.CollectionFormat = format
	return b
}
func (b *ParameterBuilder) Default(value interface{}) openapi2.Parameter {
	b.param.Default = value
	return b
}
func (b *ParameterBuilder) Maximum(max float64, exclusive bool) openapi2.Parameter {
	b.param.Maximum = &max
	b.param.ExclusiveMaximum = exclusive
	return b
}
func (b *ParameterBuilder) Minimum(min float64, exclusive bool) openapi2.Parameter {
	b.param.Minimum = &min
	b.param.ExclusiveMinimum = exclusive
	return b
}
func (b *ParameterBuilder) MaxLength(max int) openapi2.Parameter {
	b.param.MaxLength = &max
	return b
}
func (b *ParameterBuilder) MinLength(min int) openapi2.Parameter {
	b.param.MinLength = &min
	return b
}
func (b *ParameterBuilder) Pattern(pattern string) openapi2.Parameter {
	b.param.Pattern = pattern
	return b
}
func (b *ParameterBuilder) MaxItems(max int) openapi2.Parameter {
	b.param.MaxItems = &max
	return b
}
func (b *ParameterBuilder) MinItems(min int) openapi2.Parameter {
	b.param.MinItems = &min
	return b
}
func (b *ParameterBuilder) UniqueItems(unique bool) openapi2.Parameter {
	b.param.UniqueItems = unique
	return b
}
func (b *ParameterBuilder) Enum(values ...interface{}) openapi2.Parameter {
	b.param.Enum = values
	return b
}
func (b *ParameterBuilder) MultipleOf(val float64) openapi2.Parameter {
	b.param.MultipleOf = &val
	return b
}
