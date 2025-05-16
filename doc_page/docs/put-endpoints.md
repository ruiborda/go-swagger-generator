---
sidebar_position: 8
title: PUT Endpoints
---

# Documenting PUT Endpoints

This guide shows how to document PUT endpoints with go-swagger-generator using practical examples.

## Basic PUT Endpoint

Here's a simple example of documenting a PUT endpoint that updates a resource:

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
    ID        int64  `json:"id,omitempty"`
    Name      string `json:"name"`
    PhotoUrls []string `json:"photoUrls"`
    Status    string `json:"status,omitempty"` // available, pending, sold
}

// Swagger documentation for PUT /pet
var _ = swagger.Swagger().Path("/pet").
    Put(func(op openapi.Operation) {
        op.Summary("Update an existing pet").
            OperationID("updatePet").
            Tag("pet").
            Consumes(string(mime.ApplicationJSON), string(mime.ApplicationXML)).
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("Pet object that needs to be added to the store").
                    Required(true).
                    SchemaFromDTO(&Pet{})
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid ID supplied")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Pet not found")
            }).
            Response(http.StatusMethodNotAllowed, func(r openapi.Response) {
                r.Description("Validation exception")
            })
    }).
    Doc()

// Handler function
func UpdatePet(c *gin.Context) {
    var pet Pet
    if err := c.ShouldBindJSON(&pet); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID supplied"})
        return
    }
    
    // Implementation details (check if pet exists, etc.)
    
    c.JSON(http.StatusOK, pet)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.PUT("/pet", UpdatePet)
}
```

## PUT with Path Parameters

This example shows how to document a PUT endpoint that updates a resource identified by a path parameter:

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
    ID       int64  `json:"id,omitempty"`
    Username string `json:"username,omitempty"`
    FirstName string `json:"firstName,omitempty"`
    LastName  string `json:"lastName,omitempty"`
    Email    string `json:"email,omitempty"`
    Password string `json:"password,omitempty"`
    Phone    string `json:"phone,omitempty"`
}

// Swagger documentation for PUT /user/{username}
var _ = swagger.Swagger().Path("/user/{username}").
    Put(func(op openapi.Operation) {
        op.Summary("Updated user").
            Description("This can only be done by the logged in user.").
            OperationID("updateUser").
            Tag("user").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            PathParameter("username", func(p openapi.Parameter) {
                p.Description("name that need to be updated").Type("string")
            }).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("Updated user object").Required(true).SchemaFromDTO(&User{})
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid user supplied")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
    }).
    Doc()

// Handler function
func UpdateUser(c *gin.Context) {
    username := c.Param("username")
    var user User
    
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user supplied"})
        return
    }
    
    // Implementation details
    c.JSON(http.StatusOK, gin.H{
        "message": "User updated",
        "username": username,
        "user": user,
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.PUT("/user/:username", UpdateUser)
}
```

## PUT with Validation Response Codes

This example demonstrates how to document a PUT endpoint with multiple response codes for different validation scenarios:

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
    Category    string  `json:"category"`
    InStock     bool    `json:"inStock"`
}

// Swagger documentation for PUT /products/{productId}
var _ = swagger.Swagger().Path("/products/{productId}").
    Put(func(op openapi.Operation) {
        op.Summary("Update a product").
            Description("Updates an existing product in the catalog").
            OperationID("updateProduct").
            Tag("products").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON).
            PathParameter("productId", func(p openapi.Parameter) {
                p.Description("ID of the product to update").Type("integer").Format("int64")
            }).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("Product information").Required(true).SchemaFromDTO(&Product{})
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Product updated successfully").SchemaFromDTO(&Product{})
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid product data")
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                r.Description("Not authorized to update products")
            }).
            Response(http.StatusForbidden, func(r openapi.Response) {
                r.Description("Permission denied")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Product not found")
            }).
            Response(http.StatusConflict, func(r openapi.Response) {
                r.Description("Product already exists with conflicting information")
            })
    }).
    Doc()

// Handler function
func UpdateProduct(c *gin.Context) {
    productID := c.Param("productId")
    var product Product
    
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
        return
    }
    
    // Implementation details
    c.JSON(http.StatusOK, product)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.PUT("/products/:productId", UpdateProduct)
}
```

## PUT with Security Requirements

Here's how to document a PUT endpoint that requires authentication:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Settings DTO
type Settings struct {
    NotificationEnabled bool   `json:"notificationEnabled"`
    Theme               string `json:"theme"`
    Language            string `json:"language"`
}

// Swagger documentation for PUT /users/{userId}/settings
var _ = swagger.Swagger().Path("/users/{userId}/settings").
    Put(func(op openapi.Operation) {
        op.Summary("Update user settings").
            Description("Update settings for a specific user. Requires authentication.").
            OperationID("updateUserSettings").
            Tag("user").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON).
            PathParameter("userId", func(p openapi.Parameter) {
                p.Description("ID of the user whose settings to update").Type("integer").Format("int64")
            }).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("Settings object").Required(true).SchemaFromDTO(&Settings{})
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Settings updated successfully").SchemaFromDTO(&Settings{})
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid settings data")
            }).
            Response(http.StatusUnauthorized, func(r openapi.Response) {
                r.Description("Not authenticated")
            }).
            Response(http.StatusForbidden, func(r openapi.Response) {
                r.Description("Not authorized to update this user's settings")
            }).
            Security("api_key") // Requires API key authentication
    }).
    Doc()

// Handler function
func UpdateUserSettings(c *gin.Context) {
    userID := c.Param("userId")
    var settings Settings
    
    if err := c.ShouldBindJSON(&settings); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings data"})
        return
    }
    
    // Implementation details including auth check
    
    c.JSON(http.StatusOK, settings)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.PUT("/users/:userId/settings", UpdateUserSettings)
}

// Security definition (in main.go)
func setupSecurity(doc openapi.SwaggerDocBuilder) {
    doc.SecurityDefinition("api_key", func(sd openapi.SecurityScheme) {
        sd.Type("apiKey").Name("api_key").In("header")
    })
}
```