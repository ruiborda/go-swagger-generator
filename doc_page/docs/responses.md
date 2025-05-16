---
sidebar_position: 13
title: Responses
---

# Documenting Responses

This guide shows how to document API responses with go-swagger-generator using practical examples.

## Basic Response

Here's a simple example of documenting an endpoint with a basic response:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Product DTO
type Product struct {
    ID          int64   `json:"id,omitempty"`
    Name        string  `json:"name"`
    Description string  `json:"description,omitempty"`
    Price       float64 `json:"price"`
    Category    string  `json:"category,omitempty"`
}

// Swagger documentation for GET /products/{productId}
var _ = swagger.Swagger().Path("/products/{productId}").
    Get(func(op openapi.Operation) {
        op.Summary("Get product by ID").
            Description("Returns a single product").
            OperationID("getProductById").
            Tag("products").
            Produces(mime.ApplicationJSON).
            PathParameter("productId", func(p openapi.Parameter) {
                p.Description("ID of product to return").
                    Type("integer").Format("int64")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation").
                    SchemaFromDTO(&Product{})  // Generate schema from struct
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Product not found")
            })
    }).
    Doc()

// Handler function
func GetProductByID(c *gin.Context) {
    productID := c.Param("productId")
    
    // Implementation details
    product := Product{
        ID:          1,
        Name:        "Sample Product",
        Description: "This is a sample product",
        Price:       29.99,
        Category:    "Electronics",
    }
    c.JSON(http.StatusOK, product)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/products/:productId", GetProductByID)
}
```

## Array Response

This example shows how to document an endpoint that returns an array:

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

// Swagger documentation for GET /pets
var _ = swagger.Swagger().Path("/pets").
    Get(func(op openapi.Operation) {
        op.Summary("List pets").
            Description("Returns all pets in the system").
            OperationID("listPets").
            Tag("pets").
            Produces(mime.ApplicationJSON).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Array of pet objects").
                    Schema(openapi_spec.SchemaEntity{
                        Type:  "array",
                        Items: &openapi_spec.SchemaEntity{Ref: "#/definitions/Pet"},
                    })
            })
    }).
    Doc()

// Handler function
func ListPets(c *gin.Context) {
    // Implementation details
    pets := []Pet{
        {ID: 1, Name: "Max", Status: "available"},
        {ID: 2, Name: "Buddy", Status: "pending"},
        {ID: 3, Name: "Charlie", Status: "sold"},
    }
    c.JSON(http.StatusOK, pets)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/pets", ListPets)
}
```

## Multiple Response Status Codes

Here's an example of documenting endpoints with multiple response status codes:

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

// User DTO
type User struct {
    ID       int64  `json:"id,omitempty"`
    Username string `json:"username,omitempty"`
    Email    string `json:"email,omitempty"`
}

// ErrorResponse DTO
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

// Swagger documentation for POST /users
var _ = swagger.Swagger().Path("/users").
    Post(func(op openapi.Operation) {
        op.Summary("Create user").
            Description("Create a new user account").
            OperationID("createUser").
            Tag("users").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("User to create").
                    Required(true).
                    SchemaFromDTO(&User{})
            }).
            Response(http.StatusCreated, func(r openapi.Response) {
                r.Description("User created successfully").
                    SchemaFromDTO(&User{})
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid user data").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "code": {Type: "integer", Format: "int32"},
                            "message": {Type: "string"},
                        },
                    })
            }).
            Response(http.StatusConflict, func(r openapi.Response) {
                r.Description("Username already exists").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "code": {Type: "integer", Format: "int32"},
                            "message": {Type: "string"},
                        },
                    })
            }).
            Response(http.StatusInternalServerError, func(r openapi.Response) {
                r.Description("Internal server error").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "code": {Type: "integer", Format: "int32"},
                            "message": {Type: "string"},
                        },
                    })
            })
    }).
    Doc()

// Handler function
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid user data",
        })
        return
    }
    
    // Implementation details (check if username exists, etc.)
    
    // If everything is successful
    user.ID = 12345
    c.JSON(http.StatusCreated, user)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/users", CreateUser)
}
```

## Response Headers

This example demonstrates how to document response headers:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
    "time"
)

// Swagger documentation for GET /auth/token
var _ = swagger.Swagger().Path("/auth/token").
    Get(func(op openapi.Operation) {
        op.Summary("Get authentication token").
            Description("Get a new authentication token").
            OperationID("getAuthToken").
            Tag("auth").
            Produces(mime.ApplicationJSON).
            QueryParameter("username", func(p openapi.Parameter) {
                p.Description("Username for login").Required(true).Type("string")
            }).
            QueryParameter("password", func(p openapi.Parameter) {
                p.Description("Password for login").Required(true).Type("string")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    Schema(openapi_spec.SchemaEntity{Type: "string"}).
                    Header("X-Rate-Limit", func(h openapi.Header) {
                        h.Description("Rate limit per hour").
                            Type("integer").
                            Format("int32")
                    }).
                    Header("X-Expires-After", func(h openapi.Header) {
                        h.Description("Date in UTC when token expires").
                            Type("string").
                            Format("date-time")
                    }).
                    Header("X-Request-ID", func(h openapi.Header) {
                        h.Description("Unique request identifier").
                            Type("string")
                    })
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                r.Description("Invalid credentials")
            })
    }).
    Doc()

// Handler function
func GetAuthToken(c *gin.Context) {
    username := c.Query("username")
    password := c.Query("password")
    
    // Authentication logic
    
    // Set response headers
    c.Header("X-Rate-Limit", "1000")
    c.Header("X-Expires-After", time.Now().Add(24*time.Hour).UTC().Format(time.RFC3339))
    c.Header("X-Request-ID", "req-123-abc-456-def")
    
    // Return token
    c.String(http.StatusOK, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U")
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/auth/token", GetAuthToken)
}
```

