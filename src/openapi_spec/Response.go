package openapi_spec

type ResponseEntity struct {
	Description string                  `json:"description"`
	Schema      *SchemaEntity           `json:"schema,omitempty"`
	Headers     map[string]HeaderEntity `json:"headers,omitempty"`
	Examples    map[string]interface{}  `json:"examples,omitempty"`
}
