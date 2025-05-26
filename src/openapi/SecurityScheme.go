package openapi

type SecurityScheme interface {
	Type(secType string) SecurityScheme // apiKey, http, oauth2, openIdConnect
	Description(description string) SecurityScheme

	// For apiKey
	Name(name string) SecurityScheme
	In(in string) SecurityScheme // query, header, cookie

	// For http
	Scheme(scheme string) SecurityScheme // basic, bearer, etc.
	BearerFormat(format string) SecurityScheme

	// For oauth2
	Flows(config func(OAuthFlows)) SecurityScheme

	// For openIdConnect
	OpenIDConnectURL(url string) SecurityScheme

	Ref(ref string) SecurityScheme
}

type OAuthFlows interface {
	Implicit(config func(OAuthFlow)) OAuthFlows
	Password(config func(OAuthFlow)) OAuthFlows
	ClientCredentials(config func(OAuthFlow)) OAuthFlows
	AuthorizationCode(config func(OAuthFlow)) OAuthFlows
}

type OAuthFlow interface {
	AuthorizationURL(url string) OAuthFlow
	TokenURL(url string) OAuthFlow
	RefreshURL(url string) OAuthFlow
	Scope(scopeName string, description string) OAuthFlow
	Scopes(scopes map[string]string) OAuthFlow
}

type Callback interface {
	PathItem(expression string, config func(item PathItem)) Callback
	Ref(ref string) Callback
}
