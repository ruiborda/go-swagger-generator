package openapi_spec

type Tag struct {
	Name         string                 `json:"name" yaml:"name"` // REQUIRED
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}
