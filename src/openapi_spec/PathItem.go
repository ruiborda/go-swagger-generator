package openapi_spec

// PathItemEntity describes the operations available on a single path.
// It can also have parameters common to all operations in this path.
type PathItem struct {
	Ref         string            `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary     string            `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation        `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation        `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation        `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation        `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation        `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation        `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation        `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation        `json:"trace,omitempty" yaml:"trace,omitempty"` // Added Trace for completeness
	Servers     []*Server         `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []*ParameterRef   `json:"parameters,omitempty" yaml:"parameters,omitempty"` // Can be Parameter or Reference
}
