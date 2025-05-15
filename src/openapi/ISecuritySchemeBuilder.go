package openapi

type SecurityScheme interface {
	Type(secType string) SecurityScheme
	Description(description string) SecurityScheme
	Name(name string) SecurityScheme
	In(in string) SecurityScheme
	Flow(flow string) SecurityScheme
	AuthorizationURL(url string) SecurityScheme
	TokenURL(url string) SecurityScheme
	Scope(scopeName string, description string) SecurityScheme
}
