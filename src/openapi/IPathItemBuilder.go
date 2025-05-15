package openapi

type PathItem interface {
	Get(config func(Operation)) PathItem
	Post(config func(Operation)) PathItem
	Put(config func(Operation)) PathItem
	Delete(config func(Operation)) PathItem
	Options(config func(Operation)) PathItem
	Head(config func(Operation)) PathItem
	Patch(config func(Operation)) PathItem
	Parameter(name, in string, config func(Parameter)) PathItem
	Doc() SwaggerDoc
}
