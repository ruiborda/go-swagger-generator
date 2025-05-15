package openapi

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type Schema interface {
	Type(schemaType string) Schema
	Format(format string) Schema
	Ref(ref string) Schema
	Description(description string) Schema
	Enum(values ...interface{}) Schema
	Default(value interface{}) Schema
	Items(config func(Schema)) Schema
	Properties(props map[string]*openapi_spec.SchemaEntity) Schema
	Property(name string, config func(Schema)) Schema
	Required(fields ...string) Schema
	Maximum(max float64, exclusive bool) Schema
	Minimum(min float64, exclusive bool) Schema
	MaxLength(max int) Schema
	MinLength(min int) Schema
	Pattern(pattern string) Schema
	MaxItems(max int) Schema
	MinItems(min int) Schema
	UniqueItems(unique bool) Schema
	Example(example interface{}) Schema
	// Consider adding MultipleOf here if needed for direct schema building
	// MultipleOf(val float64) Schema
}
