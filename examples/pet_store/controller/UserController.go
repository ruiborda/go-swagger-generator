package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

// User DTO
type User struct {
	ID         int64  `json:"id,omitempty"`
	Username   string `json:"username,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Phone      string `json:"phone,omitempty"`
	UserStatus int32  `json:"userStatus,omitempty"` // User status
}

// UserTag defines the Swagger API tag for User
var _ = swagger.Swagger().
	Tag("user", func(tag openapi.Tag) {
		tag.Description("Operations about user").
			ExternalDocumentation("http://swagger.io", "Find out more about our store")
	})

// CreateUser swagger documentation
var _ = swagger.Swagger().Path("/user").
	Post(func(op openapi.Operation) {
		op.Summary("Create user").
			Description("This can only be done by the logged in user.").
			OperationID("createUser").
			Tag("user").
			Consumes(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			BodyParameter(func(p openapi.Parameter) {
				p.Description("Created user object").Required(true).SchemaFromDTO(&User{})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation")
			})
	}).
	Doc()

// CreateUser handler
func CreateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User created"}) // Simplified
}

// CreateUsersWithArray swagger documentation
var _ = swagger.Swagger().Path("/user/createWithArray").
	Post(func(op openapi.Operation) {
		op.Summary("Creates list of users with given input array").
			OperationID("createUsersWithArrayInput").
			Tag("user").
			Consumes(mime.ApplicationJSON).
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

// CreateUsersWithArray handler
func CreateUsersWithArray(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Users created"}) // Simplified
}

// CreateUsersWithList swagger documentation
var _ = swagger.Swagger().Path("/user/createWithList").
	Post(func(op openapi.Operation) {
		op.Summary("Creates list of users with given input array").
			OperationID("createUsersWithListInput").
			Tag("user").
			Consumes(mime.ApplicationJSON).
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

// CreateUsersWithList handler
func CreateUsersWithList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Users created"}) // Simplified
}

// LoginUser swagger documentation
var _ = swagger.Swagger().Path("/user/login").
	Get(func(op openapi.Operation) {
		op.Summary("Logs user into the system").
			OperationID("loginUser").
			Tag("user").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			QueryParameter("username", func(p openapi.Parameter) {
				p.Description("The user name for login").Required(true).Type("string")
			}).
			QueryParameter("password", func(p openapi.Parameter) {
				p.Description("The password for login in clear text").Required(true).Type("string")
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Schema(openapi_spec.SchemaEntity{Type: "string"}).
					Header("X-Rate-Limit", func(h openapi.Header) {
						h.Type("integer").Format("int32").Description("calls per hour allowed by the user")
					}).
					Header("X-Expires-After", func(h openapi.Header) {
						h.Type("string").Format("date-time").Description("date in UTC when token expires")
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid username/password supplied")
			})
	}).
	Doc()

// LoginUser handler
func LoginUser(c *gin.Context) {
	c.Header("X-Rate-Limit", "5000")
	c.Header("X-Expires-After", "2025-01-01T00:00:00Z") // Example future date
	c.String(http.StatusOK, "logged in user session token")
}

// LogoutUser swagger documentation
var _ = swagger.Swagger().Path("/user/logout").
	Get(func(op openapi.Operation) {
		op.Summary("Logs out current logged in user session").
			OperationID("logoutUser").
			Tag("user").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation")
			})
	}).
	Doc()

// LogoutUser handler
func LogoutUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User logged out"})
}

// GetUserByName swagger documentation
var _ = swagger.Swagger().Path("/user/{username}").
	Get(func(op openapi.Operation) {
		op.Summary("Get user by user name").
			OperationID("getUserByName").
			Tag("user").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			PathParameter("username", func(p openapi.Parameter) {
				p.Description("The name that needs to be fetched. Use user1 for testing. ").Type("string")
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

// GetUserByName handler
func GetUserByName(c *gin.Context) {
	username := c.Param("username")
	// Dummy response
	user := User{ID: 1, Username: username, FirstName: "John", LastName: "Doe"}
	c.JSON(http.StatusOK, user)
}

// UpdateUser swagger documentation
var _ = swagger.Swagger().Path("/user/{username}").
	Put(func(op openapi.Operation) {
		op.Summary("Updated user").
			Description("This can only be done by the logged in user.").
			OperationID("updateUser").
			Tag("user").
			Consumes(mime.ApplicationJSON).
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

// UpdateUser handler
func UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User updated"}) // Simplified
}

// DeleteUser swagger documentation
var _ = swagger.Swagger().Path("/user/{username}").
	Delete(func(op openapi.Operation) {
		op.Summary("Delete user").
			Description("This can only be done by the logged in user.").
			OperationID("deleteUser").
			Tag("user").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			PathParameter("username", func(p openapi.Parameter) {
				p.Description("The name that needs to be deleted").Type("string")
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid username supplied")
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("User not found")
			})
	}).
	Doc()

// DeleteUser handler
func DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"}) // Simplified
}
