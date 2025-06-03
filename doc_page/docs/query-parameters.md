---
sidebar_position: 11
title: Query Parameters (v2)
---

# Documenting Query Parameters (v2)

This guide shows how to document query parameters with Go-Swagger-Generator v2 for OpenAPI 3.0. Query parameters are appended to the URL after a `?` (e.g., `/search?q=term`).

## Basic Query Parameter

Here's a simple example of documenting an endpoint with an optional query parameter.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// Product DTO (example)
type Product struct { 
    ID int `json:"id" yaml:"id"` 
    Name string `json:"name" yaml:"name"` 
    Category string `json:"category,omitempty" yaml:"category,omitempty"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Product{})

// Swagger documentation for GET /products
var _ = swagger.Swagger().Path("/products"). // Path relative to server URL
    Get(func(op openapi.Operation) {
        op.Summary("List products").
            Description("Get a list of products, optionally filtered by category.").
            OperationID("getProductsV2").
            Tag("Product Operations").
            QueryParameter("category", func(p openapi.Parameter) {
                p.Description("Filter products by category name").
                  Required(false). // This parameter is optional
                  Schema(func(s openapi.Schema) { // Define schema for the query parameter
                      s.Type("string")
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - list of products returned").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&[]*Product{}) // Array of Product DTOs
                  })
            })
    }).
    Doc()

// Handler function (example)
func GetProducts(c *gin.Context) {
    category := c.Query("category") // Returns empty string if not present
    // Implementation to fetch products, filter by category if provided...
    products := []*Product{{ID: 1, Name: "Laptop", Category: category}, {ID: 2, Name: "Mouse", Category: category}}
    c.JSON(http.StatusOK, products)
}
```

## Required Query Parameter with Enum

This example shows a required query parameter that accepts an array of enum values.

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

// Swagger documentation for GET /pets/findByStatus
var _ = swagger.Swagger().Path("/pets/findByStatus").
    Get(func(op openapi.Operation) {
        op.Summary("Finds Pets by status").
            Description("Multiple status values can be provided (e.g., status=available,pending or status=available&status=pending depending on explode value).").
            OperationID("findPetsByStatusV2").
            Tag("Pet Operations").
            QueryParameter("status", func(p openapi.Parameter) {
                p.Description("Status values to filter by").
                  Required(true). // This parameter is required
                  Schema(func(s openapi.Schema) {
                      s.Type("array").
                        Items(func(itemSchema openapi.Schema) {
                            itemSchema.Type("string").
                                Enum("available", "pending", "sold"). // Allowed enum values
                                Default("available") // Default for items if applicable
                        })
                  }).
                  Style("form").   // Default style for query parameters
                  Explode(false) // For status=available,pending,sold (comma-separated)
                                 // Use Explode(true) for status=available&status=pending&status=sold
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - list of pets returned").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&[]*Pet{})
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid status value(s) supplied")
            })
    }).
    Doc()

// Handler function (example)
func FindPetsByStatus(c *gin.Context) {
    // For Explode(false) style (comma-separated): stringValue := c.Query("status") -> split by comma
    // For Explode(true) style (repeated param): arrayValue := c.QueryArray("status")
    statusValues := c.QueryArray("status") // Gin's QueryArray works for both if style is form
    // Implementation...
    pets := []*Pet{{ID: 1, Name: "Buddy", Status: "available"}}
    if len(statusValues) == 0 && c.Query("status") == "" { // Check if param was actually sent if required
         // c.JSON(http.StatusBadRequest, gin.H{"error": "status parameter is required"})
         // return
    }
    c.JSON(http.StatusOK, pets)
}
```

## Multiple Query Parameters with Various Types

An endpoint can have multiple query parameters of different types (string, integer, boolean, date) and with default values or validation.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// SearchResult DTO (example)
type SearchResultItem struct { ID string `json:"id"`; Title string `json:"title"`; }
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&SearchResultItem{})

// Swagger documentation for GET /search
var _ = swagger.Swagger().Path("/search").
    Get(func(op openapi.Operation) {
        op.Summary("Search resources with multiple filters").
            OperationID("searchResourcesV2").
            Tag("Search Operations").
            QueryParameter("q", func(p openapi.Parameter) {
                p.Description("Search query string").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("string").MinLength(1) })
            }).
            QueryParameter("type", func(p openapi.Parameter) {
                p.Description("Resource type to search for").Required(false).
                  Schema(func(s openapi.Schema) {
                      s.Type("string").Enum("article", "video", "podcast", "document")
                  })
            }).
            QueryParameter("page", func(p openapi.Parameter) {
                p.Description("Page number for pagination").Required(false).
                  Schema(func(s openapi.Schema) {
                      s.Type("integer").Format("int32").Default(1).Minimum(1, false)
                  })
            }).
            QueryParameter("pageSize", func(p openapi.Parameter) {
                p.Description("Number of items per page").Required(false).
                  Schema(func(s openapi.Schema) {
                      s.Type("integer").Format("int32").Default(10).Minimum(1, false).Maximum(100, false)
                  })
            }).
            QueryParameter("isFeatured", func(p openapi.Parameter) {
                p.Description("Filter by featured items only").Required(false).
                  Schema(func(s openapi.Schema) { s.Type("boolean").Default(false) })
            }).
            QueryParameter("publishDate", func(p openapi.Parameter) {
                p.Description("Filter by specific publication date (YYYY-MM-DD)").Required(false).
                  Schema(func(s openapi.Schema) { s.Type("string").Format("date") })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Search results returned").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) { // Custom schema for paginated result
                          s.Type("object").
                            Property("totalResults", func(ps openapi.Schema){ ps.Type("integer").Format("int32") }).
                            Property("page", func(ps openapi.Schema){ ps.Type("integer").Format("int32") }).
                            Property("pageSize", func(ps openapi.Schema){ ps.Type("integer").Format("int32") }).
                            Property("results", func(ps openapi.Schema){
                                ps.Type("array").Items(openapi.S().Ref("#/components/schemas/SearchResultItem"))
                            })
                      })
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid parameters supplied")
            })
    }).
    Doc()

// Handler function (example)
func SearchResources(c *gin.Context) {
    // q := c.Query("q")
    // pageStr := c.DefaultQuery("page", "1")
    // isFeaturedStr := c.DefaultQuery("isFeatured", "false")
    // ... parse params, perform search ...
    results := gin.H{
        "totalResults": 1, "page": 1, "pageSize": 10,
        "results": []SearchResultItem{{ID: "1", Title: "Example Result"}},
    }
    c.JSON(http.StatusOK, results)
}
```

**Key points for query parameters:**

*   **`Required(bool)`**: Specifies if the parameter is mandatory.
*   **`Schema(func(s openapi.Schema) { ... })`**: Defines the type, format, default value, enum, and validation constraints (e.g., `Minimum`, `MaxLength`) for the parameter.
*   **Array Serialization (`Style`, `Explode`)**: For array-type query parameters, `Style` and `Explode` control how the array is represented in the URL. Common styles include:
    *   `style: form, explode: false` (default for arrays if not specified): `param=value1,value2,value3`
    *   `style: form, explode: true`: `param=value1&param=value2&param=value3`
    *   Other styles exist (`spaceDelimited`, `pipeDelimited`).
*   **Default Values**: Use `Default(value)` within the schema definition to specify a default if the client doesn't provide the parameter.