## Paginated Response

Here's how to document a paginated response:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
    "strconv"
)

// Article DTO
type Article struct {
    ID      int64  `json:"id,omitempty"`
    Title   string `json:"title"`
    Content string `json:"content,omitempty"`
    Author  string `json:"author,omitempty"`
}

// PagedResponse represents a generic paginated response
type PagedResponse struct {
    Page       int         `json:"page"`
    PageSize   int         `json:"pageSize"`
    TotalPages int         `json:"totalPages"`
    TotalItems int         `json:"totalItems"`
    Data       interface{} `json:"data"`
}

// Swagger documentation for GET /articles
var _ = swagger.Swagger().Path("/articles").
    Get(func(op openapi.Operation) {
        op.Summary("List articles").
            Description("Get a paginated list of articles").
            OperationID("listArticles").
            Tag("articles").
            Produces(mime.ApplicationJSON).
            QueryParameter("page", func(p openapi.Parameter) {
                p.Description("Page number").
                    Type("integer").
                    Format("int32").
                    Minimum(1, false).
                    Default(1)
            }).
            QueryParameter("pageSize", func(p openapi.Parameter) {
                p.Description("Items per page").
                    Type("integer").
                    Format("int32").
                    Minimum(1, false).
                    Maximum(100, false).
                    Default(20)
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "page": {
                                Type:        "integer",
                                Format:      "int32",
                                Description: "Current page number",
                            },
                            "pageSize": {
                                Type:        "integer",
                                Format:      "int32",
                                Description: "Number of items per page",
                            },
                            "totalPages": {
                                Type:        "integer",
                                Format:      "int32",
                                Description: "Total number of pages",
                            },
                            "totalItems": {
                                Type:        "integer",
                                Format:      "int32",
                                Description: "Total number of items",
                            },
                            "data": {
                                Type:  "array",
                                Items: &openapi_spec.SchemaEntity{Ref: "#/definitions/Article"},
                            },
                        },
                    })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid parameters")
            })
    }).
    Doc()

// Handler function
func ListArticles(c *gin.Context) {
    pageStr := c.DefaultQuery("page", "1")
    pageSizeStr := c.DefaultQuery("pageSize", "20")
    
    page, _ := strconv.Atoi(pageStr)
    pageSize, _ := strconv.Atoi(pageSizeStr)
    
    // Implementation details
    
    // Sample data
    articles := []Article{
        {ID: 1, Title: "First Article", Content: "This is the content", Author: "John Doe"},
        {ID: 2, Title: "Second Article", Content: "More content here", Author: "Jane Smith"},
    }
    
    // Construct paginated response
    response := PagedResponse{
        Page:       page,
        PageSize:   pageSize,
        TotalPages: 5,
        TotalItems: 100,
        Data:       articles,
    }
    
    c.JSON(http.StatusOK, response)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/articles", ListArticles)
}
```

## File Download Response

This example shows how to document a file download response:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for GET /reports/{reportId}/download
var _ = swagger.Swagger().Path("/reports/{reportId}/download").
    Get(func(op openapi.Operation) {
        op.Summary("Download report").
            Description("Download a report file").
            OperationID("downloadReport").
            Tag("reports").
            Produces("application/pdf", "application/vnd.ms-excel", "text/csv").
            PathParameter("reportId", func(p openapi.Parameter) {
                p.Description("ID of the report to download").
                    Type("integer").Format("int64")
            }).
            QueryParameter("format", func(p openapi.Parameter) {
                p.Description("File format").
                    Type("string").
                    Enum("pdf", "excel", "csv").
                    Default("pdf")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("File downloaded successfully")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Report not found")
            })
    }).
    Doc()

// Handler function
func DownloadReport(c *gin.Context) {
    reportID := c.Param("reportId")
    format := c.DefaultQuery("format", "pdf")
    
    // Implementation details
    
    // Set file headers
    var contentType string
    var filename string
    
    switch format {
    case "pdf":
        contentType = "application/pdf"
        filename = "report_" + reportID + ".pdf"
    case "excel":
        contentType = "application/vnd.ms-excel"
        filename = "report_" + reportID + ".xlsx"
    case "csv":
        contentType = "text/csv"
        filename = "report_" + reportID + ".csv"
    }
    
    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Data(http.StatusOK, contentType, []byte("Sample file content"))
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/reports/:reportId/download", DownloadReport)
}
```