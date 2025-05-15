package openapi_spec

type TagEntity struct {
	Name         string                       `json:"name"`
	Description  string                       `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentationEntity `json:"externalDocs,omitempty"`
}
