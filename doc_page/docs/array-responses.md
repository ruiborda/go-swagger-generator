---
sidebar_position: 14
title: Array Responses (v2)
---

# Documenting Array Responses (v2)

This guide shows how to document API endpoints that return arrays of objects in Go-Swagger-Generator v2, for OpenAPI 3.0 specifications.

## Basic Array Response

The most common scenario is an API endpoint that returns an array of objects (e.g., a list of users). Here's how to document this type of response using `SchemaFromDTO` with a slice of pointers to your DTO.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// User DTO
type User struct {
    ID       int64  `json:"id" yaml:"id"`
    Username string `json:"username" yaml:"username"`
    Email    string `json:"email" yaml:"email"`
    FullName string `json:"fullName,omitempty" yaml:"fullName,omitempty"`
    Active   bool   `json:"active" yaml:"active"`
}

// Ensure User DTO is registered as a component schema for referencing
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&User{})

// Swagger documentation for GET /users
var _ = swagger.Swagger().Path("/users"). // Path relative to server URL
    Get(func(op openapi.Operation) {
        op.Summary("List all users").
            Description("Returns an array of all registered users.").
            OperationID("listAllUsersV2").
            Tag("User Operations").
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - list of users returned").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      // Use SchemaFromDTO with a pointer to a slice of DTO pointers for an array response.
                      // This correctly generates an array schema with items referencing the User schema.
                      mt.SchemaFromDTO(&[]*User{}) 
                  })
            })
    }).
    Doc()

// Handler function (example)
func GetAllUsers(c *gin.Context) {
    users := []*User{
        {ID: 1, Username: "alice", Email: "alice@example.com", FullName: "Alice Wonderland", Active: true},
        {ID: 2, Username: "bob", Email: "bob@example.com", FullName: "Bob The Builder", Active: true},
    }
    c.JSON(http.StatusOK, users)
}

// Minimal main and setup for context (example)
// func main() {
//     router := gin.Default()
//     doc := swagger.Swagger()
//     _, _ = doc.ComponentSchemaFromDTO(&User{}) // Register DTO
//     doc.Server("http://localhost:8080/v1", func(s openapi.Server){s.Description("Dev")})
//     // Middleware, path registration, etc.
//     router.GET("/v1/users", GetAllUsers)
//     router.Run(":8080")
// }
```

## Important Points for Array Responses

1.  **Correct DTO Usage**: Pass a pointer to a slice of pointers to your DTO type to `SchemaFromDTO()`. That is, `&[]*YourDtoType{}`.
    *   `&`: Address of the (empty) slice.
    *   `[]`: Denotes a slice.
    *   `*YourDtoType`: The elements of the slice are pointers to `YourDtoType` instances.
    This structure helps the generator understand it's an array of a specific, referenced schema.

2.  **Component Schema Registration**: It's good practice to register your DTO (e.g., `User`) as a component schema using `doc.ComponentSchemaFromDTO(&User{})` or simply `doc.SchemaFromDTO(&User{})` once in your setup. This ensures it's defined in `#/components/schemas/User` and can be referenced cleanly by the array schema.

3.  **OpenAPI 3.0 Output**: This method will generate an OpenAPI 3.0 specification where the response schema is an array, and its `items` property correctly references the schema of `YourDtoType` (e.g., `#/components/schemas/User`).

## Alternative: Manual Array Schema Definition

If you prefer or need more control, you can define the array schema manually within the response content.

```go
package main

import (
    // ... other imports
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// User DTO (as defined before)
// type User struct { ... }

// func main() {
//     doc := swagger.Swagger()
//     // First, ensure the User model is registered as a component schema
//     _, _ = doc.ComponentSchemaFromDTO(&User{})

//     // Then, manually define the array response
//     var _ = doc.Path("/users-manual-array").
//         Get(func(op openapi.Operation) {
//             op.Summary("List users (manual array schema)").Tag("User Operations").
//             Response(http.StatusOK, func(r openapi.Response) {
//                 r.Description("Array of user objects").
//                   Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
//                       mt.Schema(func(s openapi.Schema) {
//                           s.Type("array").
//                             Items(openapi.S().Ref("#/components/schemas/User")) // Reference the User schema
//                       })
//                   })
//             })
//         }).Doc()
//     // ... rest of setup ...
// }
```

## Nested Arrays in DTOs

If your DTO itself contains fields that are arrays (e.g., a `User` DTO with a slice of `Role` DTOs), `SchemaFromDTO` will handle this automatically when you register or use the parent DTO.

```go
// Role DTO
type Role struct {
    ID   int    `json:"id" yaml:"id"`
    Name string `json:"name" yaml:"name"`
}

// UserWithRoles DTO
type UserWithRoles struct {
    ID    int64   `json:"id" yaml:"id"`
    Name  string  `json:"name" yaml:"name"`
    Roles []*Role `json:"roles,omitempty" yaml:"roles,omitempty"` // Nested array of Role DTOs
}

// In your OpenAPI setup:
// doc := swagger.Swagger()
// _, _ = doc.ComponentSchemaFromDTO(&Role{})
// _, _ = doc.ComponentSchemaFromDTO(&UserWithRoles{})

// When used in a response:
// Response(http.StatusOK, func(r openapi.Response) {
//    r.Description("User with roles").
//      Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
//          mt.SchemaFromDTO(&UserWithRoles{}) // Correctly documents UserWithRoles including the Roles array
//      })
// })

// Or an array of UserWithRoles:
// Response(http.StatusOK, func(r openapi.Response) {
//    r.Description("List of users with roles").
//      Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
//          mt.SchemaFromDTO(&[]*UserWithRoles{}) 
//      })
// })
```

Go Swagger Generator v2 will correctly resolve nested DTOs and their array fields into the OpenAPI 3.0 specification. Remember to register all involved DTOs (like `Role` and `UserWithRoles`) as component schemas.