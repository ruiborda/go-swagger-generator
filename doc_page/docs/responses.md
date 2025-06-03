---
sidebar_position: 13
title: Responses (v2)
---

# Documenting Responses (v2)

This guide shows how to document API responses with Go-Swagger-Generator v2 for OpenAPI 3.0. Responses describe the output of an API operation for different HTTP status codes.

## Basic Response (Single Object)

Here's an example of documenting an endpoint that returns a single object (e.g., a Product DTO) for a successful operation.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// Product DTO
type Product struct {
    ID          int64   `json:"id,omitempty" yaml:"id,omitempty"`
    Name        string  `json:"name" yaml:"name"`
    Description string  `json:"description,omitempty" yaml:"description,omitempty"`
    Price       float64 `json:"price" yaml:"price"`
    Category    string  `json:"category,omitempty" yaml:"category,omitempty"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Product{})

// Swagger documentation for GET /products/{productId}
var _ = swagger.Swagger().Path("/products/{productId}"). // Path relative to server URL
    Get(func(op openapi.Operation) {
        op.Summary("Get product by ID").
            Description("Returns a single product by its ID.").
            OperationID("getProductByIdV2").
            Tag("Product Operations").
            PathParameter("productId", func(p openapi.Parameter) {
                p.Description("ID of product to return").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("integer").Format("int64") })
            }).
            Response(http.StatusOK, func(r openapi.Response) { // Define response for HTTP 200 OK
                r.Description("Successful operation - product found").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) { // Specify content type
                      mt.SchemaFromDTO(&Product{})  // Use Product DTO for the response schema
                  })
                // Optionally, add other content types like XML
                // r.Content(mime.ApplicationXML, func(mt openapi.MediaType) {
                //    mt.SchemaFromDTO(&Product{})
                // })
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Product not found")
                // Optionally, define a schema for the error response body
                // r.Content(mime.ApplicationJSON, func(mt openapi.MediaType) { 
                //     mt.Schema(openapi.S().Ref("#/components/schemas/ErrorModel")) 
                // })
            })
    }).
    Doc()

// Handler function (example)
func GetProductByID(c *gin.Context) {
    // productIdStr := c.Param("productId")
    // Fetch product...
    product := Product{ID: 1, Name: "Sample Laptop", Price: 999.99, Category: "Electronics"}
    c.JSON(http.StatusOK, product)
}
```

## Array Response

This example shows how to document an endpoint that returns an array of objects.

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
    Status string `json:"status,omitempty" yaml:"status,omitempty"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Pet{})

// Swagger documentation for GET /pets
var _ = swagger.Swagger().Path("/pets").
    Get(func(op openapi.Operation) {
        op.Summary("List all pets").
            Description("Returns an array of all pets in the system.").
            OperationID("listPetsV2").
            Tag("Pet Operations").
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - array of pets returned").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      // For an array response, use SchemaFromDTO with a pointer to a slice of DTO pointers
                      mt.SchemaFromDTO(&[]*Pet{}) 
                  })
            })
    }).
    Doc()

// Handler function (example)
func ListPets(c *gin.Context) {
    pets := []*Pet{{ID: 1, Name: "Buddy"}, {ID: 2, Name: "Lucy"}}
    c.JSON(http.StatusOK, pets)
}
```

## Multiple Response Status Codes with Different Schemas

Endpoints often return different responses (and schemas) for different status codes (e.g., success vs. error codes).

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// User DTO & ErrorResponse DTO
type User struct { ID int64 `json:"id"`; Username string `json:"username"`; Email string `json:"email"`; }
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&User{})
type ErrorResponse struct { Code int `json:"code"`; Message string `json:"message"`; Details string `json:"details,omitempty"`;}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&ErrorResponse{})

// Swagger documentation for POST /users
var _ = swagger.Swagger().Path("/users").
    Post(func(op openapi.Operation) {
        op.Summary("Create a new user").OperationID("createUserV2").Tag("User Operations").
            RequestBody( /* ... define request body ... */ ).
            Response(http.StatusCreated, func(r openapi.Response) {
                r.Description("User created successfully").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) { mt.SchemaFromDTO(&User{}) })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid user data supplied").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) { mt.SchemaFromDTO(&ErrorResponse{}) })
            }).
            Response(http.StatusConflict, func(r openapi.Response) {
                r.Description("Username or email already exists").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) { mt.SchemaFromDTO(&ErrorResponse{}) })
            }).
            Response(http.StatusInternalServerError, func(r openapi.Response) {
                r.Description("Internal server error encountered").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) { mt.SchemaFromDTO(&ErrorResponse{}) })
            })
    }).
    Doc()

