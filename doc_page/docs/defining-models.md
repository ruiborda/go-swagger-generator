---
sidebar_position: 14
title: Defining Models
---

# Defining Models

This guide shows how to define and use data models (schemas) with go-swagger-generator using practical examples.

## Basic Model Definition

Here's a simple example of defining a model using Go structs:

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
)

// Product defines a product in the catalog
type Product struct {
    ID          int64    `json:"id,omitempty"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    Price       float64  `json:"price"`
    Categories  []string `json:"categories,omitempty"`
    Tags        []string `json:"tags,omitempty"`
    InStock     bool     `json:"inStock"`
}

func main() {
    doc := swagger.Swagger()
    
    // Define the Product model in the Swagger documentation
    // This will make it available in the #/definitions section of the OpenAPI spec
    _, _ = doc.DefinitionFromDTO(&Product{})
    
    // Your API routes and other configurations here...
}
```

## Model with Nested Objects

This example demonstrates how to define models with nested object relationships:

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "time"
)

// Order represents a customer order
type Order struct {
    ID           int64       `json:"id,omitempty"`
    CustomerID   int64       `json:"customerId"`
    OrderDate    time.Time   `json:"orderDate"`
    ShippingInfo *Shipping   `json:"shippingInfo"`
    Items        []OrderItem `json:"items"`
    TotalPrice   float64     `json:"totalPrice"`
    Status       string      `json:"status,omitempty"` // pending, shipped, delivered
}

// Shipping represents shipping information
type Shipping struct {
    Address     string  `json:"address"`
    City        string  `json:"city"`
    State       string  `json:"state,omitempty"`
    Country     string  `json:"country"`
    PostalCode  string  `json:"postalCode"`
    PhoneNumber string  `json:"phoneNumber,omitempty"`
    TrackingNo  string  `json:"trackingNo,omitempty"`
    ShippingFee float64 `json:"shippingFee"`
}

// OrderItem represents a product in an order
type OrderItem struct {
    ProductID   int64   `json:"productId"`
    ProductName string  `json:"productName"`
    Quantity    int32   `json:"quantity"`
    UnitPrice   float64 `json:"unitPrice"`
    Subtotal    float64 `json:"subtotal"`
}

func main() {
    doc := swagger.Swagger()
    
    // Define all models in the Swagger documentation
    _, _ = doc.DefinitionFromDTO(&Order{})
    _, _ = doc.DefinitionFromDTO(&Shipping{})
    _, _ = doc.DefinitionFromDTO(&OrderItem{})
    
    // Your API routes and other configurations here...
}
```

## Models with Inheritance

Go doesn't have direct inheritance, but you can implement similar concepts in your Swagger documentation:

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
)

// Base model for all API responses
type BaseResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
}

// Data response that extends the base response
type DataResponse struct {
    BaseResponse
    Data interface{} `json:"data,omitempty"`
}

// Error response that extends the base response
type ErrorResponse struct {
    BaseResponse
    ErrorCode int    `json:"errorCode,omitempty"`
    Details   string `json:"details,omitempty"`
}

func main() {
    doc := swagger.Swagger()
    
    // Define base model
    _, _ = doc.DefinitionFromDTO(&BaseResponse{})
    
    // Define derived models manually to properly show inheritance
    doc.Definition("DataResponse", func(schema openapi.Schema) {
        schema.Type("object").
            Description("Data response that extends the base response").
            Property("success", func(prop openapi.Schema) {
                prop.Type("boolean")
                    .Description("Whether the operation was successful")
            }).
            Property("message", func(prop openapi.Schema) {
                prop.Type("string")
                    .Description("Response message")
            }).
            Property("data", func(prop openapi.Schema) {
                prop.Type("object")
                    .Description("The response data")
            })
    })

    doc.Definition("ErrorResponse", func(schema openapi.Schema) {
        schema.Type("object").
            Description("Error response that extends the base response").
            Property("success", func(prop openapi.Schema) {
                prop.Type("boolean")
                    .Description("Whether the operation was successful (always false)")
            }).
            Property("message", func(prop openapi.Schema) {
                prop.Type("string")
                    .Description("Error message")
            }).
            Property("errorCode", func(prop openapi.Schema) {
                prop.Type("integer")
                    .Format("int32")
                    .Description("Error code")
            }).
            Property("details", func(prop openapi.Schema) {
                prop.Type("string")
                    .Description("Detailed error information")
            })
    })
    
    // Your API routes and other configurations here...
}
```

## Enum Values in Models

Here's how to define models with enum values:

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
)

// Pet model with enum for status
type Pet struct {
    ID        int64    `json:"id,omitempty"`
    Name      string   `json:"name"`
    PhotoUrls []string `json:"photoUrls"`
    Status    string   `json:"status,omitempty"` // available, pending, sold
}

func main() {
    doc := swagger.Swagger()
    
    // Define the Pet model manually to specify enum values
    doc.Definition("Pet", func(schema openapi.Schema) {
        schema.Type("object").
            Property("id", func(prop openapi.Schema) {
                prop.Type("integer").Format("int64")
            }).
            Property("name", func(prop openapi.Schema) {
                prop.Type("string").Required(true)
            }).
            Property("photoUrls", func(prop openapi.Schema) {
                prop.Type("array").Required(true).
                    Items(func(items openapi.Schema) {
                        items.Type("string")
                    })
            }).
            Property("status", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Pet status in the store").
                    Enum("available", "pending", "sold")
            })
    })
    
    // Your API routes and other configurations here...
}
```

## Using Models in API Operations

Here's how to use your defined models in API operations:

```go
package main

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/middleware"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
)

