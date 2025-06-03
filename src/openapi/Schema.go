package openapi

import (
	entity "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
)

type Schema interface {
	// Base JSON Schema type fields
	Type(schemaType string) Schema // string, number, integer, boolean, array, object
	Format(format string) Schema   // e.g., int32, int64, float, double, byte, binary, date, date-time, password
	Enum(values ...interface{}) Schema
	Default(value interface{}) Schema
	Description(description string) Schema
	Title(title string) Schema // Added for OAS3

	// Validation fields for numbers
	Maximum(max float64, exclusive bool) Schema
	Minimum(min float64, exclusive bool) Schema
	MultipleOf(val float64) Schema

	// Validation fields for strings
	MaxLength(max int) Schema
	MinLength(min int) Schema
	Pattern(pattern string) Schema

	// Validation fields for arrays
	MaxItems(max int) Schema
	MinItems(min int) Schema
	UniqueItems(unique bool) Schema
	Items(config func(Schema)) Schema // For array items schema
	ItemsRef(ref string) Schema       // For array items schema reference

	// Validation fields for objects
	MaxProperties(max int) Schema
	MinProperties(min int) Schema
	Required(fields ...string) Schema
	Property(name string, config func(Schema)) Schema                 // For object properties schema
	PropertyRef(name string, ref string) Schema                       // For object property schema reference
	AdditionalProperties(allowed bool, config ...func(Schema)) Schema // bool or schema

	// OpenAPI specific schema fields
	Nullable(nullable bool) Schema     // Added for OAS3
	ReadOnly(readOnly bool) Schema     // Added for OAS3
	WriteOnly(writeOnly bool) Schema   // Added for OAS3
	Deprecated(deprecated bool) Schema // Added for OAS3
	Example(example interface{}) Schema
	ExternalDocumentation(url string, description string) Schema         // Added for OAS3
	XML(name string, config func(IXMLObjectBuilder)) Schema              // Added for OAS3
	Discriminator(propertyName string, mapping map[string]string) Schema // Added for OAS3

	// Composition fields
	AllOf(configs ...func(Schema)) Schema
	AllOfRefs(refs ...string) Schema
	OneOf(configs ...func(Schema)) Schema
	OneOfRefs(refs ...string) Schema
	AnyOf(configs ...func(Schema)) Schema
	AnyOfRefs(refs ...string) Schema
	Not(config func(Schema)) Schema
	NotRef(ref string) Schema

	Ref(ref string) Schema // To make this schema a reference to another schema
	Build() entity.Schema  // To get the built schema object
}

// IXMLObjectBuilder defines the builder for an XMLObject
type IXMLObjectBuilder interface {
	Name(name string) IXMLObjectBuilder
	Namespace(namespace string) IXMLObjectBuilder
	Prefix(prefix string) IXMLObjectBuilder
	Attribute(attribute bool) IXMLObjectBuilder
	Wrapped(wrapped bool) IXMLObjectBuilder
}
