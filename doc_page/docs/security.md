---
sidebar_position: 15
title: Security
---

# Security

In this section, you'll learn how to configure authentication and authorization schemes in your API using Go-Swagger-Generator, based on the Pet Store example.

## Defining Security Schemes

Go-Swagger-Generator allows you to define different security schemes for your API. In the Pet Store example, two types of authentication are used:

```go
// In main.go
doc.SecurityDefinition("api_key", func(sd openapi.SecurityScheme) {
    sd.Type("apiKey").Name("api_key").In("header")
})

doc.SecurityDefinition("petstore_auth", func(sd openapi.SecurityScheme) {
    sd.Type("oauth2").
        AuthorizationURL("https://petstore.swagger.io/oauth/authorize").
        Flow("implicit").
        Scope("read:pets", "read your pets").
        Scope("write:pets", "modify pets in your account")
})
```

## Types of Security Schemes

### API Key

```go
doc.SecurityDefinition("api_key", func(sd openapi.SecurityScheme) {
    sd.Type("apiKey").   // Scheme type
        Name("api_key"). // Parameter name
        In("header")     // Location (header, query, or cookie)
})
```

### OAuth2

```go
doc.SecurityDefinition("petstore_auth", func(sd openapi.SecurityScheme) {
    sd.Type("oauth2").                                             // Scheme type
        AuthorizationURL("https://petstore.swagger.io/oauth/authorize"). // Authorization URL
        Flow("implicit").                                          // OAuth2 flow
        Scope("read:pets", "read your pets").                      // Scope and description
        Scope("write:pets", "modify pets in your account")         // Another scope
})
```

## Applying Security to an Operation

Once security schemes are defined, you can apply them to specific operations:

```go
// GetPetByID swagger documentation
var _ = swagger.Swagger().Path("/pet/{petId}").
    Get(func(op openapi.Operation) {
        // ...operation configuration...
        
        // Apply API Key security scheme
        op.Security("api_key")
        
        // Apply OAuth2 security scheme with specific scopes
        op.Security("petstore_auth", "read:pets", "write:pets")
    }).
    Doc()
```

## Real Example of an Operation with OAuth2

```go
// AddPet swagger documentation
var _ = swagger.Swagger().Path("/pet").
    Post(func(operation openapi.Operation) {
        operation.Summary("Add a new pet to the store").
            OperationID("addPet").
            Tag("pet").
            Consumes(string(mime.ApplicationJSON), string(mime.ApplicationXML)).
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("Pet object that needs to be added to the store").
                    Required(true).
                    SchemaFromDTO(&Pet{})
            }).
            Response(http.StatusMethodNotAllowed, func(r openapi.Response) {
                r.Description("Invalid input")
            }).
            // This operation requires OAuth2 authentication with the "write:pets" scope
            Security("petstore_auth", "read:pets", "write:pets")
    }).
    Doc()
```

## Real Example of an Operation with API Key

```go
// GetInventory swagger documentation
var _ = swagger.Swagger().Path("/store/inventory").
    Get(func(op openapi.Operation) {
        op.Summary("Returns pet inventories by status").
            Description("Returns a map of status codes to quantities").
            OperationID("getInventory").
            Tag("store").
            Produces(mime.ApplicationJSON).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        AdditionalProperties: &openapi_spec.SchemaEntity{
                            Type:   "integer",
                            Format: "int32",
                        },
                    })
            }).
            // This operation only requires API Key
            Security("api_key")
    }).
    Doc()
```

## Backend Implementation

Security documentation in Go-Swagger-Generator describes how the user should authenticate, but it doesn't implement the authentication itself. You'll need to implement authentication and authorization logic in your API.

For this, you can create authentication middleware in Gin:

```go
// Example (not included in PetStore) of how to implement API Key authentication
func ApiKeyAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("api_key")
        if apiKey == "" || apiKey != "special-key" {  // Very basic validation for example
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "message": "API key required",
            })
            return
        }
        c.Next()
    }
}

// Use it on specific routes
router.GET("/v2/store/inventory", ApiKeyAuth(), controller.GetInventory)
```

## Best Practices

1. **Don't expose sensitive information**: Don't include tokens, keys, or credentials in your documentation.