// User model
type User struct {
    ID       int64  `json:"id,omitempty"`
    Username string `json:"username"`
    Email    string `json:"email"`
    FullName string `json:"fullName,omitempty"`
    Role     string `json:"role,omitempty"` // admin, user, guest
}

func main() {
    router := gin.Default()
    
    // Set up Swagger
    router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
        Enabled:  true,
        JSONPath: "/openapi.json",
        UIPath:   "/swagger",
    }))
    
    doc := swagger.Swagger()
    
    // Define the User model
    _, _ = doc.DefinitionFromDTO(&User{})
    
    // Define API operation that uses the User model
    doc.Path("/users/{userId}").
        Get(func(op openapi.Operation) {
            op.Summary("Get user by ID").
                Description("Returns a single user").
                OperationID("getUserById").
                Tag("users").
                Produces(mime.ApplicationJSON).
                PathParameter("userId", func(p openapi.Parameter) {
                    p.Description("ID of user to return").
                        Type("integer").Format("int64")
                }).
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("successful operation").
                        // Reference the User model defined earlier
                        SchemaFromDTO(&User{})
                }).
                Response(http.StatusNotFound, func(r openapi.Response) {
                    r.Description("User not found")
                })
        }).
        Post(func(op openapi.Operation) {
            op.Summary("Update user").
                Description("Update an existing user by ID").
                OperationID("updateUser").
                Tag("users").
                Consumes(string(mime.ApplicationJSON)).
                Produces(mime.ApplicationJSON).
                PathParameter("userId", func(p openapi.Parameter) {
                    p.Description("ID of user to update").
                        Type("integer").Format("int64")
                }).
                // Use the User model for the request body
                BodyParameter(func(p openapi.Parameter) {
                    p.Description("Updated user object").
                        Required(true).
                        SchemaFromDTO(&User{})
                }).
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("User updated successfully").
                        SchemaFromDTO(&User{})
                }).
                Response(http.StatusBadRequest, func(r openapi.Response) {
                    r.Description("Invalid user data")
                }).
                Response(http.StatusNotFound, func(r openapi.Response) {
                    r.Description("User not found")
                })
        }).
        Doc()
    
    // Handler functions
    router.GET("/users/:userId", func(c *gin.Context) {
        // Implementation omitted
    })
    
    router.POST("/users/:userId", func(c *gin.Context) {
        // Implementation omitted
    })
    
    // Run the server
    router.Run(":8080")
}
```

## Model with Validation Rules

Here's how to define models with validation rules:

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
)

func main() {
    doc := swagger.Swagger()
    
    // Define model with validation rules
    doc.Definition("SignupRequest", func(schema openapi.Schema) {
        schema.Type("object").
            Required("username", "email", "password"). // Required fields
            Property("username", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Username for the new account").
                    MinLength(3).   // Minimum length
                    MaxLength(50).  // Maximum length
                    Pattern("^[a-zA-Z0-9_-]+$") // Regex pattern for alphanumeric and some special chars
            }).
            Property("email", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Email address").
                    Format("email") // Email format validation
            }).
            Property("password", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Password").
                    Format("password"). // Password type
                    MinLength(8).       // Minimum password length
                    // Pattern for password complexity
                    Pattern("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]+$")
            }).
            Property("age", func(prop openapi.Schema) {
                prop.Type("integer").
                    Description("User's age").
                    Minimum(13, false). // Minimum age (inclusive)
                    Maximum(120, false) // Maximum age (inclusive)
            }).
            Property("website", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("User's website").
                    Format("uri") // URI format validation
            }).
            Property("interests", func(prop openapi.Schema) {
                prop.Type("array").
                    Description("User's interests").
                    MinItems(1).  // At least one interest required
                    MaxItems(10). // Maximum 10 interests
                    Items(func(items openapi.Schema) {
                        items.Type("string").
                            MinLength(2).
                            MaxLength(50)
                    })
            })
    })
    
    // Your API routes and other configurations here...
}
```

## Model with Additional Properties

Here's how to define a model that allows additional properties:

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
)

func main() {
    doc := swagger.Swagger()
    
    // Define a model with additional properties
    doc.Definition("Configuration", func(schema openapi.Schema) {
        schema.Type("object").
            Description("System configuration").
            Property("version", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Configuration version")
            }).
            Property("environment", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Deployment environment").
                    Enum("development", "testing", "staging", "production")
            }).
            // Allow additional properties of type string
            AdditionalProperties(&openapi_spec.SchemaEntity{
                Type: "string",
                Description: "Custom configuration value",
            })
    })
    
    // Define a generic map
    doc.Definition("StringMap", func(schema openapi.Schema) {
        schema.Type("object").
            Description("String map that can hold any string values").
            // Additional properties of any type
            AdditionalProperties(&openapi_spec.SchemaEntity{
                Type: "string",
            })
    })
    
    // Your API routes and other configurations here...
}
```