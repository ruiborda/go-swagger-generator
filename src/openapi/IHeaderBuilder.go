package openapi

type Header interface {
	Description(description string) Header
	Type(headerType string) Header
	Format(format string) Header
	Items(config func(Schema)) Header
	CollectionFormat(format string) Header
	Default(value interface{}) Header
	Maximum(max float64, exclusive bool) Header
	Minimum(min float64, exclusive bool) Header
	MaxLength(max int) Header
	MinLength(min int) Header
	Pattern(pattern string) Header
	MaxItems(max int) Header
	MinItems(min int) Header
	UniqueItems(unique bool) Header
	Enum(values ...interface{}) Header
	MultipleOf(val float64) Header
}
