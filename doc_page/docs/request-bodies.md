---
sidebar_position: 12
title: Request Bodies (v2)
---

# Documenting Request Bodies (v2)

This guide shows how to document request bodies with Go-Swagger-Generator v2 for OpenAPI 3.0. Request bodies are typically used with POST, PUT, and PATCH operations to send data to the server.

## Basic JSON Request Body

Here's an example of documenting an endpoint (e.g., POST for creating a resource) that accepts a JSON request body.

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
    ID        int64    `json:"id,omitempty" yaml:"id,omitempty"`
    Name      string   `json:"name" yaml:"name"`
    PhotoUrls []string `json:"photoUrls" yaml:"photoUrls"`
    Status    string   `json:"status,omitempty" yaml:"status,omitempty"` // e.g., available, pending, sold
}
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&Pet{})

// Swagger documentation for POST /pets
var _ = swagger.Swagger().Path("/pets"). // Path relative to server URL
    Post(func(op openapi.Operation) {
        op.Summary("Add a new pet to the store").
            OperationID("addPetV2").
            Tag("Pet Operations").
            RequestBody(func(rb openapi.RequestBody) { // Define the request body
                rb.Description("Pet object that needs to be added to the store").
                  Required(true).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) { // Specify content type
                      mt.SchemaFromDTO(&Pet{}) // Generate schema from Pet DTO
                  })
                // Optionally, support other content types like XML:
                // rb.Content(mime.ApplicationXML, func(mt openapi.MediaType) {
                //    mt.SchemaFromDTO(&Pet{})
                // })
            }).
            Response(http.StatusCreated, func(r openapi.Response) {
                r.Description("Pet created successfully").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.SchemaFromDTO(&Pet{}) // Return the created pet
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid input provided")
            })
    }).
    Doc()

// Handler function (example)
func AddPet(c *gin.Context) {
    var pet Pet
    if err := c.ShouldBindJSON(&pet); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    pet.ID = 123 // Simulate ID generation
    // Save pet to database...
    c.JSON(http.StatusCreated, pet)
}
```

## Request Body with Custom Inline Schema

If you don't have a DTO or need a very specific structure for a single request, you can define the schema inline.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// EventData structure for the example (could also be a DTO)
type EventData struct {
    EventName  string                 `json:"eventName"`
    Timestamp  string                 `json:"timestamp"`
    Properties map[string]interface{} `json:"properties,omitempty"`
}

// Swagger documentation for POST /analytics/events
var _ = swagger.Swagger().Path("/analytics/events").
    Post(func(op openapi.Operation) {
        op.Summary("Track analytics events").
            Description("Send one or more analytics events in a single request.").
            OperationID("trackEventsV2").
            Tag("Analytics Operations").
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("Analytics events payload").
                  Required(true).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) { // Define schema inline
                          s.Type("object").
                            Required("userId", "events").
                            Property("userId", func(ps openapi.Schema) {
                                ps.Type("string").Description("Unique identifier for the user")
                            }).
                            Property("sessionId", func(ps openapi.Schema) {
                                ps.Type("string").Description("Current session identifier (optional)")
                            }).
                            Property("events", func(ps openapi.Schema) {
                                ps.Type("array").Description("List of events to track").MinItems(1).
                                  Items(func(itemSchema openapi.Schema) { // Schema for each event in the array
                                      itemSchema.Type("object").Required("eventName", "timestamp").
                                        Property("eventName", func(prop openapi.Schema){ prop.Type("string") }).
                                        Property("timestamp", func(prop openapi.Schema){ prop.Type("string").Format("date-time") }).
                                        Property("properties", func(prop openapi.Schema){
                                            prop.Type("object").Description("Additional key-value pairs for the event").
                                               AdditionalProperties(true, openapi.S().Type("string")) // Example: all additional props are strings
                                        })
                                  })
                            })
                      })
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Events tracked successfully").
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) { // Response schema
                          s.Type("object").
                            Property("trackedCount", func(ps openapi.Schema){ ps.Type("integer") }).
                            Property("failedCount", func(ps openapi.Schema){ ps.Type("integer") })
                      })
                  })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid event data provided")
            })
    }).
    Doc()

// Handler function (example)
func TrackEvents(c *gin.Context) {
    // var requestPayload struct { ... } // Define a struct to bind JSON
    // if err := c.ShouldBindJSON(&requestPayload); ...
    c.JSON(http.StatusOK, gin.H{"trackedCount": 1, "failedCount": 0})
}
```

