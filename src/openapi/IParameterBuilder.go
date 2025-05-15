package openapi

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type Parameter interface {
	Description(description string) Parameter
	Required(required bool) Parameter
	Schema(s openapi_spec.SchemaEntity) Parameter
	SchemaFromDTO(dto interface{}) Parameter
	Type(paramType string) Parameter
	Format(format string) Parameter
	AllowEmptyValue(allow bool) Parameter
	Items(config func(Schema)) Parameter
	CollectionFormat(format string) Parameter
	Default(value interface{}) Parameter
	Maximum(max float64, exclusive bool) Parameter
	Minimum(min float64, exclusive bool) Parameter
	MaxLength(max int) Parameter
	MinLength(min int) Parameter
	Pattern(pattern string) Parameter
	MaxItems(max int) Parameter
	MinItems(min int) Parameter
	UniqueItems(unique bool) Parameter
	Enum(values ...interface{}) Parameter
	MultipleOf(val float64) Parameter
}
