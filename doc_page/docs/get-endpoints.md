---
sidebar_position: 6
title: GET Endpoints (v2)
---

# Documenting GET Endpoints (v2)

This guide shows how to document GET endpoints with Go-Swagger-Generator v2 for OpenAPI 3.0, using practical examples. GET requests are used to retrieve data.

## Basic GET Endpoint (Single Resource)

Here's an example of documenting a GET endpoint that retrieves a single resource by its ID (e.g., a user by username).

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime" // For mime types
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// User DTO (Data Transfer Object)
type User struct {
    ID       int64  `json:"id,omitempty" yaml:"id,omitempty"`
    Username string `json:"username" yaml:"username"`
    Email    string `json:"email,omitempty" yaml:"email,omitempty"`
}

// Ensure User DTO is registered as a schema component
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&User{})

// Swagger documentation for GET /users/{username} (path relative to server URL)
var _ = swagger.Swagger().Path("/users/{username}").
    Get(func(op openapi.Operation) {
        op.Summary("Get user by username").
            OperationID("getUserByNameV2").
            Tag("User Operations").
            PathParameter("username", func(p openapi.Parameter) {
                p.Description("The username of the user to retrieve").
                  Required(true).
                  Schema(func(s openapi.Schema) { // Define schema for path parameter
                      s.Type("string")
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - user found").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) { // Define content type and schema for response
                      mt.SchemaFromDTO(&User{}) // Use User DTO for response schema
                  })
                // Optionally add other content types, e.g., XML
                // r.Content(mime.ApplicationXML, func(mt openapi.MediaType) { 
                //     mt.SchemaFromDTO(&User{})
                // })
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
    }).
    Doc()

// Handler function (example)
func GetUserByName(c *gin.Context) {
    username := c.Param("username")
    // Implementation to fetch user by username...
    user := User{ID: 1, Username: username, Email: username + "@example.com"}
    c.JSON(http.StatusOK, user)
}

// Router setup (example)
// func setupRoutes(router *gin.Engine) {
//     router.GET("/v1/users/:username", GetUserByName) // Assuming server base path /v1
// }
```

## GET Collection Endpoint (Array Response)

This example shows how to document a GET endpoint that returns a collection (array) of resources, possibly with query parameters for filtering.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// Pet DTO
type Pet struct {
    ID     int64  `json:"id,omitempty" yaml:"id,omitempty"`
    Name   string `json:"name" yaml:"name"`
    Status string `json:"status,omitempty" yaml:"status,omitempty"` // e.g., available, pending, sold
}

// Ensure Pet DTO is registered: _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Pet{})

// Swagger documentation for GET /pets/findByStatus
var _ = swagger.Swagger().Path("/pets/findByStatus").
    Get(func(op openapi.Operation) {
        op.Summary("Finds Pets by status").
            Description("Multiple status values can be provided as comma-separated strings (e.g., status=available,pending).").
            OperationID("findPetsByStatusV2").
            Tag("Pet Operations").
            QueryParameter("status", func(p openapi.Parameter) {
                p.Description("Status values to filter by").
                  Required(true).
                  Schema(func(s openapi.Schema) {
                      s.Type("array").
                        Items(func(itemSchema openapi.Schema) {
                            itemSchema.Type("string").
                                Enum("available", "pending", "sold"). // Allowed enum values
                                Default("available")
                        })
                  }).
                  Style("form"). // How the parameter is serialized
                  Explode(false) // For comma-separated values with 'form' style (e.g. status=v1,v2,v3)
                                 // Use Explode(true) for status=v1&status=v2
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - list of pets returned").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      // For an array response, use SchemaFromDTO with a pointer to a slice of pointers
                      mt.SchemaFromDTO(&[]*Pet{}) 
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid status value provided")
            })
    }).
    Doc()

// Handler function (example)
func FindPetsByStatus(c *gin.Context) {
    statusValues := c.QueryArray("status") // Gin handles comma-separated if not using QueryArray
    // status := c.Query("status") // For single string like "available,pending"
    // petsToFilter := strings.Split(status, ",")
    // Implementation...
    pets := []*Pet{
        {ID: 1, Name: "Max", Status: "available"},
        {ID: 2, Name: "Bella", Status: "pending"},
    }
    c.JSON(http.StatusOK, pets)
}
```

## GET Endpoint with Multiple Parameters & Headers

Here's how to document a GET endpoint that takes multiple query parameters and might return custom headers.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
    "time"
)

// Swagger documentation for GET /auth/login (example of query params and headers)
var _ = swagger.Swagger().Path("/auth/login").
    Get(func(op openapi.Operation) {
        op.Summary("Logs user into the system (via GET)").
            OperationID("loginUserV2").
            Tag("Auth Operations").
            QueryParameter("username", func(p openapi.Parameter) {
                p.Description("The username for login").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("string") })
            }).
            QueryParameter("password", func(p openapi.Parameter) {
                p.Description("The password for login (in clear text - not recommended for GET)").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("string").Format("password") })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful login - token returned").
                  Header("X-Rate-Limit", func(h openapi.Header) {
                      h.Description("Calls per hour allowed by the user").
                        Schema(func(s openapi.Schema) { s.Type("integer").Format("int32") })
                  }).
                  Header("X-Expires-After", func(h openapi.Header) {
                      h.Description("Date in UTC when token expires").
                        Schema(func(s openapi.Schema) { s.Type("string").Format("date-time") })
                  }).
                  Content(mime.TextPlain, func(mt openapi.MediaType) { // Example: token as plain text
                      mt.Schema(func(s openapi.Schema) { s.Type("string") })
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid username/password supplied")
            })
    }).
    Doc()

// Handler function (example)
func LoginUser(c *gin.Context) {
    // username := c.Query("username")
    // password := c.Query("password")
    // Implementation...
    c.Header("X-Rate-Limit", "5000")
    c.Header("X-Expires-After", time.Now().Add(1*time.Hour).Format(time.RFC3339))
    c.String(http.StatusOK, "dummy-session-token-via-get")
}
```

## GET Endpoint with Authentication

This example shows how to document a GET endpoint that requires authentication (e.g., an API key).
Assume `ApiKeyAuth` is defined using `doc.ComponentSecurityScheme("ApiKeyAuth", ...)`.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// Swagger documentation for GET /store/inventory
var _ = swagger.Swagger().Path("/store/inventory").
    Get(func(op openapi.Operation) {
        op.Summary("Returns pet inventories by status").
            Description("Returns a map of status codes to quantities. Requires API key.").
            OperationID("getInventoryV2").
            Tag("Store Operations").
            Security(map[string][]string{ // Apply security requirement
                "ApiKeyAuth": {}, // Empty array if no scopes for this scheme type
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - inventory returned").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) {
                          s.Type("object").
                            AdditionalProperties(true, func(addPropSchema openapi.Schema) { // Map-like structure
                                addPropSchema.Type("integer").Format("int32")
                            })
                      })
                  })
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                 r.Description("API key missing or invalid")
            })
    }).
    Doc()

// Handler function (example)
func GetInventory(c *gin.Context) {
    // Auth check would be done by middleware ideally
    c.JSON(http.StatusOK, gin.H{"available": 100, "pending": 20, "sold": 500})
}

// In main OpenAPI config:
// doc.ComponentSecurityScheme("ApiKeyAuth", func(ss openapi.SecurityScheme) {
//    ss.Type("apiKey").Name("X-API-KEY").In("header")
// })
```