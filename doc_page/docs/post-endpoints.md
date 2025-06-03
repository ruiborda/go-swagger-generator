---
sidebar_position: 7
title: POST Endpoints (v2)
---

# Documenting POST Endpoints (v2)

This guide shows how to document POST endpoints with Go-Swagger-Generator v2 for OpenAPI 3.0. POST requests are typically used to create new resources or trigger actions.

## Basic POST Endpoint (Resource Creation)

Here's an example of documenting a POST endpoint that creates a new resource using a JSON request body.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// User DTO
type User struct {
    ID        int64  `json:"id,omitempty" yaml:"id,omitempty"`
    Username  string `json:"username" yaml:"username"`
    FirstName string `json:"firstName,omitempty" yaml:"firstName,omitempty"`
    LastName  string `json:"lastName,omitempty" yaml:"lastName,omitempty"`
    Email     string `json:"email" yaml:"email"`
    Password  string `json:"password,omitempty" yaml:"password,omitempty"` // Typically not returned
    Phone     string `json:"phone,omitempty" yaml:"phone,omitempty"`
}

// Ensure User DTO is registered: _, _ = swagger.Swagger().ComponentSchemaFromDTO(&User{})

// Swagger documentation for POST /users
var _ = swagger.Swagger().Path("/users"). // Path relative to server URL
    Post(func(op openapi.Operation) {
        op.Summary("Create a new user").
            Description("Creates a new user account.").
            OperationID("createUserV2").
            Tag("User Operations").
            RequestBody(func(rb openapi.RequestBody) { // Define the request body
                rb.Description("User object to be created").
                  Required(true).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&User{}) // Use User DTO for request body schema
                  })
                // Optionally add other consumable types e.g. XML
                // rb.Content(mime.ApplicationXML, func(mt openapi.MediaType) {
                //    mt.SchemaFromDTO(&User{})
                // })
            }).
            Response(http.StatusCreated, func(r openapi.Response) { // 201 Created is common for successful POST
                r.Description("User created successfully").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&User{}) // Return the created user (omitting password)
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid user data supplied")
            })
    }).
    Doc()

// Handler function (example)
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user.ID = 123 // Simulate ID generation
    // Save user to database...
    user.Password = "" // Clear password before returning
    c.JSON(http.StatusCreated, user)
}
```

## POST with Form Parameters (`application/x-www-form-urlencoded`)

This example shows how to document a POST endpoint that accepts URL-encoded form data.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    // "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// Swagger documentation for POST /pets/{petId}
var _ = swagger.Swagger().Path("/pets/{petId}"). // Path relative to server URL
    Post(func(op openapi.Operation) {
        op.Summary("Updates a pet in the store with form data").
            OperationID("updatePetWithFormV2").
            Tag("Pet Operations").
            PathParameter("petId", func(p openapi.Parameter) {
                p.Description("ID of pet that needs to be updated").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("integer").Format("int64") })
            }).
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("Pet name and status to update").
                  Required(true).
                  Content("application/x-www-form-urlencoded", func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) {
                          s.Type("object").
                            Property("name", func(propSchema openapi.Schema) {
                                propSchema.Type("string").Description("Updated name of the pet")
                            }).
                            Property("status", func(propSchema openapi.Schema) {
                                propSchema.Type("string").Description("Updated status of the pet")
                            })
                          // Form parameters are typically not individually marked as 'required' here;
                          // the requirement is on the properties within the schema if applicable, 
                          // or the entire request body can be required.
                      })
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Pet updated successfully")
            }).
            Response(http.StatusMethodNotAllowed, func(r openapi.Response) { // Or 400 Bad Request
                r.Description("Invalid input")
            })
    }).
    Doc()

// Handler function (example)
func UpdatePetWithForm(c *gin.Context) {
    // petID := c.Param("petId")
    // name := c.PostForm("name")
    // status := c.PostForm("status")
    // Implementation...
    c.JSON(http.StatusOK, gin.H{"message": "Pet updated"})
}
```

## POST with File Upload (`multipart/form-data`)

