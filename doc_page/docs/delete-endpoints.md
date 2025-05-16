---
sidebar_position: 9
title: DELETE Endpoints
---

# Documenting DELETE Endpoints

This guide shows how to document DELETE endpoints with go-swagger-generator using practical examples.

## Basic DELETE Endpoint

Here's a simple example of documenting a DELETE endpoint that removes a resource:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for DELETE /user/{username}
var _ = swagger.Swagger().Path("/user/{username}").
    Delete(func(op openapi.Operation) {
        op.Summary("Delete user").
            Description("This can only be done by the logged in user.").
            OperationID("deleteUser").
            Tag("user").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            PathParameter("username", func(p openapi.Parameter) {
                p.Description("The name that needs to be deleted").Type("string")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid username supplied")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
    }).
    Doc()

// Handler function
func DeleteUser(c *gin.Context) {
    username := c.Param("username")
    // Implementation details
    c.JSON(http.StatusOK, gin.H{"message": "User deleted", "username": username})
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.DELETE("/user/:username", DeleteUser)
}
```

## DELETE with Numeric ID Parameter

This example shows how to document a DELETE endpoint with a numeric ID parameter:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for DELETE /store/order/{orderId}
var _ = swagger.Swagger().Path("/store/order/{orderId}").
    Delete(func(op openapi.Operation) {
        op.Summary("Delete purchase order by ID").
            Description("For valid response try integer IDs with positive integer value. Negative or non-integer values will generate API errors").
            OperationID("deleteOrder").
            Tag("store").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            PathParameter("orderId", func(p openapi.Parameter) {
                p.Description("ID of the order that needs to be deleted").
                    Type("integer").Format("int64").Minimum(1, false)
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid ID supplied")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Order not found")
            })
    }).
    Doc()

// Handler function
func DeleteOrder(c *gin.Context) {
    orderID := c.Param("orderId")
    // Implementation details
    c.JSON(http.StatusOK, gin.H{"message": "Order deleted", "orderId": orderID})
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.DELETE("/store/order/:orderId", DeleteOrder)
}
```

## DELETE with Authentication Header

This example demonstrates how to document a DELETE endpoint that requires an authentication header:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for DELETE /pet/{petId}
var _ = swagger.Swagger().Path("/pet/{petId}").
    Delete(func(op openapi.Operation) {
        op.Summary("Deletes a pet").
            OperationID("deletePet").
            Tag("pet").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            HeaderParameter("api_key", func(p openapi.Parameter) {
                p.Description("API key for authentication").Required(false).Type("string")
            }).
            PathParameter("petId", func(p openapi.Parameter) {
                p.Description("Pet id to delete").Type("integer").Format("int64")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid ID supplied")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Pet not found")
            }).
            Security("petstore_auth", "write:pets", "read:pets")
    }).
    Doc()

// Handler function
func DeletePet(c *gin.Context) {
    petID := c.Param("petId")
    apiKey := c.GetHeader("api_key")
    
    // Implementation details including auth check
    
    c.JSON(http.StatusOK, gin.H{"message": "Pet deleted", "petId": petID})
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.DELETE("/pet/:petId", DeletePet)
}

// Security definition (in main.go)
func setupSecurity(doc openapi.SwaggerDocBuilder) {
    doc.SecurityDefinition("petstore_auth", func(sd openapi.SecurityScheme) {
        sd.Type("oauth2").
            AuthorizationURL("https://petstore.swagger.io/oauth/authorize").
            Flow("implicit").
            Scope("read:pets", "read your pets").
            Scope("write:pets", "modify pets in your account")
    })
}
```

## DELETE with Multiple Query Filters

Here's how to document a DELETE endpoint that accepts query parameters for filtering resources to delete:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for DELETE /logs
var _ = swagger.Swagger().Path("/logs").
    Delete(func(op openapi.Operation) {
        op.Summary("Delete system logs by criteria").
            Description("Deletes system logs matching the specified criteria").
            OperationID("deleteLogs").
            Tag("system").
            Produces(mime.ApplicationJSON).
            QueryParameter("fromDate", func(p openapi.Parameter) {
                p.Description("Start date (ISO format)").Required(false).Type("string").Format("date")
            }).
            QueryParameter("toDate", func(p openapi.Parameter) {
                p.Description("End date (ISO format)").Required(false).Type("string").Format("date")
            }).
            QueryParameter("level", func(p openapi.Parameter) {
                p.Description("Log level").Required(false).Type("string").
                    Enum("INFO", "WARNING", "ERROR", "DEBUG")
            }).
            QueryParameter("service", func(p openapi.Parameter) {
                p.Description("Service name").Required(false).Type("string")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Logs deleted").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "deletedCount": {Type: "integer", Format: "int32"},
                        },
                    })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid parameters")
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                r.Description("Not authorized to delete logs")
            }).
            Security("api_key")
    }).
    Doc()

// Handler function
func DeleteLogs(c *gin.Context) {
    fromDate := c.Query("fromDate")
    toDate := c.Query("toDate")
    level := c.Query("level")
    service := c.Query("service")
    
    // Implementation details
    
    c.JSON(http.StatusOK, gin.H{"deletedCount": 42})
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.DELETE("/logs", DeleteLogs)
}
```

## Bulk DELETE Endpoint

This example shows how to document a DELETE endpoint that deletes multiple resources:

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

// Swagger documentation for DELETE /products
var _ = swagger.Swagger().Path("/products").
    Delete(func(op openapi.Operation) {
        op.Summary("Delete multiple products").
            Description("Delete multiple products by their IDs").
            OperationID("bulkDeleteProducts").
            Tag("products").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("Array of product IDs to delete").
                    Required(true).
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "ids": {
                                Type: "array",
                                Items: &openapi_spec.SchemaEntity{
                                    Type:   "integer",
                                    Format: "int64",
                                },
                            },
                        },
                        Required: []string{"ids"},
                    })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Products deleted").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "deletedCount": {Type: "integer", Format: "int32"},
                            "failedIds": {
                                Type: "array",
                                Items: &openapi_spec.SchemaEntity{Type: "integer", Format: "int64"},
                            },
                        },
                    })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid request")
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                r.Description("Authentication required")
            }).
            Security("api_key")
    }).
    Doc()

// Handler function
func BulkDeleteProducts(c *gin.Context) {
    var request struct {
        IDs []int64 `json:"ids"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    
    // Implementation details
    
    c.JSON(http.StatusOK, gin.H{
        "deletedCount": len(request.IDs),
        "failedIds": []int64{},
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.DELETE("/products", BulkDeleteProducts)
}
```