package openapi_spec

type InfoEntity struct {
	Description    string         `json:"description,omitempty"`
	Version        string         `json:"version"`
	Title          string         `json:"title"`
	TermsOfService string         `json:"termsOfService,omitempty"`
	Contact        *ContactEntity `json:"contact,omitempty"`
	License        *LicenseEntity `json:"license,omitempty"`
}
