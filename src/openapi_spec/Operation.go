package openapi_spec

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
)

type OperationEntity struct {
	Tags         []string                     `json:"tags,omitempty"`
	Summary      string                       `json:"summary,omitempty"`
	Description  string                       `json:"description,omitempty"`
	OperationID  string                       `json:"operationId,omitempty"`
	Consumes     []mime.MimeType              `json:"consumes,omitempty"`
	Produces     []mime.MimeType              `json:"produces,omitempty"`
	Parameters   []ParameterEntity            `json:"parameters,omitempty"`
	Responses    map[string]ResponseEntity    `json:"responses"`
	Security     []map[string][]string        `json:"security,omitempty"`
	Deprecated   bool                         `json:"deprecated,omitempty"`
	ExternalDocs *ExternalDocumentationEntity `json:"externalDocs,omitempty"`
}
