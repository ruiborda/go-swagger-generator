package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

// User DTO
type User struct {
	ID         int64  `json:"id,omitempty" yaml:"id,omitempty"`
	Username   string `json:"username,omitempty" yaml:"username,omitempty"`
	FirstName  string `json:"firstName,omitempty" yaml:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty" yaml:"lastName,omitempty"`
	Email      string `json:"email,omitempty" yaml:"email,omitempty"`
	Password   string `json:"password,omitempty" yaml:"password,omitempty"`
	Phone      string `json:"phone,omitempty" yaml:"phone,omitempty"`
	UserStatus int32  `json:"userStatus,omitempty" yaml:"userStatus,omitempty"` // User status
}

// UserTag defines the Swagger API tag for User
var _ = swagger.Swagger().
	Tag("user", func(tag openapi.Tag) {
		tag.Description("Operations about user").
			ExternalDocumentation("http://swagger.io", "Find out more about our store") // Example had this, but seems like Petstore link
	})

// CreateUser swagger documentation
var _ = swagger.Swagger().Path("/user").
	Post(func(op openapi.Operation) {
		op.Summary("Create user").
			Description("This can only be done by the logged in user.").
			OperationID("createUser").
			Tag("user").
			RequestBody(func(rb openapi.RequestBody) {
				rb.Description("Created user object").Required(true).
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&User{})
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) { // Original default response, can be improved
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
			RequestBody(func(rb openapi.RequestBody) {
				rb.Description("List of user object").Required(true).
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.Schema(func(s openapi.Schema) {
							s.Type("array").ItemsRef("#/components/schemas/User")
						})
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation")
				// Potentially return array of users or confirmation
			})
	}).
	Doc()

// CreateUsersWithArray handler
func CreateUsersWithArray(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Users created"})
}

// CreateUsersWithList swagger documentation
var _ = swagger.Swagger().Path("/user/createWithList").
	Post(func(op openapi.Operation) {
		op.Summary("Creates list of users with given input list"). // Title adjusted slightly for clarity
										OperationID("createUsersWithListInput").
										Tag("user").
										RequestBody(func(rb openapi.RequestBody) {
				rb.Description("List of user object").Required(true).
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.Schema(func(s openapi.Schema) {
							s.Type("array").ItemsRef("#/components/schemas/User")
						})
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation")
			})
	}).
	Doc()

// CreateUsersWithList handler
func CreateUsersWithList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Users created"})
}

// LoginUser swagger documentation
var _ = swagger.Swagger().Path("/user/login").
	Get(func(op openapi.Operation) {
		op.Summary("Logs user into the system").
			OperationID("loginUser").
			Tag("user").
			QueryParameter("username", func(p openapi.Parameter) {
				p.Description("The user name for login").Required(true).
					Schema(func(s openapi.Schema) { s.Type("string") })
			}).
			QueryParameter("password", func(p openapi.Parameter) {
				p.Description("The password for login in clear text").Required(true).
					Schema(func(s openapi.Schema) { s.Type("string") })
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Header("X-Rate-Limit", func(h openapi.Header) {
						h.Description("calls per hour allowed by the user").
							Schema(func(s openapi.Schema) { s.Type("integer").Format("int32") })
					}).
					Header("X-Expires-After", func(h openapi.Header) {
						h.Description("date in UTC when token expires").
							Schema(func(s openapi.Schema) { s.Type("string").Format("date-time") })
					}).
					Content(mime.TextPlain, func(mt openapi.MediaType) { // Assuming token is plain text
						mt.Schema(func(s openapi.Schema) { s.Type("string") })
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
	c.Header("X-Expires-After", "2025-01-01T00:00:00Z")
	c.String(http.StatusOK, "logged in user session token")
}

// LogoutUser swagger documentation
var _ = swagger.Swagger().Path("/user/logout").
	Get(func(op openapi.Operation) {
		op.Summary("Logs out current logged in user session").
			OperationID("logoutUser").
			Tag("user").
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation")
				// No content needed for logout usually
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
			PathParameter("username", func(p openapi.Parameter) {
				p.Description("The name that needs to be fetched. Use user1 for testing.").Required(true).
					Schema(func(s openapi.Schema) { s.Type("string") })
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&User{})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&User{})
					})
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
	user := User{ID: 1, Username: username, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Phone: "123-456-7890", UserStatus: 1}
	c.JSON(http.StatusOK, user)
}

// UpdateUser swagger documentation
var _ = swagger.Swagger().Path("/user/{username}").
	Put(func(op openapi.Operation) {
		op.Summary("Update user"). // Updated from "Updated user"
						Description("This can only be done by the logged in user.").
						OperationID("updateUser").
						Tag("user").
						PathParameter("username", func(p openapi.Parameter) {
				p.Description("name that need to be updated").Required(true).
					Schema(func(s openapi.Schema) { s.Type("string") })
			}).
			RequestBody(func(rb openapi.RequestBody) {
				rb.Description("Updated user object").Required(true).
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&User{})
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) { // Added success response
				r.Description("User updated successfully")
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
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// DeleteUser swagger documentation
var _ = swagger.Swagger().Path("/user/{username}").
	Delete(func(op openapi.Operation) {
		op.Summary("Delete user").
			Description("This can only be done by the logged in user.").
			OperationID("deleteUser").
			Tag("user").
			PathParameter("username", func(p openapi.Parameter) {
				p.Description("The name that needs to be deleted").Required(true).
					Schema(func(s openapi.Schema) { s.Type("string") })
			}).
			Response(http.StatusOK, func(r openapi.Response) { // Added success response
				r.Description("User deleted successfully")
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
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
