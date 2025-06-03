---
sidebar_position: 14
title: Defining Models (Schemas) (v2)
---

# Defining Models (Schemas) (v2)

This guide shows how to define and use data models (schemas) with Go Swagger Generator v2 for OpenAPI 3.0 documentation, using practical examples. In OpenAPI 3.0, schemas are typically defined under `#/components/schemas/`.

## Basic Model Definition

Here's a simple example of defining a model using Go structs. `doc.SchemaFromDTO()` or `doc.ComponentSchemaFromDTO()` can be used to register a Go struct as a schema component.

```go
package main

import (
    // "github.com/gin-gonic/gin" // if building a full app
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

// Product defines a product in the catalog
type Product struct {
    ID          int64    `json:"id,omitempty" yaml:"id,omitempty"`
    Name        string   `json:"name" yaml:"name"`
    Description string   `json:"description,omitempty" yaml:"description,omitempty"`
    Price       float64  `json:"price" yaml:"price"`
    Categories  []string `json:"categories,omitempty" yaml:"categories,omitempty"`
    Tags        []string `json:"tags,omitempty" yaml:"tags,omitempty"`
    InStock     bool     `json:"inStock" yaml:"inStock"`
}

func main() {
    doc := swagger.Swagger() // Initializes OpenAPI 3.0 builder by default
    
    // Define the Product model in the Swagger documentation as a component schema.
    // This will make it available as #/components/schemas/Product in the OpenAPI spec.
    _, _ = doc.ComponentSchemaFromDTO(&Product{}) 
    // Alternatively, doc.SchemaFromDTO(&Product{}) also works and registers it if not present.
    
    // Your API routes and other configurations here...
    // Example: doc.Path("/products").Get(func(op openapi.Operation) { ... })
}
```

## Model with Nested Objects

This example demonstrates how to define models with nested object relationships. Each struct will be registered as a separate schema component.

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi" // Assuming this path for openapi types
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "time"
)

// Order represents a customer order
type Order struct {
    ID           int64       `json:"id,omitempty" yaml:"id,omitempty"`
    CustomerID   int64       `json:"customerId" yaml:"customerId"`
    OrderDate    time.Time   `json:"orderDate" yaml:"orderDate"`
    ShippingInfo *Shipping   `json:"shippingInfo" yaml:"shippingInfo"`
    Items        []OrderItem `json:"items" yaml:"items"`
    TotalPrice   float64     `json:"totalPrice" yaml:"totalPrice"`
    Status       string      `json:"status,omitempty" yaml:"status,omitempty"` // e.g., pending, shipped, delivered
}

// Shipping represents shipping information
type Shipping struct {
    Address     string  `json:"address" yaml:"address"`
    City        string  `json:"city" yaml:"city"`
    State       string  `json:"state,omitempty" yaml:"state,omitempty"`
    Country     string  `json:"country" yaml:"country"`
    PostalCode  string  `json:"postalCode" yaml:"postalCode"`
    PhoneNumber string  `json:"phoneNumber,omitempty" yaml:"phoneNumber,omitempty"`
    TrackingNo  string  `json:"trackingNo,omitempty" yaml:"trackingNo,omitempty"`
    ShippingFee float64 `json:"shippingFee" yaml:"shippingFee"`
}

// OrderItem represents a product in an order
type OrderItem struct {
    ProductID   int64   `json:"productId" yaml:"productId"`
    ProductName string  `json:"productName" yaml:"productName"`
    Quantity    int32   `json:"quantity" yaml:"quantity"`
    UnitPrice   float64 `json:"unitPrice" yaml:"unitPrice"`
    Subtotal    float64 `json:"subtotal" yaml:"subtotal"`
}

