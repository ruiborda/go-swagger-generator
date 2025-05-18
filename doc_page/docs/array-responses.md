---
sidebar_position: 14
title: Array Responses
---

# Documenting Array Responses

This guide shows how to document API endpoints that return arrays of objects in go-swagger-generator.

## Basic Array Response

The most common scenario is an API endpoint that returns an array of objects. Here's how to document this type of response:

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
    ID       int64  `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    FullName string `json:"fullName"`
    Active   bool   `json:"active"`
}

// Swagger documentation for GET /users
var _ = swagger.Swagger().Path("/users").
    Get(func(op openapi.Operation) {
        op.Summary("List all users").
            Description("Returns an array of all registered users").
            OperationID("listUsers").
            Tag("users").
            Produces(mime.ApplicationJSON).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation").
                    // Use the array format with pointer to support typed arrays
                    SchemaFromDTO(&[]*User{})
            })
    }).
    Doc()

// Handler function
func GetAllUsers(c *gin.Context) {
    // Return a sample array of users
    users := []*User{
        {
            ID:       1,
            Username: "user1",
            Email:    "user1@example.com",
            FullName: "User One",
            Active:   true,
        },
        {
            ID:       2,
            Username: "user2",
            Email:    "user2@example.com",
            FullName: "User Two",
            Active:   true,
        },
    }
    
    c.JSON(http.StatusOK, users)
}
```

## Important Points to Remember

When documenting array responses, there are a few key things to keep in mind:

1. **Use the correct format**: Pass `&[]*YourType{}` to `SchemaFromDTO()` to properly document an array of objects.

2. **Pointer syntax**: Notice that we use a double pointer syntax:
   - `&` - A pointer to the slice
   - `[]` - The slice itself
   - `*YourType` - Pointers to the elements in the slice

3. **Results in Swagger UI**: This will generate proper OpenAPI documentation showing that the endpoint returns an array of objects with the schema defined by your struct.

## Understanding the Internal Implementation

Behind the scenes, the `SchemaFromDTO()` method detects when you pass it a pointer to a slice of object pointers and automatically:

1. Creates a schema for the object type if it doesn't already exist in the definitions
2. Creates an array schema with items referencing the object schema
3. Sets the response schema to the array schema

## Alternative: Manual Schema Definition

If you prefer to define the schema manually, you can also use this approach:

```go
import (
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// First, define the User model
_, _ = doc.DefinitionFromDTO(&User{})

// Then reference it in an array schema
Response(http.StatusOK, func(r openapi.Response) {
    r.Description("Array of user objects").
        Schema(openapi_spec.SchemaEntity{
            Type:  "array",
            Items: &openapi_spec.SchemaEntity{Ref: "#/definitions/User"},
        })
})
```

## Nested Arrays

When dealing with more complex structures that contain arrays, define your structs as normal:

```go
// Department with employees
type Department struct {
    ID        int64    `json:"id"`
    Name      string   `json:"name"`
    Employees []*User  `json:"employees"`
}

// Then use it in your response
Response(http.StatusOK, func(r openapi.Response) {
    r.Description("Department with employee list").
        SchemaFromDTO(&Department{})
})
```

The nested array will be correctly documented in the OpenAPI specification.

## Working Example

For a complete working example, check the array_response example in the examples directory:

```bash
cd examples/array_response
go run main.go
```

This example demonstrates an API that returns both single objects and arrays of objects with proper OpenAPI documentation.