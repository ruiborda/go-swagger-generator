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

	// Define API routes
	// These handlers would typically serve actual data.
	router.GET("/v1/userinfo", func(c *gin.Context) {
		c.JSON(http.StatusOK, GenericResponse[UserData]{
			Success: true,
			Data:    UserData{ID: 1, Username: "johndoe", Email: "john.doe@example.com"},
		})
	})
	router.GET("/v1/productinfo", func(c *gin.Context) {
		c.JSON(http.StatusOK, GenericResponse[ProductData]{
			Success: true,
			Data:    ProductData{SKU: "PROD123", ProductName: "Awesome Gadget", Price: 99.99},
		})
	})
	router.GET("/v1/userslist", func(c *gin.Context) {
		c.JSON(http.StatusOK, GenericResponse[[]*UserData]{
			Success: true,
			Data: []*UserData{
				{ID: 1, Username: "johndoe", Email: "john.doe@example.com"},
				{ID: 2, Username: "janedoe", Email: "jane.doe@example.com"},
			},
		})
	})

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("OpenAPI 3.0 JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

// ConfigureOpenAPI sets up the OpenAPI 3.0 documentation.
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

	// Register DTOs. SchemaFromDTO will now handle generic types.
	// The generator should create component names like:
	// - main_GenericResponse_main_UserData
	// - main_GenericResponse_main_ProductData
	// - main_GenericResponse_ListOf_main_UserData (or similar for slice type argument)
	// It will also register UserData and ProductData if not already registered.
	_, _ = doc.SchemaFromDTO(&GenericResponse[UserData]{})
	_, _ = doc.SchemaFromDTO(&GenericResponse[ProductData]{})
	_, _ = doc.SchemaFromDTO(&GenericResponse[[]*UserData]{}) // Generic with a slice type argument
	// Individual DTOs also need to be registered if they are to be referenced directly or are complex.
	_, _ = doc.SchemaFromDTO(&UserData{})
	_, _ = doc.SchemaFromDTO(&ProductData{})

	// Document an endpoint that returns GenericResponse[UserData]
	doc.Path("/userinfo").
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

	// Document an endpoint that returns GenericResponse[ProductData]
	doc.Path("/productinfo").
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
	
	// Document an endpoint that returns GenericResponse[[]*UserData]
	doc.Path("/userslist").
		Get(func(op openapi.Operation) {
			op.Summary("Get List of Users Information").
				Tag("Generic Examples").
				OperationID("getUsersListInfo").
				Response(http.StatusOK, func(r openapi.Response) {
					r.Description("Successful operation - list of user data returned in generic response").
						Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
							// Note: The SchemaFromDTO for GenericResponse[[]*UserData] will register
							// a component like 'main_GenericResponse_ListOf_main_UserData'.
							// The 'ListOf_main_UserData' part itself might also be registered as a separate component
							// if `doc.SchemaFromDTO(&[]*UserData{})` was called, or implicitly if needed.
							mt.SchemaFromDTO(&GenericResponse[[]*UserData]{})
						})
				})
		}).Doc()
}
