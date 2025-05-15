package openapi_spec

type SchemaEntity struct {
	Ref                  string                       `json:"$ref,omitempty"`
	Format               string                       `json:"format,omitempty"`
	Title                string                       `json:"title,omitempty"`
	Description          string                       `json:"description,omitempty"`
	Default              interface{}                  `json:"default,omitempty"`
	MultipleOf           *float64                     `json:"multipleOf,omitempty"`
	Maximum              *float64                     `json:"maximum,omitempty"`
	ExclusiveMaximum     bool                         `json:"exclusiveMaximum,omitempty"`
	Minimum              *float64                     `json:"minimum,omitempty"`
	ExclusiveMinimum     bool                         `json:"exclusiveMinimum,omitempty"`
	MaxLength            *int                         `json:"maxLength,omitempty"`
	MinLength            *int                         `json:"minLength,omitempty"`
	Pattern              string                       `json:"pattern,omitempty"`
	MaxItems             *int                         `json:"maxItems,omitempty"`
	MinItems             *int                         `json:"minItems,omitempty"`
	UniqueItems          bool                         `json:"uniqueItems,omitempty"`
	MaxProperties        *int                         `json:"maxProperties,omitempty"`
	MinProperties        *int                         `json:"minProperties,omitempty"`
	Required             []string                     `json:"required,omitempty"`
	Enum                 []interface{}                `json:"enum,omitempty"`
	Type                 string                       `json:"type,omitempty"`
	Items                *SchemaEntity                `json:"items,omitempty"`
	AllOf                []*SchemaEntity              `json:"allOf,omitempty"`
	Properties           map[string]*SchemaEntity     `json:"properties,omitempty"`
	AdditionalProperties interface{}                  `json:"additionalProperties,omitempty"` // Puede ser bool o SchemaEntity
	Discriminator        string                       `json:"discriminator,omitempty"`
	ReadOnly             bool                         `json:"readOnly,omitempty"`
	XML                  *XMLObjectEntity             `json:"xml,omitempty"`
	ExternalDocs         *ExternalDocumentationEntity `json:"externalDocs,omitempty"`
	Example              interface{}                  `json:"example,omitempty"`
}