2. **Protect sensitive endpoints**: Make sure operations that modify data or access sensitive information require authentication.

3. **Multiple security schemes**: You can use different schemes depending on the level of access required.

4. **Clearly document requirements**: Specify which security schemes and scopes are needed for each operation.

In the next chapter, we'll see [configuration options for production environments](production.md).

---
sidebar_position: 11
title: Security Schemes
---

# Documenting Security Schemes

This guide shows how to document security schemes with go-swagger-generator using practical examples.

## API Key Authentication

Here's how to document API key authentication:

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/middleware"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

func main() {
    router := gin.Default()

    // Set up Swagger UI
    router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
        Enabled:  true,
        JSONPath: "/openapi.json",
        UIPath:   "/swagger",
    }))

    doc := swagger.Swagger()
    
    // Define API key security scheme
    doc.SecurityDefinition("api_key", func(sd openapi.SecurityScheme) {
        sd.Type("apiKey").  // Type is "apiKey"
            Name("X-API-Key").  // Header name
            In("header")  // Location (header, query, cookie)
    })

    // Documentation for protected endpoint
    doc.Path("/protected-resource").
        Get(func(op openapi.Operation) {
            op.Summary("Get protected resource").
                Description("Get a resource that requires API key authentication").
                OperationID("getProtectedResource").
                Tag("secure").
                Produces(mime.ApplicationJSON).
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("successful operation").
                        Schema(openapi_spec.SchemaEntity{
                            Type: "object",
                            Properties: map[string]*openapi_spec.SchemaEntity{
                                "message": {Type: "string"},
                                "timestamp": {Type: "string", Format: "date-time"},
                            },
                        })
                }).
                Response(http.StatusUnauthorized, func(r openapi.Response) {
                    r.Description("API key missing or invalid")
                }).
                Security("api_key")  // Specify that this endpoint requires the api_key security scheme
        }).
        Doc()

    // Handler function
    router.GET("/protected-resource", func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == "" || apiKey != "valid-api-key" { // Simplified validation
            c.JSON(http.StatusUnauthorized, gin.H{"error": "API key missing or invalid"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message": "You have access to the protected resource",
            "timestamp": time.Now().Format(time.RFC3339),
        })
    })

    fmt.Println("Server running on http://localhost:8080")
    fmt.Println("Swagger UI available at http://localhost:8080/swagger")
    _ = router.Run(":8080")
}
```

## OAuth2 Authentication

This example demonstrates how to document OAuth2 authentication:

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/middleware"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

func main() {
    router := gin.Default()

    // Set up Swagger UI
    router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
        Enabled:  true,
        JSONPath: "/openapi.json",
        UIPath:   "/swagger",
    }))

    doc := swagger.Swagger()

    // Define OAuth2 security scheme with implicit flow
    doc.SecurityDefinition("petstore_auth", func(sd openapi.SecurityScheme) {
        sd.Type("oauth2").  // Type is "oauth2"
            AuthorizationURL("https://petstore.swagger.io/oauth/authorize").  // Auth URL
            Flow("implicit").  // OAuth flow type (implicit, password, application, accessCode)
            Scope("read:pets", "read your pets").  // Available scopes
            Scope("write:pets", "modify pets in your account")
    })

    // Documentation for endpoint requiring OAuth2
    doc.Path("/pets").
        Post(func(op openapi.Operation) {
            op.Summary("Add a new pet").
                Description("Add a new pet to the store").
                OperationID("addPet").
                Tag("pets").
                Consumes(string(mime.ApplicationJSON)).
                Produces(mime.ApplicationJSON).
                // Body parameter definition here
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("Pet added successfully")
                }).
                Response(http.StatusUnauthorized, func(r openapi.Response) {
                    r.Description("Authentication required")
                }).
                Response(http.StatusForbidden, func(r openapi.Response) {
                    r.Description("Insufficient permissions")
                }).
                // Require specific OAuth scopes
                Security("petstore_auth", "write:pets", "read:pets")
        }).
        Doc()

    // Documentation for endpoint requiring only read scope
    doc.Path("/pets").
        Get(func(op openapi.Operation) {
            op.Summary("List all pets").
                Description("Returns all pets from the store").
                OperationID("listPets").
                Tag("pets").
                Produces(mime.ApplicationJSON).
                // Response definition here
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("Successful operation")
                }).
                // Only read permission required
                Security("petstore_auth", "read:pets")
        }).
        Doc()

    // Implementation details omitted...

    fmt.Println("Server running on http://localhost:8080")
    fmt.Println("Swagger UI available at http://localhost:8080/swagger")
    _ = router.Run(":8080")
}
```

