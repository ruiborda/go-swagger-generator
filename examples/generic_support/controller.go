package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-service-common/dto"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/v2/src/swagger"
	"net/http"
)

type UserController struct {
}

var _ = swagger.Swagger().Path("/userinfo").
	Get(func(op openapi.Operation) {
		op.Summary("Get User Information").
			Tag("Generic Examples").
			OperationID("getUserInfo").
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation - user data returned in generic response").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&GenericResponse[UserData]{})
					})
			})
	}).Doc()

func (this *UserController) GetUserInfo(c *gin.Context) {
	c.JSON(http.StatusOK, GenericResponse[UserData]{
		Success: true,
		Data:    UserData{ID: 1, Username: "johndoe", Email: "john.doe@example.com"},
	})
}

var _ = swagger.Swagger().Path("/productinfo").
	Get(func(op openapi.Operation) {
		op.Summary("Get Product Information").
			Tag("Generic Examples").
			OperationID("getProductInfo").
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation - product data returned in generic response").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&GenericResponse[ProductData]{})
					})
			})
	}).Doc()

func (this *UserController) GetProductInfo(c *gin.Context) {
	c.JSON(http.StatusOK, GenericResponse[ProductData]{
		Success: true,
		Data:    ProductData{SKU: "PROD123", ProductName: "Awesome Gadget", Price: 99.99},
	})
}

var _ = swagger.Swagger().Path("/userslist").
	Get(func(op openapi.Operation) {
		op.Summary("Get List of Users Information").
			Tag("Generic Examples").
			OperationID("getUsersListInfo").
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation - list of user data returned in generic response").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&GenericResponse[[]*UserData]{})
					})
			})
	}).Doc()

func (this *UserController) GetUsersList(c *gin.Context) {
	c.JSON(http.StatusOK, GenericResponse[[]*UserData]{
		Success: true,
		Data: []*UserData{
			{ID: 1, Username: "johndoe", Email: "john.doe@example.com"},
			{ID: 2, Username: "janedoe", Email: "jane.doe@example.com"},
		},
	})
}

var _ = swagger.Swagger().Path("/api/v1/login-with-google").
	Post(func(op openapi.Operation) {
		op.Summary("Login or signup with Google OAuth").
			OperationID("LoginWithGoogle").
			Tag("auth").
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Login response with user details and JWT token").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&dto.Response[UserData]{})
					})
			})
	}).
	Doc()