Here's how to document a POST endpoint that handles file uploads along with other form data.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// ApiResponse DTO for file upload response
type ApiResponse struct {
    Code    int32  `json:"code,omitempty" yaml:"code,omitempty"`
    Type    string `json:"type,omitempty" yaml:"type,omitempty"`
    Message string `json:"message,omitempty" yaml:"message,omitempty"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&ApiResponse{})

// Swagger documentation for POST /pets/{petId}/uploadImage
var _ = swagger.Swagger().Path("/pets/{petId}/uploadImage").
    Post(func(op openapi.Operation) {
        op.Summary("Uploads an image for a pet").
            OperationID("uploadPetImageV2").
            Tag("Pet Operations").
            PathParameter("petId", func(p openapi.Parameter) {
                p.Description("ID of pet to update").Required(true).
                  Schema(func(s openapi.Schema) { s.Type("integer").Format("int64") })
            }).
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("Image file and additional metadata").
                  Required(true).
                  Content("multipart/form-data", func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) {
                          s.Type("object").
                            Property("additionalMetadata", func(propSchema openapi.Schema) {
                                propSchema.Type("string").Description("Additional data to pass to server")
                            }).
                            Property("file", func(propSchema openapi.Schema) {
                                propSchema.Type("string").Format("binary").Description("Image file to upload")
                            })
                          // To mark parts as required inside multipart form:
                          // s.Required("file") 
                      })
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Image uploaded successfully").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&ApiResponse{})
                  })
            })
    }).
    Doc()

// Handler function (example)
func UploadImage(c *gin.Context) {
    // petID := c.Param("petId")
    // metadata := c.PostForm("additionalMetadata")
    // file, _ := c.FormFile("file")
    // Implementation...
    c.JSON(http.StatusOK, ApiResponse{Code: 200, Type: "success", Message: "Image uploaded"})
}
```

## POST with Array Input

This example demonstrates documenting a POST endpoint that accepts an array of objects in the request body.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// User DTO (defined earlier)
// type User struct { ... }
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&User{})

// Swagger documentation for POST /users/createWithArray
var _ = swagger.Swagger().Path("/users/createWithArray").
    Post(func(op openapi.Operation) {
        op.Summary("Creates a list of users from an array").
            OperationID("createUsersWithArrayInputV2").
            Tag("User Operations").
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("Array of user objects to create").Required(true).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      // For an array request body, use SchemaFromDTO with a pointer to a slice of pointers
                      mt.SchemaFromDTO(&[]*User{}) 
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) { // Or 201 Created
                r.Description("Successfully created users (or partial success report)")
                // Response could be a summary or array of created users
            })
    }).
    Doc()

// Handler function (example)
func CreateUsersWithArray(c *gin.Context) {
    var users []*User
    if err := c.ShouldBindJSON(&users); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Implementation to create multiple users...
    c.JSON(http.StatusOK, gin.H{"message": "Users processed", "count": len(users)})
}
```

## POST for Actions (No Specific Resource Creation)

A POST request can also be used to trigger an action that doesn't necessarily create a resource identifiable by the request URL.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    // "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
    "net/http"
)

// ActionRequest DTO for the action
type ActionRequest struct {
    ActionName string                 `json:"actionName" yaml:"actionName"`
    Parameters map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&ActionRequest{})

// Swagger documentation for POST /system/actions
var _ = swagger.Swagger().Path("/system/actions").
    Post(func(op openapi.Operation) {
        op.Summary("Trigger a system action").
            OperationID("triggerSystemActionV2").
            Tag("System Operations").
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("Action to perform").Required(true).
                  Content("application/json", func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&ActionRequest{})
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Action triggered successfully").
                  Content("application/json", func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) { // Custom success response schema
                          s.Type("object").
                            Property("status", func(ps openapi.Schema){ ps.Type("string") }).
                            Property("details", func(ps openapi.Schema){ ps.Type("string") })
                      })
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid action request")
            })
    }).
    Doc()

// Handler function (example)
func TriggerSystemAction(c *gin.Context) {
    var req ActionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action request"})
        return
    }
    // Perform action based on req.ActionName and req.Parameters...
    c.JSON(http.StatusOK, gin.H{"status": "success", "details": "Action '" + req.ActionName + "' processed."})
}
```