---
sidebar_position: 7
title: POST Endpoints
---

# Documenting POST Endpoints

This guide shows how to document POST endpoints with go-swagger-generator using practical examples.

## Basic POST Endpoint

Here's a simple example of documenting a POST endpoint that creates a resource:

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

// Swagger documentation for POST /user
var _ = swagger.Swagger().Path("/user").
    Post(func(op openapi.Operation) {
        op.Summary("Create user").
            Description("This can only be done by the logged in user.").
            OperationID("createUser").
            Tag("user").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("Created user object").Required(true).SchemaFromDTO(&User{})
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation")
            })
    }).
    Doc()

// Handler function
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Implementation details
    c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/user", CreateUser)
}
```

## POST with Form Parameters

This example shows how to document a POST endpoint that accepts form data:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for POST /pet/{petId}
var _ = swagger.Swagger().Path("/pet/{petId}").
    Post(func(op openapi.Operation) {
        op.Summary("Updates a pet in the store with form data").
            OperationID("updatePetWithForm").
            Tag("pet").
            Consumes("application/x-www-form-urlencoded").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            PathParameter("petId", func(p openapi.Parameter) {
                p.Description("ID of pet that needs to be updated").Type("integer").Format("int64")
            }).
            FormParameter("name", func(p openapi.Parameter) {
                p.Description("Updated name of the pet").Required(false).Type("string")
            }).
            FormParameter("status", func(p openapi.Parameter) {
                p.Description("Updated status of the pet").Required(false).Type("string")
            }).
            Response(http.StatusMethodNotAllowed, func(r openapi.Response) {
                r.Description("Invalid input")
            })
    }).
    Doc()

// Handler function
func UpdatePetWithForm(c *gin.Context) {
    petID := c.Param("petId")
    name := c.PostForm("name")
    status := c.PostForm("status")
    // Implementation details
    c.JSON(http.StatusOK, gin.H{"message": "Pet updated", "id": petID, "name": name, "status": status})
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/pet/:petId", UpdatePetWithForm)
}
```

## POST with File Upload

Here's how to document a POST endpoint that handles file uploads:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// API Response DTO
type ApiResponse struct {
    Code    int32  `json:"code"`
    Type    string `json:"type"`
    Message string `json:"message"`
}

// Swagger documentation for POST /pet/{petId}/uploadImage
var _ = swagger.Swagger().Path("/pet/{petId}/uploadImage").
    Post(func(operation openapi.Operation) {
        operation.Summary("uploads an image").
            OperationID("uploadFile").
            Tag("pet").
            Consumes("multipart/form-data").
            Produce(mime.ApplicationJSON).
            PathParameter("petId", func(p openapi.Parameter) {
                p.Description("ID of pet to update").
                    MaxLength(64).
                    Type("integer").Format("int64")
            }).
            FormParameter("additionalMetadata", func(p openapi.Parameter) {
                p.Description("Additional data to pass to server").Required(false).Type("string")
            }).
            FormParameter("file", func(p openapi.Parameter) {
                p.Description("file to upload").Required(false).Type("file")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").SchemaFromDTO(&ApiResponse{})
            })
    }).
    Doc()

// Handler function
func UploadImage(c *gin.Context) {
    petID := c.Param("petId")
    metadata := c.PostForm("additionalMetadata")
    
    // Get uploaded file
    file, _ := c.FormFile("file")
    // Implementation details
    
    c.JSON(http.StatusOK, ApiResponse{
        Code: 200,
        Type: "success",
        Message: "Image uploaded for pet " + petID + " with metadata: " + metadata,
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/pet/:petId/uploadImage", UploadImage)
}
```

## POST with Array Input

This example demonstrates how to document a POST endpoint that accepts an array of objects:

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

// Swagger documentation for POST /user/createWithArray
var _ = swagger.Swagger().Path("/user/createWithArray").
    Post(func(op openapi.Operation) {
        op.Summary("Creates list of users with given input array").
            OperationID("createUsersWithArrayInput").
            Tag("user").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("List of user object").Required(true).
                    Schema(openapi_spec.SchemaEntity{
                        Type:  "array",
                        Items: &openapi_spec.SchemaEntity{Ref: "#/definitions/User"},
                    })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation")
            })
    }).
    Doc()

// Handler function
func CreateUsersWithArray(c *gin.Context) {
    var users []User
    if err := c.ShouldBindJSON(&users); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Implementation details
    c.JSON(http.StatusOK, gin.H{"message": "Users created", "count": len(users)})
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/user/createWithArray", CreateUsersWithArray)
}
```

## POST with Object Creation and Return

Here's how to document a POST endpoint that creates an object and returns it:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
    "time"
)

// Order DTO
type Order struct {
    ID       int64     `json:"id,omitempty"`
    PetID    int64     `json:"petId,omitempty"`
    Quantity int32     `json:"quantity,omitempty"`
    ShipDate time.Time `json:"shipDate,omitempty"`
    Status   string    `json:"status,omitempty"` // placed, approved, delivered
    Complete bool      `json:"complete,omitempty"`
}

// Swagger documentation for POST /store/order
var _ = swagger.Swagger().Path("/store/order").
    Post(func(op openapi.Operation) {
        op.Summary("Place an order for a pet").
            OperationID("placeOrder").
            Tag("store").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("order placed for purchasing the pet").Required(true).SchemaFromDTO(&Order{})
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").SchemaFromDTO(&Order{})
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid Order")
            })
    }).
    Doc()

// Handler function
func PlaceOrder(c *gin.Context) {
    var order Order
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order"})
        return
    }
    
    // Implementation details
    if order.ID == 0 {
        order.ID = 123 // Generate ID
    }
    order.ShipDate = time.Now()
    order.Status = "placed"
    
    c.JSON(http.StatusOK, order)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/store/order", PlaceOrder)
}
```