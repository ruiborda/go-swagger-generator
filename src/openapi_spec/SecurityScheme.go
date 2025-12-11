package openapi_spec

type SecurityScheme struct {
	Type        string `json:"type" yaml:"type"` // "apiKey", "http", "oauth2", "openIdConnect"
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// For "apiKey"
	Name string `json:"name,omitempty" yaml:"name,omitempty"` // REQUIRED if type is "apiKey"
	In   string `json:"in,omitempty" yaml:"in,omitempty"`     // REQUIRED if type is "apiKey". "query", "header", "cookie"

	// For "http"
	Scheme       string `json:"scheme,omitempty" yaml:"scheme,omitempty"`             // REQUIRED if type is "http". E.g. "basic", "bearer"
	BearerFormat string `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"` // Hint for "bearer" scheme

	// For "oauth2"
	Flows *OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"` // REQUIRED if type is "oauth2"

	// For "openIdConnect"
	OpenIdConnectUrl string `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"` // REQUIRED if type is "openIdConnect"
}

// OAuthFlows holds the configurations for different OAuth2 flows.
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// OAuthFlow describes a single OAuth2 flow.
type OAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"` // REQUIRED for implicit, authorizationCode
	TokenUrl         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`                 // REQUIRED for password, clientCredentials, authorizationCode
	RefreshUrl       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes" yaml:"scopes"` // REQUIRED. Key is scope name, value is description.
}
