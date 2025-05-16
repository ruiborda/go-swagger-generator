---
sidebar_position: 12
title: Request Bodies
---

# Documenting Request Bodies

This guide shows how to document request bodies with go-swagger-generator using practical examples.

## Basic JSON Request Body

Here's a simple example of documenting an endpoint that accepts a JSON request body:

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
    ID        int64     `json:"id,omitempty"`
    Name      string    `json:"name"`
    PhotoUrls []string  `json:"photoUrls"`
    Status    string    `json:"status,omitempty"` // available, pending, sold
}

// Swagger documentation for POST /pet
var _ = swagger.Swagger().Path("/pet").
    Post(func(operation openapi.Operation) {
        operation.Summary("Add a new pet to the store").
            OperationID("addPet").
            Tag("pet").
            Consumes(string(mime.ApplicationJSON), string(mime.ApplicationXML)).
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("Pet object that needs to be added to the store").
                    Required(true).
                    SchemaFromDTO(&Pet{})  // Generate schema from struct
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Pet successfully added")
            }).
            Response(http.StatusMethodNotAllowed, func(r openapi.Response) {
                r.Description("Invalid input")
            })
    }).
    Doc()

// Handler function
func AddPet(c *gin.Context) {
    var pet Pet
    if err := c.ShouldBindJSON(&pet); err != nil {
        c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Invalid input"})
        return
    }
    // Implementation details
    c.JSON(http.StatusOK, pet)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/pet", AddPet)
}
```

## Request Body with Custom Schema

This example shows how to document a request body with a custom schema:

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

// Swagger documentation for POST /analytics/events
var _ = swagger.Swagger().Path("/analytics/events").
    Post(func(op openapi.Operation) {
        op.Summary("Track analytics events").
            Description("Send multiple analytics events in a single request").
            OperationID("trackEvents").
            Tag("analytics").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("Events to track").
                    Required(true).
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "userId": {
                                Type:        "string",
                                Description: "Unique identifier for the user",
                                Required:    []string{"userId"},
                            },
                            "sessionId": {
                                Type:        "string",
                                Description: "Current session identifier",
                            },
                            "events": {
                                Type:        "array",
                                Description: "List of events to track",
                                Items: &openapi_spec.SchemaEntity{
                                    Type: "object",
                                    Properties: map[string]*openapi_spec.SchemaEntity{
                                        "eventName": {
                                            Type:        "string",
                                            Description: "Name of the event",
                                        },
                                        "timestamp": {
                                            Type:        "string",
                                            Format:      "date-time",
                                            Description: "When the event occurred",
                                        },
                                        "properties": {
                                            Type:        "object",
                                            Description: "Additional properties for the event",
                                            AdditionalProperties: &openapi_spec.SchemaEntity{
                                                Type: "string",
                                            },
                                        },
                                    },
                                    Required: []string{"eventName", "timestamp"},
                                },
                            },
                        },
                        Required: []string{"userId", "events"},
                    })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Events tracked successfully").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "tracked": {Type: "integer", Format: "int32"},
                            "failed": {Type: "integer", Format: "int32"},
                        },
                    })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid event data")
            })
    }).
    Doc()

// Handler function
func TrackEvents(c *gin.Context) {
    var request struct {
        UserID    string `json:"userId"`
        SessionID string `json:"sessionId"`
        Events    []struct {
            EventName  string                 `json:"eventName"`
            Timestamp  string                 `json:"timestamp"`
            Properties map[string]interface{} `json:"properties"`
        } `json:"events"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
        return
    }

    // Implementation details
    c.JSON(http.StatusOK, gin.H{
        "tracked": len(request.Events),
        "failed": 0,
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/analytics/events", TrackEvents)
}
```

## Array in Request Body

Here's an example of documenting a request body that contains an array:

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
    Status   string `json:"status,omitempty"`
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
                p.Description("List of user objects").Required(true).
                    Schema(openapi_spec.SchemaEntity{
                        Type:  "array",
                        Items: &openapi_spec.SchemaEntity{Ref: "#/definitions/User"},
                    })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid user data")
            })
    }).
    Doc()

// Handler function
func CreateUsersWithArray(c *gin.Context) {
    var users []User
    if err := c.ShouldBindJSON(&users); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
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

## Request Body with Validation

This example demonstrates how to document a request body with validation constraints:

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

// Register DTO
type RegisterRequest struct {
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  string `json:"password"`
    FirstName string `json:"firstName,omitempty"`
    LastName  string `json:"lastName,omitempty"`
    Age       int    `json:"age,omitempty"`
}

// Swagger documentation for POST /auth/register
var _ = swagger.Swagger().Path("/auth/register").
    Post(func(op openapi.Operation) {
        op.Summary("Register a new user").
            Description("Create a new user account").
            OperationID("registerUser").
            Tag("auth").
            Consumes(string(mime.ApplicationJSON)).
            Produces(mime.ApplicationJSON).
            BodyParameter(func(p openapi.Parameter) {
                p.Description("User registration details").
                    Required(true).
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "username": {
                                Type:        "string",
                                Description: "Username for the new account",
                                MinLength:   4,
                                MaxLength:   20,
                                Pattern:     "^[a-zA-Z0-9_]+$", // Alphanumeric and underscore only
                            },
                            "email": {
                                Type:        "string",
                                Description: "Email address",
                                Format:      "email",
                            },
                            "password": {
                                Type:        "string",
                                Description: "Password (min 8 chars, must include uppercase, lowercase, and number)",
                                MinLength:   8,
                                Pattern:     "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).+$",
                            },
                            "firstName": {
                                Type:        "string",
                                Description: "First name",
                            },
                            "lastName": {
                                Type:        "string",
                                Description: "Last name",
                            },
                            "age": {
                                Type:        "integer",
                                Description: "Age in years",
                                Format:      "int32",
                                Minimum:     13, // Minimum age
                                Maximum:     120, // Maximum age
                            },
                        },
                        Required: []string{"username", "email", "password"},
                    })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("User registered successfully").
                    Schema(openapi_spec.SchemaEntity{
                        Type: "object",
                        Properties: map[string]*openapi_spec.SchemaEntity{
                            "userId": {Type: "integer", Format: "int64"},
                            "username": {Type: "string"},
                        },
                    })
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid registration data")
            }).
            Response(http.StatusConflict, func(r openapi.Response) {
                r.Description("Username or email already exists")
            })
    }).
    Doc()

// Handler function
func RegisterUser(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration data"})
        return
    }

    // Validation and implementation details
    
    c.JSON(http.StatusOK, gin.H{
        "userId": 12345,
        "username": req.Username,
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/auth/register", RegisterUser)
}
```

## Form Data Request Body

Here's how to document an endpoint that accepts form data:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Swagger documentation for POST /contact
var _ = swagger.Swagger().Path("/contact").
    Post(func(op openapi.Operation) {
        op.Summary("Send contact message").
            Description("Submit a contact form message").
            OperationID("sendContactMessage").
            Tag("contact").
            Consumes("application/x-www-form-urlencoded").
            Produces(mime.ApplicationJSON).
            FormParameter("name", func(p openapi.Parameter) {
                p.Description("Your full name").
                    Type("string").
                    Required(true).
                    MinLength(2).
                    MaxLength(100)
            }).
            FormParameter("email", func(p openapi.Parameter) {
                p.Description("Your email address").
                    Type("string").
                    Required(true).
                    Format("email")
            }).
            FormParameter("subject", func(p openapi.Parameter) {
                p.Description("Message subject").
                    Type("string").
                    Required(true).
                    MinLength(5).
                    MaxLength(200)
            }).
            FormParameter("message", func(p openapi.Parameter) {
                p.Description("Message content").
                    Type("string").
                    Required(true).
                    MinLength(10).
                    MaxLength(2000)
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Message sent successfully")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid form data")
            })
    }).
    Doc()

// Handler function
func SendContactMessage(c *gin.Context) {
    name := c.PostForm("name")
    email := c.PostForm("email")
    subject := c.PostForm("subject")
    message := c.PostForm("message")
    
    // Validation and implementation details
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Thank you for your message!",
    })
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.POST("/contact", SendContactMessage)
}
```