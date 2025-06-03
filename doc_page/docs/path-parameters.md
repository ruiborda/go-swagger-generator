---
sidebar_position: 10
title: Path Parameters (v2)
---

# Documenting Path Parameters (v2)

This guide shows how to document path parameters with Go-Swagger-Generator v2 for OpenAPI 3.0. Path parameters are variable parts of a URL path, enclosed in curly braces (e.g., `/users/{userId}`). They are always required.

## Basic Path Parameter

Here's a simple example of documenting an endpoint with a single path parameter.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Pet DTO
type Pet struct {
    ID     int64  `json:"id,omitempty" yaml:"id,omitempty"`
    Name   string `json:"name" yaml:"name"`
    Status string `json:"status,omitempty" yaml:"status,omitempty"` // e.g., available, pending, sold
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Pet{})

// Swagger documentation for GET /pets/{petId}
var _ = swagger.Swagger().Path("/pets/{petId}"). // Path relative to server URL
    Get(func(op openapi.Operation) {
        op.Summary("Find pet by ID").
            Description("Returns a single pet by its ID.").
            OperationID("getPetByIdV2").
            Tag("Pet Operations").
            PathParameter("petId", func(p openapi.Parameter) {
                p.Description("ID of the pet to retrieve").
                  Required(true). // Path parameters are inherently required
                  Schema(func(s openapi.Schema) { // Define schema for the path parameter
                      s.Type("integer").Format("int64")
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - pet found").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Pet{})
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid ID supplied (e.g., not an integer)")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Pet not found")
            })
    }).
    Doc()

// Handler function (example)
func GetPetByID(c *gin.Context) {
    // petIDStr := c.Param("petId")
    // Convert petIDStr to int64, fetch pet...
    pet := Pet{ID: 1, Name: "Doggie", Status: "available"}
    c.JSON(http.StatusOK, pet)
}
```

## Path Parameter with Validation

This example shows how to document a path parameter with validation constraints (e.g., minimum/maximum values) defined within its schema.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Order DTO
type Order struct {
    ID       int64  `json:"id,omitempty" yaml:"id,omitempty"`
    PetID    int64  `json:"petId,omitempty" yaml:"petId,omitempty"`
    Quantity int32  `json:"quantity,omitempty" yaml:"quantity,omitempty"`
    Status   string `json:"status,omitempty" yaml:"status,omitempty"` // e.g., placed, approved, delivered
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Order{})

// Swagger documentation for GET /store/orders/{orderId}
var _ = swagger.Swagger().Path("/store/orders/{orderId}").
    Get(func(op openapi.Operation) {
        op.Summary("Find purchase order by ID").
            Description("Retrieves a purchase order. Order ID must be between 1 and 1000 (inclusive).").
            OperationID("getOrderByIdV2").
            Tag("Store Operations").
            PathParameter("orderId", func(p openapi.Parameter) {
                p.Description("ID of the purchase order to retrieve").
                  Schema(func(s openapi.Schema) {
                      s.Type("integer").Format("int64").
                        Minimum(1, false).   // Minimum value 1 (inclusive)
                        Maximum(1000, false) // Maximum value 1000 (inclusive)
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - order found").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Order{})
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid Order ID supplied (e.g., out of range or not an integer)")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Order not found")
            })
    }).
    Doc()

// Handler function (example)
func GetOrderByID(c *gin.Context) {
    // orderIDStr := c.Param("orderId")
    // Convert, validate range, fetch order...
    order := Order{ID: 1, PetID: 100, Quantity: 1, Status: "placed"}
    c.JSON(http.StatusOK, order)
}
```

## String Path Parameter with Pattern

Here's an example of documenting a string path parameter with pattern, minLength, and maxLength validations.

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
    ID       int64  `json:"id,omitempty" yaml:"id,omitempty"`
    Username string `json:"username" yaml:"username"`
    Email    string `json:"email,omitempty" yaml:"email,omitempty"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&User{})

// Swagger documentation for GET /users/{username}
var _ = swagger.Swagger().Path("/users/{username}").
    Get(func(op openapi.Operation) {
        op.Summary("Get user by username").
            OperationID("getUserByUsernameV2").
            Tag("User Operations").
            PathParameter("username", func(p openapi.Parameter) {
                p.Description("The username for login (alphanumeric, 3-20 chars)").
                  Schema(func(s openapi.Schema) {
                      s.Type("string").
                        Pattern("^[a-zA-Z0-9]+$"). // Regex for alphanumeric
                        MinLength(3).              // Minimum length
                        MaxLength(20)              // Maximum length
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - user found").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&User{})
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid username supplied (e.g., fails pattern or length constraints)")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
    }).
    Doc()

// Handler function (example)
func GetUserByName(c *gin.Context) {
    // username := c.Param("username")
    // Validate username format (Gin can also do this with bindings in routes), fetch user...
    user := User{ID: 1, Username: "testuser", Email: "testuser@example.com"}
    c.JSON(http.StatusOK, user)
}
```

## Multiple Path Parameters

This example demonstrates documenting an endpoint with multiple path parameters.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Comment DTO
type Comment struct {
    ID      int64  `json:"id,omitempty" yaml:"id,omitempty"`
    PostID  int64  `json:"postId,omitempty" yaml:"postId,omitempty"` // Matches path param name if desired
    Content string `json:"content" yaml:"content"`
    Author  string `json:"author" yaml:"author"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Comment{})

// Swagger documentation for GET /archives/{year}/{month}/posts
var _ = swagger.Swagger().Path("/archives/{year}/{month}/posts").
    Get(func(op openapi.Operation) {
        op.Summary("Get posts from a specific month and year").
            Description("Returns a list of posts for the given year and month.").
            OperationID("getArchivedPostsV2").
            Tag("Archive Operations").
            PathParameter("year", func(p openapi.Parameter) {
                p.Description("The year of the archive (e.g., 2023)").
                  Schema(func(s openapi.Schema) { 
                      s.Type("integer").Format("int32").Minimum(1900, false).Maximum(2100, false)
                  })
            }).
            PathParameter("month", func(p openapi.Parameter) {
                p.Description("The month of the archive (1-12)").
                  Schema(func(s openapi.Schema) { 
                      s.Type("integer").Format("int32").Minimum(1, false).Maximum(12, false)
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - list of posts returned").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&[]*Comment{}) // Example, should be Post DTO
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid year or month supplied")
            })
    }).
    Doc()

// Handler function (example)
func GetArchivedPosts(c *gin.Context) {
    // yearStr := c.Param("year")
    // monthStr := c.Param("month")
    // Convert, validate, fetch posts...
    comments := []*Comment{{ID:1, PostID:1, Content:"A post", Author:"Test"}}
    c.JSON(http.StatusOK, comments)
}
```

## Path Parameter with Enum Values

Here's how to document a path parameter that must be one of a predefined set of values (enum).

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Report DTO (example for response)
type ReportData struct {
    Title string `json:"title" yaml:"title"`
    Data  string `json:"data" yaml:"data"` 
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&ReportData{})

// Swagger documentation for GET /reports/{reportType}/download
var _ = swagger.Swagger().Path("/reports/{reportType}/download").
    Get(func(op openapi.Operation) {
        op.Summary("Download a specific type of report").
            Description("Downloads a report. The type of report must be one of the allowed values.").
            OperationID("downloadReportByTypeV2").
            Tag("Report Operations").
            PathParameter("reportType", func(p openapi.Parameter) {
                p.Description("Type of report to download").
                  Schema(func(s openapi.Schema) {
                      s.Type("string").
                        Enum("sales", "inventory", "customers", "analytics") // Allowed enum values
                  })
            }).
            QueryParameter("format", func(p openapi.Parameter) { // Example of an additional query param
                 p.Description("Format of the report (e.g., pdf, csv)").Required(false).
                   Schema(func(s openapi.Schema){ s.Type("string").Default("pdf") })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Report downloaded successfully (binary data)").
                  Content(mime.ApplicationOctetStream, func(mt openapi.MediaType) { // Example for generic binary file
                      mt.Schema(func(s openapi.Schema) { s.Type("string").Format("binary") })
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid report type or format supplied")
            })
    }).
    Doc()

// Handler function (example)
func DownloadReport(c *gin.Context) {
    // reportType := c.Param("reportType")
    // format := c.Query("format")
    // Generate report based on type and format, then send as binary data...
    c.Data(http.StatusOK, mime.ApplicationOctetStream, []byte("Report binary data..."))
}
```