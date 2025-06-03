---
sidebar_position: 9
title: DELETE Endpoints (v2)
---

# Documenting DELETE Endpoints (v2)

This guide shows how to document DELETE endpoints with Go-Swagger-Generator v2 for OpenAPI 3.0. DELETE requests are used to remove a resource.

## Basic DELETE Endpoint

Here's a simple example of documenting a DELETE endpoint that removes a resource identified by a path parameter.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    // "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime" // Not strictly needed if no response body
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for DELETE /users/{username}
var _ = swagger.Swagger().Path("/users/{username}"). // Path relative to server URL
    Delete(func(op openapi.Operation) {
        op.Summary("Delete user by username").
            Description("This operation deletes a user account.").
            OperationID("deleteUserV2").
            Tag("User Operations").
            PathParameter("username", func(p openapi.Parameter) {
                p.Description("The username of the user to delete").
                  Required(true).
                  Schema(func(s openapi.Schema) { // Define schema for path parameter
                      s.Type("string")
                  })
            }).
            Response(http.StatusNoContent, func(r openapi.Response) { // 204 No Content is typical for successful DELETE
                r.Description("User deleted successfully")
                // No response body for 204
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid username supplied")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
    }).
    Doc()

// Handler function (example)
func DeleteUser(c *gin.Context) {
    // username := c.Param("username")
    // Implementation to delete user...
    c.Status(http.StatusNoContent) // Respond with 204 No Content
}
```

## DELETE with Numeric ID Parameter

This example shows a DELETE endpoint using a numeric ID with validation.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for DELETE /store/orders/{orderId}
var _ = swagger.Swagger().Path("/store/orders/{orderId}").
    Delete(func(op openapi.Operation) {
        op.Summary("Delete purchase order by ID").
            Description("Deletes a single purchase order. Ensure the ID is a positive integer.").
            OperationID("deleteOrderV2").
            Tag("Store Operations").
            PathParameter("orderId", func(p openapi.Parameter) {
                p.Description("ID of the order to delete").
                  Required(true).
                  Schema(func(s openapi.Schema) {
                      s.Type("integer").Format("int64").Minimum(1, false) // e.g., ID must be >= 1
                  })
            }).
            Response(http.StatusNoContent, func(r openapi.Response) {
                r.Description("Order deleted successfully")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid Order ID supplied (e.g., not a positive integer)")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Order not found")
            })
    }).
    Doc()

// Handler function (example)
func DeleteOrder(c *gin.Context) {
    // orderIDStr := c.Param("orderId")
    // Convert orderIDStr to int64, validate, and then delete...
    c.Status(http.StatusNoContent)
}
```

## DELETE with Authentication

This example demonstrates a DELETE endpoint that requires authentication (e.g., OAuth2 with specific scopes).
Assume `OAuth2WriteAccess` security scheme is defined using `doc.ComponentSecurityScheme("OAuth2WriteAccess", ...)`.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for DELETE /pets/{petId}
var _ = swagger.Swagger().Path("/pets/{petId}").
    Delete(func(op openapi.Operation) {
        op.Summary("Deletes a pet").
            OperationID("deletePetV2").
            Tag("Pet Operations").
            Security(map[string][]string{ // Apply security requirement
                "OAuth2WriteAccess": {"write:pets"}, // Example: OAuth2 with 'write:pets' scope
            }).
            PathParameter("petId", func(p openapi.Parameter) {
                p.Description("Pet ID to delete").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("integer").Format("int64") })
            }).
            // Optionally, an API key could be passed as a header parameter if needed for some legacy systems,
            // but it's better to rely on standardized security schemes like OAuth2 or Bearer tokens.
            // Example of a non-standard header parameter if strictly necessary:
            // op.HeaderParameter("X-Legacy-Api-Key", func(p openapi.Parameter) {
            //    p.Description("Legacy API key for this operation (if applicable)").Required(false).
            //      Schema(func(s openapi.Schema){ s.Type("string") })
            // })
            Response(http.StatusNoContent, func(r openapi.Response) {
                r.Description("Pet deleted successfully")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid Pet ID supplied")
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                r.Description("Authentication information is missing or invalid")
            }).
            Response(http.StatusForbidden, func(r openapi.Response) {
                r.Description("Authenticated user does not have permission to delete this pet")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Pet not found")
            })
    }).
    Doc()

// Handler function (example)
func DeletePet(c *gin.Context) {
    // Auth check by middleware...
    // Pet deletion logic...
    c.Status(http.StatusNoContent)
}

// In main OpenAPI config:
// doc.ComponentSecurityScheme("OAuth2WriteAccess", func(ss openapi.SecurityScheme) {
//    ss.Type("oauth2").Flows(func(f openapi.OAuthFlows) { /* ... */ })
// })
```

## Bulk DELETE Endpoint (using Request Body)

If you need to delete multiple resources based on criteria or a list of IDs provided in a request body, you would typically use a POST request (as DELETE with a body is sometimes problematic/discouraged, though allowed by HTTP spec). However, if you must use DELETE with a body:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// BulkDeleteRequest DTO
type BulkDeleteRequest struct {
    IDs []int64 `json:"ids" yaml:"ids"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&BulkDeleteRequest{})

// BulkDeleteResponse DTO
type BulkDeleteResponse struct {
    DeletedCount int     `json:"deletedCount" yaml:"deletedCount"`
    FailedIDs    []int64 `json:"failedIds,omitempty" yaml:"failedIds,omitempty"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&BulkDeleteResponse{})

// Swagger documentation for DELETE /products (with request body)
var _ = swagger.Swagger().Path("/products/bulk-delete"). // Using a more specific path for clarity
    Post(func(op openapi.Operation) { // Changed to POST for safer body handling, could be DELETE
        op.Summary("Bulk delete products by IDs").
            Description("Deletes multiple products based on a list of IDs provided in the request body.").
            OperationID("bulkDeleteProductsV2").
            Tag("Product Operations").
            // Security(...) if needed
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("List of product IDs to delete").Required(true).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&BulkDeleteRequest{})
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Products deleted (or attempted). Check response for details.").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&BulkDeleteResponse{})
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid request body or IDs")
            })
    }).
    Doc()

// Handler function (example)
func BulkDeleteProducts(c *gin.Context) {
    var req BulkDeleteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    // Implementation: Iterate req.IDs, attempt deletion, collect results...
    c.JSON(http.StatusOK, BulkDeleteResponse{DeletedCount: len(req.IDs) - 1, FailedIDs: []int64{req.IDs[0]}} /* example */)
}
```
**Note:** While HTTP DELETE can have a body, it's not universally supported by all clients/proxies. Using POST for bulk operations with a body is often more robust.