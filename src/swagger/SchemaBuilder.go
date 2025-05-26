package swagger

import (
	openapi "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type SchemaBuilder struct {
	schema     *entity.Schema
	docBuilder *SwaggerDocBuilder
}

func (b *SchemaBuilder) Type(schemaType string) openapi.Schema {
	b.schema.Type = schemaType
	return b
}

func (b *SchemaBuilder) Format(format string) openapi.Schema {
	b.schema.Format = format
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

func (b *SchemaBuilder) Description(description string) openapi.Schema {
	b.schema.Description = description
	return b
}

func (b *SchemaBuilder) Title(title string) openapi.Schema {
	b.schema.Title = title
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

func (b *SchemaBuilder) MultipleOf(val float64) openapi.Schema {
	b.schema.MultipleOf = &val
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

func (b *SchemaBuilder) Items(config func(openapi.Schema)) openapi.Schema {
	itemsSchema := &entity.Schema{}
	itemSchemaBuilder := &SchemaBuilder{schema: itemsSchema, docBuilder: b.docBuilder}
	config(itemSchemaBuilder)
	b.schema.Items = &entity.SchemaRef{Schema: itemsSchema}
	return b
}

func (b *SchemaBuilder) ItemsRef(ref string) openapi.Schema {
	b.schema.Items = &entity.SchemaRef{Ref: ref}
	return b
}

func (b *SchemaBuilder) MaxProperties(max int) openapi.Schema {
	b.schema.MaxProperties = &max
	return b
}

func (b *SchemaBuilder) MinProperties(min int) openapi.Schema {
	b.schema.MinProperties = &min
	return b
}

func (b *SchemaBuilder) Required(fields ...string) openapi.Schema {
	b.schema.Required = append(b.schema.Required, fields...)
	return b
}

func (b *SchemaBuilder) Property(name string, config func(openapi.Schema)) openapi.Schema {
	if b.schema.Properties == nil {
		b.schema.Properties = make(map[string]*entity.SchemaRef)
	}
	propSchema := &entity.Schema{}
	propSchemaBuilder := &SchemaBuilder{schema: propSchema, docBuilder: b.docBuilder}
	config(propSchemaBuilder)
	b.schema.Properties[name] = &entity.SchemaRef{Schema: propSchema}
	return b
}

func (b *SchemaBuilder) PropertyRef(name string, ref string) openapi.Schema {
	if b.schema.Properties == nil {
		b.schema.Properties = make(map[string]*entity.SchemaRef)
	}
	b.schema.Properties[name] = &entity.SchemaRef{Ref: ref}
	return b
}

func (b *SchemaBuilder) AdditionalProperties(allowed bool, config ...func(openapi.Schema)) openapi.Schema {
	if !allowed {
		b.schema.AdditionalProperties = false
		return b
	}
	if len(config) > 0 && config[0] != nil {
		apSchema := &entity.Schema{}
		apBuilder := &SchemaBuilder{schema: apSchema, docBuilder: b.docBuilder}
		config[0](apBuilder)
		b.schema.AdditionalProperties = &entity.SchemaRef{Schema: apSchema}
	} else {
		b.schema.AdditionalProperties = true // Default to true if allowed and no schema provided
	}
	return b
}

func (b *SchemaBuilder) Nullable(nullable bool) openapi.Schema {
	b.schema.Nullable = nullable
	return b
}

func (b *SchemaBuilder) ReadOnly(readOnly bool) openapi.Schema {
	b.schema.ReadOnly = readOnly
	return b
}

func (b *SchemaBuilder) WriteOnly(writeOnly bool) openapi.Schema {
	b.schema.WriteOnly = writeOnly
	return b
}

func (b *SchemaBuilder) Deprecated(deprecated bool) openapi.Schema {
	b.schema.Deprecated = deprecated
	return b
}

func (b *SchemaBuilder) Example(example interface{}) openapi.Schema {
	b.schema.Example = example
	return b
}

func (b *SchemaBuilder) ExternalDocumentation(url string, description string) openapi.Schema {
	b.schema.ExternalDocs = &entity.ExternalDocumentation{URL: url, Description: description}
	return b
}

func (b *SchemaBuilder) XML(name string, config func(openapi.IXMLObjectBuilder)) openapi.Schema {
	xmlObj := entity.XMLObject{Name: name}
	// xmlBuilder := &XMLObjectBuilder{xml: &xmlObj} // Assuming XMLObjectBuilder exists
	// config(xmlBuilder)
	b.schema.XML = &xmlObj
	return b
}

func (b *SchemaBuilder) Discriminator(propertyName string, mapping map[string]string) openapi.Schema {
	b.schema.Discriminator = &entity.Discriminator{PropertyName: propertyName, Mapping: mapping}
	return b
}

func (b *SchemaBuilder) addSchemaRefs(refs []*entity.SchemaRef, configs []func(openapi.Schema), stringRefs []string) []*entity.SchemaRef {
	for _, cfg := range configs {
		s := &entity.Schema{}
		builder := &SchemaBuilder{schema: s, docBuilder: b.docBuilder}
		cfg(builder)
		refs = append(refs, &entity.SchemaRef{Schema: s})
	}
	for _, strRef := range stringRefs {
		refs = append(refs, &entity.SchemaRef{Ref: strRef})
	}
	return refs
}

func (b *SchemaBuilder) AllOf(configs ...func(openapi.Schema)) openapi.Schema {
	b.schema.AllOf = b.addSchemaRefs(b.schema.AllOf, configs, nil)
	return b
}

func (b *SchemaBuilder) AllOfRefs(refs ...string) openapi.Schema {
	b.schema.AllOf = b.addSchemaRefs(b.schema.AllOf, nil, refs)
	return b
}

func (b *SchemaBuilder) Ref(ref string) openapi.Schema {
	b.schema.Ref = ref
	return b
}

func (b *SchemaBuilder) Build() entity.Schema {
	return *b.schema
}

func (b *SchemaBuilder) OneOf(configs ...func(openapi.Schema)) openapi.Schema {
	b.schema.OneOf = b.addSchemaRefs(b.schema.OneOf, configs, nil)
	return b
}
func (b *SchemaBuilder) OneOfRefs(refs ...string) openapi.Schema {
	b.schema.OneOf = b.addSchemaRefs(b.schema.OneOf, nil, refs)
	return b
}
func (b *SchemaBuilder) AnyOf(configs ...func(openapi.Schema)) openapi.Schema {
	b.schema.AnyOf = b.addSchemaRefs(b.schema.AnyOf, configs, nil)
	return b
}
func (b *SchemaBuilder) AnyOfRefs(refs ...string) openapi.Schema {
	b.schema.AnyOf = b.addSchemaRefs(b.schema.AnyOf, nil, refs)
	return b
}
func (b *SchemaBuilder) Not(config func(openapi.Schema)) openapi.Schema {
	s := &entity.Schema{}
	builder := &SchemaBuilder{schema: s, docBuilder: b.docBuilder}
	config(builder)
	b.schema.Not = &entity.SchemaRef{Schema: s}
	return b
}
func (b *SchemaBuilder) NotRef(ref string) openapi.Schema {
	b.schema.Not = &entity.SchemaRef{Ref: ref}
	return b
}
