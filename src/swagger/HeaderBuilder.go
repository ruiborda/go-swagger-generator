package swagger

import (
	openapi2 "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity2 "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type HeaderBuilder struct {
	header *entity2.HeaderEntity
	// docBuilder *SwaggerDocBuilder // Uncomment if needed for SchemaFromDTO in Items
}

func (b *HeaderBuilder) Description(description string) openapi2.Header {
	b.header.Description = description
	return b
}
func (b *HeaderBuilder) Type(headerType string) openapi2.Header {
	b.header.Type = headerType
	return b
}
func (b *HeaderBuilder) Format(format string) openapi2.Header {
	b.header.Format = format
	return b
}
func (b *HeaderBuilder) Items(config func(openapi2.Schema)) openapi2.Header {
	itemsSchema := &entity2.SchemaEntity{}
	schemaBuilder := &SchemaBuilder{schema: itemsSchema /* docBuilder: b.docBuilder */}
	config(schemaBuilder)
	b.header.Items = itemsSchema
	return b
}
func (b *HeaderBuilder) CollectionFormat(format string) openapi2.Header {
	b.header.CollectionFormat = format
	return b
}
func (b *HeaderBuilder) Default(value interface{}) openapi2.Header {
	b.header.Default = value
	return b
}
func (b *HeaderBuilder) Maximum(max float64, exclusive bool) openapi2.Header {
	b.header.Maximum = &max
	b.header.ExclusiveMaximum = exclusive
	return b
}
func (b *HeaderBuilder) Minimum(min float64, exclusive bool) openapi2.Header {
	b.header.Minimum = &min
	b.header.ExclusiveMinimum = exclusive
	return b
}
func (b *HeaderBuilder) MaxLength(max int) openapi2.Header {
	b.header.MaxLength = &max
	return b
}
func (b *HeaderBuilder) MinLength(min int) openapi2.Header {
	b.header.MinLength = &min
	return b
}
func (b *HeaderBuilder) Pattern(pattern string) openapi2.Header {
	b.header.Pattern = pattern
	return b
}
func (b *HeaderBuilder) MaxItems(max int) openapi2.Header {
	b.header.MaxItems = &max
	return b
}
func (b *HeaderBuilder) MinItems(min int) openapi2.Header {
	b.header.MinItems = &min
	return b
}
func (b *HeaderBuilder) UniqueItems(unique bool) openapi2.Header {
	b.header.UniqueItems = unique
	return b
}
func (b *HeaderBuilder) Enum(values ...interface{}) openapi2.Header {
	b.header.Enum = values
	return b
}
func (b *HeaderBuilder) MultipleOf(val float64) openapi2.Header {
	b.header.MultipleOf = &val
	return b
}
