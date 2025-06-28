package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/v2/src/middleware"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

// GenericResponse is a generic DTO for API responses.
// The type parameter T will be replaced by a concrete type upon instantiation.
type GenericResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

// UserData is a concrete DTO for user information.
type UserData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ProductData is another concrete DTO for product information.
type ProductData struct {
	SKU         string  `json:"sku"`
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
}

func main() {
	router := gin.Default()

	// Configure OpenAPI 3.0 documentation
	ConfigureOpenAPI(router)

	router.GET("/v1/userinfo", GetUserInfo)
	router.GET("/v1/productinfo", GetProductInfo)
	router.GET("/v1/userslist", GetUsersList)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("OpenAPI 3.0 JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

func ConfigureOpenAPI(router *gin.Engine) {
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
		Title:    "Generic DTO Support API (v2)",
	}))

	doc := swagger.Swagger()

	doc.Info(func(info openapi.Info) {
		info.Title("API with Generic DTO Support").
			Version("1.0.0").
			Description("Demonstrates OpenAPI documentation for generic DTOs like GenericResponse[T].")
	})

	doc.Server("http://localhost:8080/v1", func(server openapi.Server) {
		server.Description("Local development server (v1)")
	})
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

func GetUserInfo(c *gin.Context) {
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

func GetProductInfo(c *gin.Context) {
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

func GetUsersList(c *gin.Context) {
	c.JSON(http.StatusOK, GenericResponse[[]*UserData]{
		Success: true,
		Data: []*UserData{
			{ID: 1, Username: "johndoe", Email: "john.doe@example.com"},
			{ID: 2, Username: "janedoe", Email: "jane.doe@example.com"},
		},
	})
}
