---
sidebar_position: 10
title: Path Parameters
---

# Documenting Path Parameters

This guide shows how to document path parameters with go-swagger-generator using practical examples.

## Basic Path Parameter

Here's a simple example of documenting an endpoint with a path parameter:

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
    Status    string `json:"status,omitempty"` // available, pending, sold
}

// Swagger documentation for GET /pet/{petId}
var _ = swagger.Swagger().Path("/pet/{petId}").
    Get(func(op openapi.Operation) {
        op.Summary("Find pet by ID").
            Description("Returns a single pet").
            OperationID("getPetById").
            Tag("pet").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            PathParameter("petId", func(p openapi.Parameter) {
                p.Description("ID of pet to return").
                    Type("integer").
                    Format("int64")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").SchemaFromDTO(&Pet{})
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid ID supplied")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Pet not found")
            })
    }).
    Doc()

// Handler function
func GetPetByID(c *gin.Context) {
    petID := c.Param("petId")
    // Implementation details
    pet := Pet{ID: 1, Name: "Doggie", Status: "available"}
    c.JSON(http.StatusOK, pet)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/pet/:petId", GetPetByID)
}
```

## Path Parameter with Validation

This example shows how to document a path parameter with validation constraints:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Order DTO
type Order struct {
    ID       int64  `json:"id,omitempty"`
    PetID    int64  `json:"petId,omitempty"`
    Quantity int32  `json:"quantity,omitempty"`
    Status   string `json:"status,omitempty"` // placed, approved, delivered
}

// Swagger documentation for GET /store/order/{orderId}
var _ = swagger.Swagger().Path("/store/order/{orderId}").
    Get(func(op openapi.Operation) {
        op.Summary("Find purchase order by ID").
            Description("For valid response try integer IDs with value >= 1 and <= 10. Other values will generate exceptions").
            OperationID("getOrderById").
            Tag("store").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            PathParameter("orderId", func(p openapi.Parameter) {
                p.Description("ID of pet that needs to be fetched").
                    Type("integer").Format("int64").
                    Minimum(1, false).  // Minimum value (exclusive=false means inclusive)
                    Maximum(10, false)  // Maximum value (exclusive=false means inclusive)
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").SchemaFromDTO(&Order{})
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
func GetOrderByID(c *gin.Context) {
    orderID := c.Param("orderId")
    // Implementation details
    order := Order{ID: 1, PetID: 1, Quantity: 5, Status: "placed"}
    c.JSON(http.StatusOK, order)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/store/order/:orderId", GetOrderByID)
}
```

## String Path Parameter

Here's an example of documenting a path parameter that's a string with pattern validation:

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
    Email    string `json:"email,omitempty"`
}

// Swagger documentation for GET /user/{username}
var _ = swagger.Swagger().Path("/user/{username}").
    Get(func(op openapi.Operation) {
        op.Summary("Get user by user name").
            OperationID("getUserByName").
            Tag("user").
            Produces(mime.ApplicationJSON, mime.ApplicationXML).
            PathParameter("username", func(p openapi.Parameter) {
                p.Description("The name that needs to be fetched").
                    Type("string").
                    Pattern("^[a-zA-Z0-9]+$").  // Only alphanumeric characters
                    MinLength(3).               // Minimum length
                    MaxLength(50)               // Maximum length
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").SchemaFromDTO(&User{})
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
func GetUserByName(c *gin.Context) {
    username := c.Param("username")
    // Implementation details
    user := User{ID: 1, Username: username, Email: "user@example.com"}
    c.JSON(http.StatusOK, user)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/user/:username", GetUserByName)
}
```

## Multiple Path Parameters

This example demonstrates how to document an endpoint with multiple path parameters:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/openapi"
    "github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/src/swagger"
    "net/http"
)

// Comment DTO
type Comment struct {
    ID      int64  `json:"id,omitempty"`
    PostID  int64  `json:"postId,omitempty"`
    Content string `json:"content"`
    Author  string `json:"author"`
}

// Swagger documentation for GET /posts/{postId}/comments/{commentId}
var _ = swagger.Swagger().Path("/posts/{postId}/comments/{commentId}").
    Get(func(op openapi.Operation) {
        op.Summary("Get a specific comment on a post").
            Description("Returns a single comment from a specific post").
            OperationID("getPostComment").
            Tag("comments").
            Produces(mime.ApplicationJSON).
            PathParameter("postId", func(p openapi.Parameter) {
                p.Description("ID of the post").
                    Type("integer").Format("int64").
                    Required(true)
            }).
            PathParameter("commentId", func(p openapi.Parameter) {
                p.Description("ID of the comment").
                    Type("integer").Format("int64").
                    Required(true)
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").SchemaFromDTO(&Comment{})
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid ID supplied")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Comment not found")
            })
    }).
    Doc()

// Handler function
func GetPostComment(c *gin.Context) {
    postID := c.Param("postId")
    commentID := c.Param("commentId")
    
    // Implementation details
    comment := Comment{
        ID: 1,
        PostID: 1,
        Content: "This is a great post!",
        Author: "John Doe",
    }
    c.JSON(http.StatusOK, comment)
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/posts/:postId/comments/:commentId", GetPostComment)
}
```

## Path Parameter with Enum Values

Here's how to document a path parameter that only accepts certain values:

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

// Swagger documentation for GET /reports/{reportType}/download
var _ = swagger.Swagger().Path("/reports/{reportType}/download").
    Get(func(op openapi.Operation) {
        op.Summary("Download a report").
            Description("Downloads a report in the specified format").
            OperationID("downloadReport").
            Tag("reports").
            Produces("application/pdf", "application/vnd.ms-excel", "text/csv").
            PathParameter("reportType", func(p openapi.Parameter) {
                p.Description("Type of report to download").
                    Type("string").
                    Enum("sales", "inventory", "customers", "analytics")  // Only these values are allowed
            }).
            QueryParameter("year", func(p openapi.Parameter) {
                p.Description("Year for the report data").
                    Type("integer").Format("int32").
                    Minimum(2000, false).Maximum(2030, false)
            }).
            QueryParameter("month", func(p openapi.Parameter) {
                p.Description("Month for the report data (1-12)").
                    Type("integer").Format("int32").
                    Minimum(1, false).Maximum(12, false)
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Report downloaded successfully")
            }).
            Response(http.StatusBadRequest, func(r openapi.Response) {
                r.Description("Invalid parameters")
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("Report not found")
            })
    }).
    Doc()

// Handler function
func DownloadReport(c *gin.Context) {
    reportType := c.Param("reportType")
    year := c.Query("year")
    month := c.Query("month")
    
    // Implementation details
    c.Header("Content-Disposition", "attachment; filename=report.pdf")
    c.Data(http.StatusOK, "application/pdf", []byte("Report data"))
}

// Router setup
func setupRoutes(router *gin.Engine) {
    router.GET("/reports/:reportType/download", DownloadReport)
}
```