func main() {
    doc := swagger.Swagger()
    
    // Define all models in #/components/schemas/
    _, _ = doc.ComponentSchemaFromDTO(&Order{})
    _, _ = doc.ComponentSchemaFromDTO(&Shipping{})
    _, _ = doc.ComponentSchemaFromDTO(&OrderItem{})
    
    // Your API routes and other configurations here...
}
```

## Models with Composition (Simulating Inheritance)

Go uses composition (embedding structs) rather than inheritance. OpenAPI 3.0 supports `allOf` for schema composition.
When you embed a struct, `ComponentSchemaFromDTO` will typically flatten properties or you might need to define it manually for explicit `allOf`.
For explicit `allOf` (manual definition):

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

// Base model for all API responses
type BaseResponse struct {
    Success bool   `json:"success" yaml:"success"`
    Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

// Data response that extends the base response
type DataResponse struct {
    BaseResponse // Embedded
    Data interface{} `json:"data,omitempty" yaml:"data,omitempty"`
}

func main() {
    doc := swagger.Swagger()
    
    // Register BaseResponse first
    _, _ = doc.ComponentSchemaFromDTO(&BaseResponse{})

    // Manually define DataResponse to use allOf for clear composition
    doc.ComponentSchema("DataResponse", func(schema openapi.Schema) {
        schema.Description("Data response that composes BaseResponse and adds a data field.").
            AllOf( // Use AllOf to combine schemas
                openapi.S().Ref("#/components/schemas/BaseResponse"), // Reference to BaseResponse schema
                openapi.S().Type("object").Property("data", func(prop openapi.Schema) {
                    prop.Type("object").Description("The response data")
                    // Could also be more specific: prop.Ref("#/components/schemas/SpecificDataModel")
                }),
            )
    })
    
    // If you register DataResponse directly using ComponentSchemaFromDTO(&DataResponse{}), 
    // it will likely inline the properties of BaseResponse.
    // For explicit allOf, manual definition is better.
}
```

## Enum Values in Models

Here's how to define models with enum values, typically done by manually defining or extending the schema property.

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

// Pet model with enum for status
type Pet struct {
    ID        int64    `json:"id,omitempty" yaml:"id,omitempty"`
    Name      string   `json:"name" yaml:"name"`
    PhotoUrls []string `json:"photoUrls" yaml:"photoUrls"`
    Status    string   `json:"status,omitempty" yaml:"status,omitempty"` // Intended values: available, pending, sold
}

func main() {
    doc := swagger.Swagger()
    
    // Define the Pet model component schema, then customize the 'status' property
    sPet, _ := doc.ComponentSchemaFromDTO(&Pet{})
    sPet.Property("status", func(prop openapi.Schema) {
        prop.Description("Pet status in the store").
             Enum("available", "pending", "sold") // Add Enum constraint
    })
    // Alternatively, define the entire schema manually if more control is needed:
    /*
    doc.ComponentSchema("Pet", func(schema openapi.Schema) {
        schema.Type("object").
            Property("id", func(prop openapi.Schema) {
                prop.Type("integer").Format("int64")
            }).
            Property("name", func(prop openapi.Schema) {
                prop.Type("string").Required("name") // Use Required on the parent schema object
            }).
            Property("photoUrls", func(prop openapi.Schema) {
                prop.Type("array").Required("photoUrls").
                    Items(func(items openapi.Schema) {
                        items.Type("string")
                    })
            }).
            Property("status", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Pet status in the store").
                    Enum("available", "pending", "sold")
            }).Required("name", "photoUrls") // List required properties at the object level
    })
    */
}
```

## Using Models in API Operations (Request Bodies and Responses)

Here's how to use your defined models (schemas) in API operations for request bodies and responses.

```go
package main

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/middleware"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

// User model
type User struct {
    ID       int64  `json:"id,omitempty" yaml:"id,omitempty"`
    Username string `json:"username" yaml:"username"`
    Email    string `json:"email" yaml:"email"`
    FullName string `json:"fullName,omitempty" yaml:"fullName,omitempty"`
    Role     string `json:"role,omitempty" yaml:"role,omitempty"` // e.g., admin, user, guest
}

