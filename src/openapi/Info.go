package openapi

type Info interface {
	Title(title string) Info
	Version(version string) Info
	Description(description string) Info
	TermsOfService(terms string) Info
	Contact(config func(Contact)) Info
	License(config func(License)) Info
}
