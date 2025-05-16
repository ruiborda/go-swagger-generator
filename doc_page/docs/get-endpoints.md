---
sidebar_position: 6
title: GET Endpoints
---

# Documenting GET Endpoints

This guide shows how to document GET endpoints with go-swagger-generator using practical examples.

## Basic GET Endpoint

Here's a simple example of documenting a GET endpoint that returns a single resource:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// User DTO
type User struct {
    ID       int64  `json:"id,omitempty"`
    Username string `json:"username,omitempty"`
    Email    string `json:"email,omitempty"`
}

// Swagger documentation for GET /user/{username}
var _ = swagger.Swagger().Path("/user/{username}").
    Get(func(op openapi.Operation) {
        op.Summary("Get user by user name").
            OperationID("getUserByName").
            Tag("user").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            PathParameter("username", func(p openapi.Parameter) {
                p.Description("The name that needs to be fetched").Type("string")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").SchemaFromDTO(&User{})
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
    }).
    Doc()

// Handler function
func GetUserByName(c *gin.Context) {
    username := c.Param("username")
    // Implementation details
    user := User{ID: 1, Username: username, Email: "user@example.com"}
    c.JSON(http.StatusOK, user)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/user/:username", GetUserByName)
}
```

## GET Collection Endpoint

This example shows how to document a GET endpoint that returns a collection of resources:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Pet DTO
type Pet struct {
    ID     int64  `json:"id,omitempty"`
    Name   string `json:"name"`
    Status string `json:"status,omitempty"` // available, pending, sold
}

// Swagger documentation for GET /pet/findByStatus
var _ = swagger.Swagger().Path("/pet/findByStatus").
    Get(func(op openapi.Operation) {
        op.Summary("Finds Pets by status").
            Description("Multiple status values can be provided with comma separated strings").
            OperationID("findPetsByStatus").
            Tag("pet").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            QueryParameter("status", func(p openapi.Parameter) {
                p.Description("Status values that need to be considered for filter").
                    Required(true).
                    Type("array").
                    CollectionFormat("multi").
                    Items(func(item openapi.Schema) {
                        item.Type("string").
                            Enum("available", "pending", "sold").
                            Default("available")
                    })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    Schema(openapi_spec.SchemaEntity{
                        Type:  "array",
                        Items: &openapi_spec.SchemaEntity{Ref: "#/definitions/Pet"},
                    })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid status value")
            })
    }).
    Doc()

// Handler function
func FindPetsByStatus(c *gin.Context) {
    status := c.QueryArray("status")
    // Implementation details
    pets := []Pet{
        {ID: 1, Name: "doggie", Status: "available"},
        {ID: 2, Name: "kitty", Status: "pending"},
    }
    c.JSON(http.StatusOK, pets)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/pet/findByStatus", FindPetsByStatus)
}
```

## GET Endpoint with Multiple Parameters

Here's how to document a GET endpoint that takes multiple parameters:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for GET /user/login
var _ = swagger.Swagger().Path("/user/login").
    Get(func(op openapi.Operation) {
        op.Summary("Logs user into the system").
            OperationID("loginUser").
            Tag("user").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            QueryParameter("username", func(p openapi.Parameter) {
                p.Description("The user name for login").Required(true).Type("string")
            }).
            QueryParameter("password", func(p openapi.Parameter) {
                p.Description("The password for login in clear text").Required(true).Type("string")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    Schema(openapi_spec.SchemaEntity{Type: "string"}).
                    Header("X-Rate-Limit", func(h openapi.Header) {
                        h.Type("integer").Format("int32").Description("calls per hour allowed by the user")
                    }).
                    Header("X-Expires-After", func(h openapi.Header) {
                        h.Type("string").Format("date-time").Description("date in UTC when token expires")
                    })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid username/password supplied")
            })
    }).
    Doc()

// Handler function
func LoginUser(c *gin.Context) {
    username := c.Query("username")
    password := c.Query("password")
    // Implementation details
    c.Header("X-Rate-Limit", "5000")
    c.Header("X-Expires-After", "2025-01-01T00:00:00Z")
    c.String(http.StatusOK, "logged in user session token")
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/user/login", LoginUser)
}
```

## GET Endpoint with Authentication

This example shows how to document a GET endpoint that requires authentication:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for GET /store/inventory
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
            Security("api_key") // Requires API key authentication
    }).
    Doc()

// Handler function
func GetInventory(c *gin.Context) {
    // Implementation details
    c.JSON(http.StatusOK, gin.H{"available": 10, "pending": 5, "sold": 2})
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/store/inventory", GetInventory)
}

// Security definition (in main.go)
func setupSecurity(doc openapi.SwaggerDocBuilder) {
    doc.SecurityDefinition("api_key", func(sd openapi.SecurityScheme) {
        sd.Type("apiKey").Name("api_key").In("header")
    })
}
```