func main() {
    router := gin.Default()
    
    doc := swagger.Swagger()
    _, _ = doc.ComponentSchemaFromDTO(&User{}) // Ensure User schema is registered

    // Set up Swagger UI middleware
    router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
        Enabled:  true,
        JSONPath: "/openapi.json",
        UIPath:   "/", // Serve Swagger UI at root
        Title:    "User API v2",
    }))
    doc.Server("http://localhost:8080/v2", func(s openapi.Server) { s.Description("Dev Server"); })
    
    // Define API operation that uses the User model
    doc.Path("/users/{userId}"). // Path relative to server URL
        Get(func(op openapi.Operation) {
            op.Summary("Get user by ID").
                Description("Returns a single user").
                OperationID("getUserById").
                Tag("users").
                PathParameter("userId", func(p openapi.Parameter) {
                    p.Description("ID of user to return").Required(true).
                        Schema(func(s openapi.Schema) { s.Type("integer").Format("int64") })
                }).
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("successful operation").
                        Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                            mt.SchemaFromDTO(&User{}) // Reference User schema for response
                        })
                }).
                Response(http.StatusNotFound, func(r openapi.Response) {
                    r.Description("User not found")
                })
        }).
        Post(func(op openapi.Operation) {
            op.Summary("Create or Update user"). // POST can be for create, PUT for update typically
                Description("Creates a new user or updates an existing user by ID.").
                OperationID("createOrUpdateUser").
                Tag("users").
                PathParameter("userId", func(p openapi.Parameter) {
                    p.Description("ID of user to create/update").Required(true).
                        Schema(func(s openapi.Schema) { s.Type("integer").Format("int64") })
                }).
                RequestBody(func(rb openapi.RequestBody) { // Define Request Body for POST
                    rb.Description("User object for creation/update").
                        Required(true).
                        Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                            mt.SchemaFromDTO(&User{}) // Use User schema for request body
                        })
                }).
                Response(http.StatusOK, func(r openapi.Response) {
                    r.Description("User updated successfully").
                        Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                            mt.SchemaFromDTO(&User{})
                        })
                }).
                Response(http.StatusCreated, func(r openapi.Response) {
                    r.Description("User created successfully").
                        Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                            mt.SchemaFromDTO(&User{})
                        })
                }).
                Response(http.StatusBadRequest, func(r openapi.Response) {
                    r.Description("Invalid user data supplied")
                })
        }).
        Doc()
    
    // Handler functions (omitted for brevity)
    router.GET("/v2/users/:userId", func(c *gin.Context) { /* ... */ })
    router.POST("/v2/users/:userId", func(c *gin.Context) { /* ... */ })
    
    router.Run(":8080")
}
```

## Model with Validation Rules

Define models with validation rules using OpenAPI 3.0 keywords within the schema definition.

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

func main() {
    doc := swagger.Swagger()
    
    // Define model with validation rules as a component schema
    doc.ComponentSchema("SignupRequest", func(schema openapi.Schema) {
        schema.Type("object").
            Required("username", "email", "password"). // Required properties at object level
            Property("username", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Username for the new account").
                    MinLength(3).   // Minimum length
                    MaxLength(50).  // Maximum length
                    Pattern("^[a-zA-Z0-9_-]+$") // Regex pattern
            }).
            Property("email", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Email address").
                    Format("email") // Email format validation
            }).
            Property("password", func(prop openapi.Schema) {
                prop.Type("string").
                    Description("Password").
                    Format("password"). // Suggests input type, not strong validation
                    MinLength(8).       // Minimum password length
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
                    Items(func(itemsSchema openapi.Schema) {
                        itemsSchema.Type("string").
                            MinLength(2).
                            MaxLength(50)
                    })
            })
    })
    
    // Your API routes and other configurations here...
}
```

## Model with Additional Properties

Here's how to define a model that allows additional properties (like a map or dictionary).

```go
package main

import (
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

func main() {
    doc := swagger.Swagger()
    
    // Define a model with additional properties
    doc.ComponentSchema("Configuration", func(schema openapi.Schema) {
        schema.Type("object").
            Description("System configuration with potentially dynamic keys").
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
            AdditionalProperties(true, func(addPropsSchema openapi.Schema) {
                addPropsSchema.Type("string").Description("Custom configuration value")
            })
            // To allow any type for additional properties (OpenAPI 3.0.x):
            // AdditionalProperties(true, openapi.S()) // or AdditionalProperties(true, nil) for 'true'
            // For OpenAPI 3.1+, AdditionalProperties can be a full schema object directly.
    })
    
    // Your API routes and other configurations here...
}
```