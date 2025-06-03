package openapi

import "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"

type Response interface {
	Description(description string) Response
	Header(name string, config func(Header)) Response
	Content(mimeType mime.MimeType, config func(MediaType)) Response // Replaces Schema/SchemaFromDTO/SchemaRef
	Link(name string, config func(Link)) Response
	Ref(ref string) Response // For referencing a response in components
}

// Link defines the builder for a Link object
type Link interface {
	OperationRef(ref string) Link
	OperationID(id string) Link
	Parameter(name string, expressionOrValue interface{}) Link
	RequestBody(expressionOrValue interface{}) Link
	Description(description string) Link
	Server(url string, config func(Server)) Link
	Ref(ref string) Link // For referencing a link in components
}

// Example defines the builder for an Example object
type Example interface {
	Summary(summary string) Example
	Description(description string) Example
	Value(value interface{}) Example
	ExternalValue(url string) Example
	Ref(ref string) Example // For referencing an example in components
}