// Handler function (example)
func CreateUser(c *gin.Context) {
    // ... bind user data ...
    // if validationFails { c.JSON(http.StatusBadRequest, ErrorResponse{...}); return }
    // if conflict { c.JSON(http.StatusConflict, ErrorResponse{...}); return }
    createdUser := User{ID: 1, Username: "newuser", Email: "new@example.com"}
    c.JSON(http.StatusCreated, createdUser)
}
```

## Response Headers

This example demonstrates how to document custom headers returned in a response.

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

// Swagger documentation for POST /auth/login (example with headers)
var _ = swagger.Swagger().Path("/auth/login").
    Post(func(op openapi.Operation) {
        op.Summary("User login").OperationID("loginUserV2").Tag("Auth Operations").
            RequestBody( /* ... define login request body ... */ ).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Login successful, token returned in body, rate limits in headers").
                  Header("X-Rate-Limit-Limit", func(h openapi.Header) {
                      h.Description("The number of allowed requests in the current period").
                        Schema(func(s openapi.Schema) { s.Type("integer") })
                  }).
                  Header("X-Rate-Limit-Remaining", func(h openapi.Header) {
                      h.Description("The number of remaining requests in the current period").
                        Schema(func(s openapi.Schema) { s.Type("integer") })
                  }).
                  Header("X-Token-Expires-At", func(h openapi.Header) {
                      h.Description("Timestamp when the authentication token expires").
                        Schema(func(s openapi.Schema) { s.Type("string").Format("date-time") })
                  }).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) { // Schema for login token response
                          s.Type("object").Property("token", func(ps openapi.Schema){ ps.Type("string") })
                      })
                  })
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                r.Description("Invalid credentials")
            })
    }).
    Doc()

// Handler function (example)
func LoginUser(c *gin.Context) {
    // ... authenticate user ...
    c.Header("X-Rate-Limit-Limit", "1000")
    c.Header("X-Rate-Limit-Remaining", "999")
    c.Header("X-Token-Expires-At", time.Now().Add(1*time.Hour).Format(time.RFC3339))
    c.JSON(http.StatusOK, gin.H{"token": "dummy-jwt-token"})
}
```

## Paginated Response

Here's how to document a response that is paginated, often returning metadata alongside the data array.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
    // "strconv"
)

// Article DTO
type Article struct { ID int64 `json:"id"`; Title string `json:"title"`; }
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Article{})

// PagedResponse - Generic DTO for paginated responses. Could also be defined inline.
// type PagedResponse struct {
// Page       int         `json:"page"`
// PageSize   int         `json:"pageSize"`
// TotalPages int         `json:"totalPages"`
// TotalItems int         `json:"totalItems"`
// Data       interface{} `json:"data"` // Use specific type like []*Article here
// }

// Swagger documentation for GET /articles
var _ = swagger.Swagger().Path("/articles").
    Get(func(op openapi.Operation) {
        op.Summary("List articles with pagination").OperationID("listArticlesV2").Tag("Article Operations").
            QueryParameter("page", func(p openapi.Parameter) { /* ... */ }).
            QueryParameter("pageSize", func(p openapi.Parameter) { /* ... */ }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - paginated list of articles").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) { // Inline schema for paginated response structure
                          s.Type("object").
                            Property("page", func(ps openapi.Schema){ ps.Type("integer").Description("Current page number") }).
                            Property("pageSize", func(ps openapi.Schema){ ps.Type("integer").Description("Number of items per page") }).
                            Property("totalPages", func(ps openapi.Schema){ ps.Type("integer").Description("Total number of pages") }).
                            Property("totalItems", func(ps openapi.Schema){ ps.Type("integer").Description("Total number of items") }).
                            Property("data", func(ps openapi.Schema){ // The actual data array
                                ps.Type("array").Items(openapi.S().Ref("#/components/schemas/Article"))
                            })
                      })
                  })
            })
    }).
    Doc()

// Handler function (example)
func ListArticles(c *gin.Context) {
    // page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    // pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
    articles := []*Article{{ID: 1, Title: "Intro to OpenAPI"}}
    response := gin.H{
        "page": 1, "pageSize": 20, "totalPages": 5, "totalItems": 100,
        "data": articles,
    }
    c.JSON(http.StatusOK, response)
}
```

## File Download Response

Documenting a file download, where the response body is the file content.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// Swagger documentation for GET /reports/{reportId}/download
var _ = swagger.Swagger().Path("/reports/{reportId}/download").
    Get(func(op openapi.Operation) {
        op.Summary("Download a report file").OperationID("downloadReportV2").Tag("Report Operations").
            PathParameter("reportId", func(p openapi.Parameter) { /* ... */ }).
            QueryParameter("format", func(p openapi.Parameter) { // e.g., format=pdf or format=csv
                p.Description("File format for download").Required(false).
                  Schema(func(s openapi.Schema){ s.Type("string").Enum("pdf", "csv", "xlsx").Default("pdf") })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Report file downloaded successfully").
                  // Multiple content types can be specified if the format query param changes the output type
                  Content(mime.ApplicationPdf, func(mt openapi.MediaType) { // For PDF
                      mt.Schema(func(s openapi.Schema) { s.Type("string").Format("binary") })
                  }).
                  Content("text/csv", func(mt openapi.MediaType) { // For CSV
                      mt.Schema(func(s openapi.Schema) { s.Type("string").Format("binary") })
                  })
                // For a single known binary type, only one Content entry is needed.
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Report not found or format unavailable")
            })
    }).
    Doc()

// Handler function (example)
func DownloadReport(c *gin.Context) {
    // reportID := c.Param("reportId")
    format := c.DefaultQuery("format", "pdf")
    var contentType string
    var fileData []byte // Your report data
    
    switch format {
    case "pdf":
        contentType = mime.ApplicationPdf
        fileData = []byte("%PDF-1.4 sample PDF content")
    case "csv":
        contentType = "text/csv"
        fileData = []byte("col1,col2\nval1,val2")
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported format"})
        return
    }
    c.Header("Content-Disposition", "attachment; filename=report."+format)
    c.Data(http.StatusOK, contentType, fileData)
}
```