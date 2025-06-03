---
sidebar_position: 15
title: Security Schemes (v2)
---

# Security Schemes (v2)

In this section, you'll learn how to configure authentication and authorization schemes in your API using Go-Swagger-Generator v2 for OpenAPI 3.0. Security schemes are defined in the `components.securitySchemes` section of the OpenAPI document and then applied to operations or globally.

## Defining Security Schemes (Components)

Go-Swagger-Generator v2 allows you to define different security schemes using `doc.ComponentSecurityScheme()`.

```go
// In your main configuration (e.g., ConfigureOpenAPI function)
doc := swagger.Swagger()

// API Key Example (in header)
doc.ComponentSecurityScheme("ApiKeyAuth", func(ss openapi.SecurityScheme) {
    ss.Type("apiKey").
       Name("X-API-KEY"). // Name of the header or query parameter
       In("header").       // Location: "query", "header" or "cookie"
       Description("API Key for authentication.")
})

// OAuth2 Example (Implicit Flow - common for older UIs, consider PKCE for newer apps)
doc.ComponentSecurityScheme("OAuth2Implicit", func(ss openapi.SecurityScheme) {
    ss.Type("oauth2").
       Description("OAuth2 Implicit Grant for accessing pet store resources.").
       Flows(func(flows openapi.OAuthFlows) {
           flows.Implicit(func(flow openapi.OAuthFlow) {
               flow.AuthorizationURL("https://petstore.swagger.io/oauth/authorize").
                  Scope("write:pets", "Modify pets in your account").
                  Scope("read:pets", "Read your pets")
               // TokenURL is not used for implicit flow
           })
       })
})

// JWT Bearer Authentication Example
doc.ComponentSecurityScheme("BearerAuth", func(ss openapi.SecurityScheme) {
    ss.Type("http").
       Scheme("bearer").
       BearerFormat("JWT"). // Optional, a hint
       Description("JWT Authorization header using the Bearer scheme. Example: \"Authorization: Bearer {token}\"")
})

// Basic Authentication Example
doc.ComponentSecurityScheme("BasicAuth", func(ss openapi.SecurityScheme) {
    ss.Type("http").
       Scheme("basic").
       Description("HTTP Basic authentication.")
})
```

## Applying Security to an Operation

Once security schemes are defined as components, you can apply them to specific operations. An operation can require one or more security schemes. If multiple schemes are listed in a single security requirement object, they are treated as AND. To specify OR, you would list multiple security requirement objects (though the library's direct support for OR on a single operation needs verification - typically, global security or a single requirement per op is common).

```go
// Example: Applying ApiKeyAuth to an operation
var _ = swagger.Swagger().Path("/secure/data").
    Get(func(op openapi.Operation) {
        op.Summary("Get secure data").
           Tag("Secure Endpoints").
           // Apply ApiKeyAuth security scheme
           Security(map[string][]string{
               "ApiKeyAuth": {}, // Empty array for scopes if not applicable
           })
        // ... other operation configurations (responses, etc.)
    }).
    Doc()

// Example: Applying OAuth2 with specific scopes
var _ = swagger.Swagger().Path("/pets").
    Post(func(op openapi.Operation) {
        op.Summary("Add a new pet").
           Tag("Pets").
           Security(map[string][]string{
               "OAuth2Implicit": {"write:pets", "read:pets"}, // Scopes required for this operation
           })
        // ... requestBody, responses, etc.
    }).
    Doc()

// Example: Applying JWT Bearer Authentication
var _ = swagger.Swagger().Path("/user/profile").
    Get(func(op openapi.Operation) {
        op.Summary("Get user profile").
           Tag("User").
           Security(map[string][]string{
               "BearerAuth": {},
           })
        // ... responses, etc.
    }).
    Doc()
```

## Global Security

You can also apply security schemes globally to all operations, unless overridden at the operation level.

```go
doc := swagger.Swagger()
// Define ComponentSecurityScheme("ApiKeyAuth", ...) as above

// Apply ApiKeyAuth globally
doc.GlobalSecurity(map[string][]string{
    "ApiKeyAuth": {},
})

// To make a specific operation unsecured when global security is active:
op.Security(map[string][]string{}) // An empty security requirement
// Or more explicitly, if the library supports it, an empty array for the operation's security field.
// op.Security([]map[string][]string{}...) // Check library specifics for this syntax
```

## Backend Implementation

Documenting security schemes in OpenAPI describes how clients should authenticate. **You must implement the actual security logic (e.g., authentication middleware) in your Gin application.**

```go
// Example Gin middleware for API Key authentication (simplified)
func ApiKeyAuthMiddleware(expectedApiKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-KEY") // Matches 'Name' in ApiKeyAuth definition
        if apiKey == "" || apiKey != expectedApiKey {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required or invalid"})
            return
        }
        c.Next()
    }
}

// Usage in Gin:
// router.GET("/v2/secure/data", ApiKeyAuthMiddleware("your-secret-api-key"), controller.GetSecureDataHandler)
```

## Best Practices for v2/OpenAPI 3.0

1.  **Use `components.securitySchemes`**: Define all security schemes centrally using `doc.ComponentSecurityScheme()`.
2.  **Be Specific with OAuth2 Flows**: OpenAPI 3.0 has detailed flow definitions (authorizationCode, implicit, password, clientCredentials). Choose the one that matches your IdP.
3.  **HTTP Authentication**: Use `type: http` for schemes like Basic and Bearer (JWT).
4.  **Clear Descriptions**: Provide good descriptions for your security schemes and how to use them (e.g., where to get an API key or token).
5.  **Scopes**: Clearly define and use OAuth2 scopes if applicable.

In the next chapter, we'll see [configuration options for production environments](production.md).