## Basic Authentication

This example shows how to document basic authentication:

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/middleware"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
    "strings"
)

func main() {
    router := gin.Default()

    // Set up Swagger UI
    router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
        Enabled:  true,
        JSONPath: "/openapi.json",
        UIPath:   "/swagger",
    }))

    doc := swagger.Swagger()

    // Define basic authentication security scheme
    doc.SecurityDefinition("basic_auth", func(sd openapi.SecurityScheme) {
        sd.Type("basic")  // Type is "basic" for Basic Authentication
    })

    // Documentation for protected endpoint
    doc.Path("/admin/dashboard").
        Get(func(op openapi.Operation) {
            op.Summary("Admin Dashboard").
                Description("Access the admin dashboard").
                OperationID("getAdminDashboard").
                Tag("admin").
                Produces(mime.ApplicationJSON, "text/html").
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("successful operation")
                }).
                Response(http.StatusUnauthorized, func(r openapi.Response) {
                    r.Description("Authentication required")
                }).
                Response(http.StatusForbidden, func(r openapi.Response) {
                    r.Description("User is not an admin")
                }).
                Security("basic_auth")  // Requires basic auth
        }).
        Doc()

    // Handler function with basic auth validation
    router.GET("/admin/dashboard", func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        
        // Basic auth validation (simplified)
        if !strings.HasPrefix(auth, "Basic ") {
            c.Header("WWW-Authenticate", "Basic realm=\"Admin Dashboard\"")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            return
        }
        
        // In a real app, decode and validate the credentials
        // For this example, we'll just return a successful response
        
        c.JSON(http.StatusOK, gin.H{
            "message": "Admin dashboard data",
            "stats": gin.H{
                "users": 1024,
                "activeUsers": 512,
                "newSignups": 64,
            },
        })
    })

    fmt.Println("Server running on http://localhost:8080")
    fmt.Println("Swagger UI available at http://localhost:8080/swagger")
    _ = router.Run(":8080")
}
```

## Multiple Security Schemes

This example demonstrates how to document endpoints that support multiple authentication methods:

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/middleware"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

func main() {
    router := gin.Default()

    // Set up Swagger UI
    router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
        Enabled:  true,
        JSONPath: "/openapi.json",
        UIPath:   "/swagger",
    }))

    doc := swagger.Swagger()

    // Define API key security scheme
    doc.SecurityDefinition("api_key", func(sd openapi.SecurityScheme) {
        sd.Type("apiKey").Name("X-API-Key").In("header")
    })

    // Define OAuth2 security scheme
    doc.SecurityDefinition("oauth2", func(sd openapi.SecurityScheme) {
        sd.Type("oauth2").
            AuthorizationURL("https://example.com/oauth/authorize").
            Flow("implicit").
            Scope("read", "Read access").
            Scope("write", "Write access")
    })

    // Documentation for endpoint that accepts either authentication method
    doc.Path("/resources").
        Get(func(op openapi.Operation) {
            op.Summary("Get resource").
                Description("Get a resource using either API key or OAuth2 authentication").
                OperationID("getResource").
                Tag("resources").
                Produces(mime.ApplicationJSON).
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("successful operation")
                }).
                // To specify alternative security requirements, make multiple Security() calls
                Security("api_key").
                Security("oauth2", "read")
        }).
        Doc()

    // Documentation for endpoint that requires both authentication methods
    // Note: This is just an example; requiring both API key and OAuth2 is unusual
    doc.Path("/admin/settings").
        Put(func(op openapi.Operation) {
            op.Summary("Update settings").
                Description("Update system settings (requires both authentication methods)").
                OperationID("updateSettings").
                Tag("admin").
                Consumes(string(mime.ApplicationJSON)).
                Produces(mime.ApplicationJSON).
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("Settings updated")
                }).
                // To require multiple security schemes, provide them in a single Security call
                Security("api_key", "oauth2:write")
        }).
        Doc()

    // Implementation details omitted...

    fmt.Println("Server running on http://localhost:8080")
    fmt.Println("Swagger UI available at http://localhost:8080/swagger")
    _ = router.Run(":8080")
}
```

