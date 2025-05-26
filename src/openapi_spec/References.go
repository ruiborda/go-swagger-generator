package openapi_spec

// This file contains Ref structs for all major referencable components in OpenAPI 3.0.
// These allow a component to be defined inline or via a $ref.
// For simplicity in this refactoring, we will make component map values be pointers to the actual structs (e.g. `*Schema` instead of `SchemaRef`).
// The `$ref` field is part of the component struct itself (e.g. `Schema.Ref`).
// However, to explicitly model the choice (value OR reference), Ref types are common.
// The prompt implies a more direct struct usage. The current Schemas already include a $ref field.
// For now, we define them for conceptual clarity; actual usage might simplify to direct pointers.

// SchemaRef can be a Schema object or a reference to one.
type SchemaRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Schema `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// ResponseRef can be a Response object or a reference to one.
type ResponseRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Response `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// ParameterRef can be a Parameter object or a reference to one.
type ParameterRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Parameter `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// ExampleRef can be an Example object or a reference to one.
type ExampleRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Example `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// RequestBodyRef can be a RequestBody object or a reference to one.
type RequestBodyRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*RequestBody `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// HeaderRef can be a Header object or a reference to one.
type HeaderRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Header `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// SecuritySchemeRef can be a SecurityScheme object or a reference to one.
type SecuritySchemeRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*SecurityScheme `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// LinkRef can be a Link object or a reference to one.
type LinkRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Link `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// CallbackRef can be a Callback object (PathItem) or a reference to one.
type CallbackRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	// Callback is essentially a PathItem, so map[string]*PathItem
	// For simplicity here, we'll assume it's a map of PathItems directly or refs to them.
	// Actual callback structure: map[expression]PathItem | map[expression]$ref
	// Using PathItem directly for now as placeholder
	*PathItem `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// Link describes a relationship between an operation and a path.
type Link struct {
	OperationRef string                 `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationId  string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"` // Value can be an expression or any value.
	RequestBody  interface{}            `json:"requestBody,omitempty" yaml:"requestBody,omitempty"` // Value can be an expression or any value.
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Server       *Server                `json:"server,omitempty" yaml:"server,omitempty"`
}
