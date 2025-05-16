---
sidebar_position: 11
title: Query Parameters
---

# Documenting Query Parameters

This guide shows how to document query parameters with go-swagger-generator using practical examples.

## Basic Query Parameter

Here's a simple example of documenting an endpoint with a query parameter:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for GET /products
var _ = swagger.Swagger().Path("/products").
    Get(func(op openapi.Operation) {
        op.Summary("Get products").
            Description("Get a list of products, optionally filtered by category").
            OperationID("getProducts").
            Tag("products").
            Produces(mime.ApplicationJSON).
            QueryParameter("category", func(p openapi.Parameter) {
                p.Description("Filter products by category").
                    Type("string").
                    Required(false)
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "array",
                        Items: &openapi_spec.SchemaEntity{
                            Ref: "#/definitions/Product",
                        },
                    })
            })
    }).
    Doc()

// Handler function
func GetProducts(c *gin.Context) {
    category := c.Query("category")
    // Implementation details
    c.JSON(http.StatusOK, []gin.H{
        {"id": 1, "name": "Product 1", "category": category},
        {"id": 2, "name": "Product 2", "category": category},
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/products", GetProducts)
}
```

## Required Query Parameter

This example shows how to document a required query parameter:

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
                    Required(true).  // This parameter is required
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

## Multiple Query Parameters

Here's an example of documenting multiple query parameters:

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

// Swagger documentation for GET /search
var _ = swagger.Swagger().Path("/search").
    Get(func(op openapi.Operation) {
        op.Summary("Search resources").
            Description("Search for resources with various filters").
            OperationID("searchResources").
            Tag("search").
            Produces(mime.ApplicationJSON).
            QueryParameter("q", func(p openapi.Parameter) {
                p.Description("Search query string").
                    Type("string").
                    Required(true)
            }).
            QueryParameter("type", func(p openapi.Parameter) {
                p.Description("Resource type").
                    Type("string").
                    Enum("article", "video", "podcast", "document").
                    Required(false)
            }).
            QueryParameter("category", func(p openapi.Parameter) {
                p.Description("Category filter").
                    Type("string").
                    Required(false)
            }).
            QueryParameter("tags", func(p openapi.Parameter) {
                p.Description("Tags filter").
                    Type("array").
                    CollectionFormat("multi").
                    Items(func(item openapi.Schema) {
                        item.Type("string")
                    }).
                    Required(false)
            }).
            QueryParameter("page", func(p openapi.Parameter) {
                p.Description("Page number").
                    Type("integer").
                    Format("int32").
                    Default(1).
                    Minimum(1, false)
            }).
            QueryParameter("pageSize", func(p openapi.Parameter) {
                p.Description("Number of items per page").
                    Type("integer").
                    Format("int32").
                    Default(10).
                    Minimum(1, false).
                    Maximum(100, false)
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "totalResults": {Type: "integer", Format: "int32"},
                            "page": {Type: "integer", Format: "int32"},
                            "pageSize": {Type: "integer", Format: "int32"},
                            "results": {
                                Type: "array",
                                Items: &openapi_spec.SchemaEntity{
                                    Type: "object",
                                },
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
func SearchResources(c *gin.Context) {
    query := c.Query("q")
    resourceType := c.Query("type")
    category := c.Query("category")
    tags := c.QueryArray("tags")
    // Implementation details
    c.JSON(http.StatusOK, gin.H{
        "totalResults": 42,
        "page": 1,
        "pageSize": 10,
        "results": []gin.H{
            {"id": 1, "title": "Result 1", "type": resourceType},
            {"id": 2, "title": "Result 2", "type": resourceType},
        },
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/search", SearchResources)
}
```

## Query Parameter with Validation

This example demonstrates how to document a query parameter with validation constraints:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for GET /products/price-range
var _ = swagger.Swagger().Path("/products/price-range").
    Get(func(op openapi.Operation) {
        op.Summary("Filter products by price").
            Description("Get products within the specified price range").
            OperationID("getProductsByPriceRange").
            Tag("products").
            Produces(mime.ApplicationJSON).
            QueryParameter("minPrice", func(p openapi.Parameter) {
                p.Description("Minimum price").
                    Type("number").
                    Format("double").
                    Minimum(0.0, true).  // exclusive=true means value must be > 0
                    Required(true)
            }).
            QueryParameter("maxPrice", func(p openapi.Parameter) {
                p.Description("Maximum price").
                    Type("number").
                    Format("double").
                    Minimum(0.01, false).  // inclusive minimum
                    Maximum(10000.0, false).  // inclusive maximum
                    Required(true)
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Products found")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid price range")
            })
    }).
    Doc()

// Handler function
func GetProductsByPriceRange(c *gin.Context) {
    minPrice := c.Query("minPrice")
    maxPrice := c.Query("maxPrice")
    // Implementation details
    c.JSON(http.StatusOK, []gin.H{
        {"id": 1, "name": "Budget Item", "price": 19.99},
        {"id": 2, "name": "Premium Item", "price": 99.99},
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/products/price-range", GetProductsByPriceRange)
}
```

## Boolean and Date Query Parameters

Here's how to document boolean and date format query parameters:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for GET /events
var _ = swagger.Swagger().Path("/events").
    Get(func(op openapi.Operation) {
        op.Summary("Get events").
            Description("Get a list of events with various filters").
            OperationID("getEvents").
            Tag("events").
            Produces(mime.ApplicationJSON).
            QueryParameter("startDate", func(p openapi.Parameter) {
                p.Description("Filter events starting from this date (ISO 8601)").
                    Type("string").
                    Format("date")  // ISO 8601 date format (YYYY-MM-DD)
            }).
            QueryParameter("endDate", func(p openapi.Parameter) {
                p.Description("Filter events ending before this date (ISO 8601)").
                    Type("string").
                    Format("date")
            }).
            QueryParameter("includePrivate", func(p openapi.Parameter) {
                p.Description("Include private events in results").
                    Type("boolean").
                    Default(false)
            }).
            QueryParameter("featured", func(p openapi.Parameter) {
                p.Description("Show only featured events").
                    Type("boolean").
                    Default(false)
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid parameters")
            })
    }).
    Doc()

// Handler function
func GetEvents(c *gin.Context) {
    startDate := c.Query("startDate")
    endDate := c.Query("endDate")
    includePrivate := c.Query("includePrivate") == "true"
    featured := c.Query("featured") == "true"
    
    // Implementation details
    c.JSON(http.StatusOK, []gin.H{
        {
            "id": 1,
            "title": "Conference",
            "date": "2025-06-15",
            "private": false,
            "featured": true,
        },
        {
            "id": 2,
            "title": "Workshop",
            "date": "2025-07-20",
            "private": includePrivate,
            "featured": featured,
        },
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/events", GetEvents)
}
```