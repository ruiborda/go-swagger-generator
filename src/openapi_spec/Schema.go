package openapi_spec

// Schema allows the definition of input and output data types. These types can be objects, primitives, or arrays.
// This object is an extended subset of the JSON Schema Specification Wright Draft 00.
// For more information about the properties, see JSON Schema Core and JSON Schema Validation.
// Unless specified otherwise, JSON Schema properties definitions also apply to the OpenAPI Schema Object.
// Properties 'type', 'format', 'items', 'properties', 'required', 'enum', validation keywords, etc. are defined here.
type Schema struct {
	// JSON Schema fields
	Title                string                 `json:"title,omitempty" yaml:"title,omitempty"`
	MultipleOf           *float64               `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum              *float64               `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum     bool                   `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"` // OAS3 uses boolean, not pointer
	Minimum              *float64               `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum     bool                   `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"` // OAS3 uses boolean, not pointer
	MaxLength            *int                   `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength            *int                   `json:"minLength,omitempty" yaml:"minLength,omitempty"` // Pointer to distinguish 0 from not set
	Pattern              string                 `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MaxItems             *int                   `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems             *int                   `json:"minItems,omitempty" yaml:"minItems,omitempty"` // Pointer to distinguish 0 from not set
	UniqueItems          bool                   `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	MaxProperties        *int                   `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        *int                   `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required             []string               `json:"required,omitempty" yaml:"required,omitempty"`
	Enum                 []interface{}          `json:"enum,omitempty" yaml:"enum,omitempty"`

	// OpenAPI Schema Object specific fields
	Type                 string                 `json:"type,omitempty" yaml:"type,omitempty"` // E.g. "string", "integer", "object", "array"
	AllOf                []*SchemaRef           `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	OneOf                []*SchemaRef           `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf                []*SchemaRef           `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Not                  *SchemaRef             `json:"not,omitempty" yaml:"not,omitempty"`
	Items                *SchemaRef             `json:"items,omitempty" yaml:"items,omitempty"` // For type: "array"
	Properties           map[string]*SchemaRef  `json:"properties,omitempty" yaml:"properties,omitempty"` // For type: "object"
	AdditionalProperties interface{}            `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"` // Can be bool or *SchemaRef
	Description          string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Format               string                 `json:"format,omitempty" yaml:"format,omitempty"`
	Default              interface{}            `json:"default,omitempty" yaml:"default,omitempty"`

	// OpenAPI specific fields
	Nullable        bool                   `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	Discriminator   *Discriminator         `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
	ReadOnly        bool                   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly       bool                   `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	XML             *XMLObject             `json:"xml,omitempty" yaml:"xml,omitempty"`
	ExternalDocs    *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Example         interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Deprecated      bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	// Reference to another schema
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

type Discriminator struct {
	PropertyName string            `json:"propertyName" yaml:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}