## JWT Authentication

Here's an example documenting JWT-based authentication:

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/middleware"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
    "strings"
)

func main() {
    router := gin.Default()

    // Set up Swagger UI
    router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
        Enabled:  true,
        JSONPath: "/openapi.json",
        UIPath:   "/swagger",
    }))

    doc := swagger.Swagger()

    // Define JWT security scheme (as a Bearer token)
    doc.SecurityDefinition("jwt_auth", func(sd openapi.SecurityScheme) {
        sd.Type("apiKey").
            Name("Authorization").  // Header name
            In("header").           // Location
            Description("JWT Authorization header using the Bearer scheme. Example: \"Authorization: Bearer {token}\"")
    })

    // Login endpoint to get JWT token
    doc.Path("/auth/login").
        Post(func(op openapi.Operation) {
            op.Summary("Login").
                Description("Login to get JWT token").
                OperationID("loginUser").
                Tag("auth").
                Consumes(string(mime.ApplicationJSON)).
                Produces(mime.ApplicationJSON).
                BodyParameter(func(p openapi.Parameter) {
                    p.Description("Login credentials").
                        Required(true).
                        Schema(openapi_spec.SchemaEntity{
                            Type: "object",
                            Properties: map[string]*openapi_spec.SchemaEntity{
                                "username": {Type: "string"},
                                "password": {Type: "string", Format: "password"},
                            },
                            Required: []string{"username", "password"},
                        })
                }).
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("Login successful").
                        Schema(openapi_spec.SchemaEntity{
                            Type: "object",
                            Properties: map[string]*openapi_spec.SchemaEntity{
                                "token": {Type: "string"},
                                "expires_at": {Type: "string", Format: "date-time"},
                            },
                        })
                }).
                Response(http.StatusUnauthorized, func(r openapi.Response) {
                    r.Description("Invalid credentials")
                })
        }).
        Doc()

    // Protected endpoint using JWT
    doc.Path("/user/profile").
        Get(func(op openapi.Operation) {
            op.Summary("Get user profile").
                Description("Get current user's profile (requires JWT authentication)").
                OperationID("getUserProfile").
                Tag("user").
                Produces(mime.ApplicationJSON).
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("successful operation").
                        Schema(openapi_spec.SchemaEntity{
                            Type: "object",
                            Properties: map[string]*openapi_spec.SchemaEntity{
                                "id": {Type: "integer", Format: "int64"},
                                "username": {Type: "string"},
                                "email": {Type: "string"},
                                "name": {Type: "string"},
                            },
                        })
                }).
                Response(http.StatusUnauthorized, func(r openapi.Response) {
                    r.Description("Missing or invalid JWT token")
                }).
                Security("jwt_auth")
        }).
        Doc()

    // Handler functions (simplified implementation)
    router.POST("/auth/login", func(c *gin.Context) {
        var login struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        
        if err := c.ShouldBindJSON(&login); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
            return
        }
        
        // Validate credentials (simplified)
        if login.Username == "user" && login.Password == "password" {
            // Generate token (simplified)
            c.JSON(http.StatusOK, gin.H{
                "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
                "expires_at": "2025-05-16T00:00:00Z",
            })
            return
        }
        
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
    })

    router.GET("/user/profile", func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        
        // Check for Bearer token (simplified)
        if !strings.HasPrefix(auth, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing JWT token"})
            return
        }
        
        // In a real app, validate the token
        // Here we'll just return a profile
        c.JSON(http.StatusOK, gin.H{
            "id": 12345,
            "username": "user",
            "email": "user@example.com",
            "name": "Jane Doe",
        })
    })

    fmt.Println("Server running on http://localhost:8080")
    fmt.Println("Swagger UI available at http://localhost:8080/swagger")
    _ = router.Run(":8080")
}
```