package openapi_spec

type License struct {
	Name string `json:"name" yaml:"name"` // Name is required in OAS3
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}
