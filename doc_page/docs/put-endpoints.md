---
sidebar_position: 8
title: PUT Endpoints (v2)
---

# Documenting PUT Endpoints (v2)

This guide shows how to document PUT endpoints with Go-Swagger-Generator v2 for OpenAPI 3.0. PUT requests are typically used to update an existing resource entirely or create it if it doesn't exist (though idempotency is key).

## Basic PUT Endpoint (Full Resource Update)

Here's an example of documenting a PUT endpoint that updates an existing resource. The request body usually contains the complete representation of the resource.

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
    ID        int64    `json:"id" yaml:"id"` // ID should be part of the DTO for updates
    Name      string   `json:"name" yaml:"name"`
    PhotoUrls []string `json:"photoUrls" yaml:"photoUrls"`
    Status    string   `json:"status,omitempty" yaml:"status,omitempty"` // e.g., available, pending, sold
}

// Ensure Pet DTO is registered: _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Pet{})

// Swagger documentation for PUT /pets/{petId}
// Assumes the resource is identified by petId in the path.
var _ = swagger.Swagger().Path("/pets/{petId}"). // Path relative to server URL
    Put(func(op openapi.Operation) {
        op.Summary("Update an existing pet").
            OperationID("updatePetV2").
            Tag("Pet Operations").
            PathParameter("petId", func(p openapi.Parameter) {
                p.Description("ID of pet to update").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("integer").Format("int64") })
            }).
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("Pet object that needs to be updated in the store. ID in body should match path or be omitted.").
                  Required(true).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Pet{}) // Full Pet object for update
                  })
                // Optionally add other consumable types e.g. XML
                // rb.Content(mime.ApplicationXML, func(mt openapi.MediaType) {
                //    mt.SchemaFromDTO(&Pet{})
                // })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Pet updated successfully").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Pet{}) // Return the updated pet
                  })
            }).
            Response(http.StatusCreated, func(r openapi.Response) {
                r.Description("Pet created successfully (if PUT allows creation)").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Pet{}) // Return the created pet
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid ID supplied or invalid Pet data")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Pet not found")
            })
    }).
    Doc()

// Handler function (example)
func UpdatePet(c *gin.Context) {
    // petID := c.Param("petId")
    var pet Pet
    if err := c.ShouldBindJSON(&pet); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Pet data"})
        return
    }
    // Implementation: find pet by petID, update its fields with data from 'pet' DTO.
    // Ensure pet.ID from body matches petID from path, or handle appropriately.
    c.JSON(http.StatusOK, pet) // Return updated pet
}
```

## PUT with Validation Response Codes

This example demonstrates documenting a PUT endpoint with multiple response codes for different validation scenarios and outcomes.

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
    Category    string  `json:"category" yaml:"category"`
    InStock     bool    `json:"inStock" yaml:"inStock"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Product{})

// Swagger documentation for PUT /products/{productId}
var _ = swagger.Swagger().Path("/products/{productId}").
    Put(func(op openapi.Operation) {
        op.Summary("Update a product").
            Description("Updates an existing product in the catalog or creates it if allowed by implementation.").
            OperationID("updateProductV2").
            Tag("Product Operations").
            PathParameter("productId", func(p openapi.Parameter) {
                p.Description("ID of the product to update or create").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("integer").Format("int64") })
            }).
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("Product information").Required(true).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Product{})
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Product updated successfully").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Product{})
                  })
            }).
            Response(http.StatusCreated, func(r openapi.Response) {
                r.Description("Product created successfully (if PUT allows creation)").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Product{})
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid product data supplied")
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                r.Description("Authentication required to update products")
            }).
            Response(http.StatusForbidden, func(r openapi.Response) {
                r.Description("Authorization failed (e.g., insufficient permissions)")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Product not found (if PUT only updates existing and it's missing)")
            })
            // Potentially StatusConflict (409) if update violates a unique constraint.
    }).
    Doc()

// Handler function (example)
func UpdateProduct(c *gin.Context) {
    // productID := c.Param("productId")
    var product Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
        return
    }
    // Implementation: Check auth, find product, update, or create...
    c.JSON(http.StatusOK, product) // Or http.StatusCreated if new
}
```

## PUT with Security Requirements

Here's how to document a PUT endpoint that requires authentication (e.g., an API key or OAuth2 token).
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

// Settings DTO
type Settings struct {
    NotificationEnabled bool   `json:"notificationEnabled" yaml:"notificationEnabled"`
    Theme               string `json:"theme" yaml:"theme"`
    Language            string `json:"language" yaml:"language"`
}
// _, _ = swagger.Swagger().ComponentSecuritySchemeFromDTO(&Settings{})

// Swagger documentation for PUT /users/{userId}/settings
var _ = swagger.Swagger().Path("/users/{userId}/settings").
    Put(func(op openapi.Operation) {
        op.Summary("Update user settings").
            Description("Update settings for a specific user. Requires authentication.").
            OperationID("updateUserSettingsV2").
            Tag("User Operations").
            Security(map[string][]string{ // Apply security requirement
                "ApiKeyAuth": {}, // Assuming ApiKeyAuth is defined
            }).
            PathParameter("userId", func(p openapi.Parameter) {
                p.Description("ID of the user whose settings to update").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("integer").Format("int64") })
            }).
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("Settings object").Required(true).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Settings{})
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Settings updated successfully").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Settings{})
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid settings data")
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                r.Description("Authentication credentials missing or invalid")
            }).
            Response(http.StatusForbidden, func(r openapi.Response) {
                r.Description("Not authorized to update this user's settings")
            })
    }).
    Doc()

// Handler function (example)
func UpdateUserSettings(c *gin.Context) {
    // userID := c.Param("userId")
    // Auth check would be done by middleware ideally
    var settings Settings
    if err := c.ShouldBindJSON(&settings); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings data"})
        return
    }
    // Implementation...
    c.JSON(http.StatusOK, settings)
}

// In main OpenAPI config:
// doc.ComponentSecurityScheme("ApiKeyAuth", func(ss openapi.SecurityScheme) {
//    ss.Type("apiKey").Name("X-API-KEY").In("header")
// })
```