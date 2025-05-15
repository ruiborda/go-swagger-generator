package swagger

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type SchemaBuilder struct {
	schema     *openapi_spec.SchemaEntity
	docBuilder *SwaggerDocBuilder // For SchemaFromDTO if used within items/properties
}

func (b *SchemaBuilder) Type(schemaType string) openapi.Schema {
	b.schema.Type = schemaType
	return b
}
func (b *SchemaBuilder) Format(format string) openapi.Schema {
	b.schema.Format = format
	return b
}
func (b *SchemaBuilder) Ref(ref string) openapi.Schema {
	b.schema.Ref = ref
	return b
}
func (b *SchemaBuilder) Description(description string) openapi.Schema {
	b.schema.Description = description
	return b
}
func (b *SchemaBuilder) Enum(values ...interface{}) openapi.Schema {
	b.schema.Enum = values
	return b
}
func (b *SchemaBuilder) Default(value interface{}) openapi.Schema {
	b.schema.Default = value
	return b
}
func (b *SchemaBuilder) Items(config func(openapi.Schema)) openapi.Schema {
	itemsSchema := &openapi_spec.SchemaEntity{}
	itemSchemaBuilder := &SchemaBuilder{schema: itemsSchema, docBuilder: b.docBuilder}
	config(itemSchemaBuilder)
	b.schema.Items = itemsSchema
	return b
}
func (b *SchemaBuilder) Properties(props map[string]*openapi_spec.SchemaEntity) openapi.Schema {
	b.schema.Properties = props
	return b
}
func (b *SchemaBuilder) Property(name string, config func(openapi.Schema)) openapi.Schema {
	if b.schema.Properties == nil {
		b.schema.Properties = make(map[string]*openapi_spec.SchemaEntity)
	}
	propSchema := &openapi_spec.SchemaEntity{}
	propSchemaBuilder := &SchemaBuilder{schema: propSchema, docBuilder: b.docBuilder}
	config(propSchemaBuilder)
	b.schema.Properties[name] = propSchema
	return b
}
func (b *SchemaBuilder) Required(fields ...string) openapi.Schema {
	b.schema.Required = append(b.schema.Required, fields...)
	return b
}
func (b *SchemaBuilder) Maximum(max float64, exclusive bool) openapi.Schema {
	b.schema.Maximum = &max
	b.schema.ExclusiveMaximum = exclusive
	return b
}
func (b *SchemaBuilder) Minimum(min float64, exclusive bool) openapi.Schema {
	b.schema.Minimum = &min
	b.schema.ExclusiveMinimum = exclusive
	return b
}
func (b *SchemaBuilder) MaxLength(max int) openapi.Schema {
	b.schema.MaxLength = &max
	return b
}
func (b *SchemaBuilder) MinLength(min int) openapi.Schema {
	b.schema.MinLength = &min
	return b
}
func (b *SchemaBuilder) Pattern(pattern string) openapi.Schema {
	b.schema.Pattern = pattern
	return b
}
func (b *SchemaBuilder) MaxItems(max int) openapi.Schema {
	b.schema.MaxItems = &max
	return b
}
func (b *SchemaBuilder) MinItems(min int) openapi.Schema {
	b.schema.MinItems = &min
	return b
}
func (b *SchemaBuilder) UniqueItems(unique bool) openapi.Schema {
	b.schema.UniqueItems = unique
	return b
}
func (b *SchemaBuilder) Example(example interface{}) openapi.Schema {
	b.schema.Example = example
	return b
}