## Array as Request Body

If the entire request body is an array of objects (e.g., bulk creation).

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// User DTO (defined in earlier examples)
// type User struct { ... }
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&User{})

// Swagger documentation for POST /users/bulkCreate
var _ = swagger.Swagger().Path("/users/bulkCreate").
    Post(func(op openapi.Operation) {
        op.Summary("Creates multiple users from a list").
            OperationID("createUsersWithListV2").
            Tag("User Operations").
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("List of user objects to create").Required(true).
                  Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                      // For an array request body, use SchemaFromDTO with a pointer to a slice of DTO pointers
                      mt.SchemaFromDTO(&[]*User{}) 
                  })
            }).
            Response(http.StatusCreated, func(r openapi.Response) {
                r.Description("Users created successfully (or report on status)")
                // Response could be a summary, e.g., number created/failed
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid user data in the list")
            })
    }).
    Doc()

// Handler function (example)
func CreateUsersWithList(c *gin.Context) {
    var users []*User
    if err := c.ShouldBindJSON(&users); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
        return
    }
    // Implementation to create multiple users...
    c.JSON(http.StatusCreated, gin.H{"message": "Users processed", "count": len(users)})
}
```

## Form Data Request Body (`application/x-www-form-urlencoded` or `multipart/form-data`)

For form data, the `Content` type is typically `application/x-www-form-urlencoded` or `multipart/form-data`. The schema describes the form fields.

**`application/x-www-form-urlencoded` Example:**

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for POST /submit-feedback
var _ = swagger.Swagger().Path("/submit-feedback").
    Post(func(op openapi.Operation) {
        op.Summary("Submit user feedback via form").
            OperationID("submitFeedbackV2").
            Tag("Feedback Operations").
            RequestBody(func(rb openapi.RequestBody) {
                rb.Description("Feedback form data").Required(true).
                  Content(mime.ApplicationFormUrlEncoded, func(mt openapi.MediaType) {
                      mt.Schema(func(s openapi.Schema) {
                          s.Type("object").
                            Required("email", "message"). // Specify required form fields
                            Property("name", func(ps openapi.Schema){ ps.Type("string").Description("Your name (optional)") }).
                            Property("email", func(ps openapi.Schema){ ps.Type("string").Format("email").Description("Your email address") }).
                            Property("subject", func(ps openapi.Schema){ ps.Type("string").Description("Subject of feedback").Default("General Feedback") }).
                            Property("message", func(ps openapi.Schema){ ps.Type("string").Description("Your feedback message").MinLength(10) })
                      })
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Feedback submitted successfully")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid form data provided")
            })
    }).
    Doc()

// Handler function (example)
func SubmitFeedback(c *gin.Context) {
    // name := c.PostForm("name")
    // email := c.PostForm("email")
    // ... process form data ...
    c.JSON(http.StatusOK, gin.H{"message": "Feedback received!"})
}
```
For `multipart/form-data` (e.g., including file uploads), see the [POST Endpoints (v2)](./post-endpoints.md#post-with-file-upload-multipartform-data) guide.

**Key elements for `RequestBody`:**
*   `Description(string)`: A description of the request body.
*   `Required(bool)`: Whether the request body is mandatory.
*   `Content(mimeType string, func(mt openapi.MediaType))`: Defines one or more media types the endpoint consumes.
    *   Inside `Content`, `mt.SchemaFromDTO(&YourType{})` or `mt.Schema(func(s openapi.Schema){...})` defines the structure of the data for that media type.