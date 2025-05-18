package controller

import (
	"net/http"
	"strconv"

	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

// Pet DTOs
type Pet struct {
	ID        int64     `json:"id,omitempty"`
	Category  *Category `json:"category,omitempty"`
	Name      string    `json:"name"`
	PhotoUrls []string  `json:"photoUrls"`
	Tags      []Tag     `json:"tags,omitempty"`
	Status    string    `json:"status,omitempty"` // can be "available", "pending", "sold"
}

type Category struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Tag struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ApiResponse struct {
	Code    int32  `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// PetTag defines the Swagger API tag for Pet
var _ = swagger.Swagger().
	Tag("pet", func(tag openapi.Tag) {
		tag.Description("Everything about your Pets").
			ExternalDocumentation("http://swagger.io", "Find out more")
	})

// UploadImage swagger documentation
var _ = swagger.Swagger().Path("/pet/{petId}/uploadImage").
	Post(func(operation openapi.Operation) {
		operation.Summary("uploads an image").
			OperationID("uploadFile").
			Tag("pet").
			Consumes("multipart/form-data").
			Produce(mime.ApplicationJSON).
			PathParameter("petId", func(p openapi.Parameter) {
				p.Description("ID of pet to update").
					MaxLength(64).
					Type("integer").Format("int64")
			}).
			FormParameter("additionalMetadata", func(p openapi.Parameter) {
				p.Description("Additional data to pass to server").Required(false).Type("string")
			}).
			FormParameter("file", func(p openapi.Parameter) {
				p.Description("file to upload").Required(false).Type("file")
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").SchemaFromDTO(&ApiResponse{})
			}).
			Security("petstore_auth", "read:pets", "write:pets")
	}).
	Doc()

// UploadImage handler
func UploadImage(c *gin.Context) {
	c.JSON(http.StatusOK, ApiResponse{Code: 200, Type: "success", Message: "Image uploaded"})
}

// FindByTags swagger documentation
var _ = swagger.Swagger().Path("/pet/findByTags").
	Get(func(op openapi.Operation) {
		op.Summary("Finds Pets by tags").
			Description("Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.").
			OperationID("findPetsByTags").
			Tag("pet").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			QueryParameter("tags", func(p openapi.Parameter) {
				p.Description("Tags to filter by").
					Required(true).
					Type("array").
					CollectionFormat("multi").
					Items(func(item openapi.Schema) { item.Type("string") })
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Schema(openapi_spec.SchemaEntity{
						Type:  "array",
						Items: &openapi_spec.SchemaEntity{Ref: "#/definitions/Pet"},
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid tag value")
			}).
			Deprecated(true).
			Security("petstore_auth", "read:pets", "write:pets")
	}).
	Doc()

// FindByTags handler
func FindByTags(c *gin.Context) {
	tags := c.QueryArray("tags")
	if len(tags) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag value"})
		return
	}
	// Dummy response
	pets := []Pet{
		{ID: 1, Name: "doggie", Tags: []Tag{{ID: 1, Name: "tag1"}}, PhotoUrls: []string{"http://example.com/photo1.jpg"}, Status: "available"},
	}
	c.JSON(http.StatusOK, pets)
}

// AddPet swagger documentation
var _ = swagger.Swagger().Path("/pet").
	Post(func(operation openapi.Operation) {
		operation.Summary("Add a new pet to the store").
			OperationID("addPet").
			Tag("pet").
			Consumes(mime.ApplicationJSON, mime.ApplicationXML).
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			BodyParameter(func(p openapi.Parameter) {
				p.Description("Pet object that needs to be added to the store").
					Required(true).
					SchemaFromDTO(&Pet{})
			}).
			Response(http.StatusMethodNotAllowed, func(r openapi.Response) {
				r.Description("Invalid input")
			}).
			Security("petstore_auth", "read:pets", "write:pets")
	}).
	Doc()

// AddPet handler
func AddPet(c *gin.Context) {
	var pet Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Invalid input"})
		return
	}
	c.JSON(http.StatusOK, pet)
}

// UpdatePet swagger documentation
var _ = swagger.Swagger().Path("/pet").
	Put(func(op openapi.Operation) {
		op.Summary("Update an existing pet").
			OperationID("updatePet").
			Tag("pet").
			Consumes(mime.ApplicationJSON, mime.ApplicationXML).
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			BodyParameter(func(p openapi.Parameter) {
				p.Description("Pet object that needs to be added to the store").
					Required(true).
					SchemaFromDTO(&Pet{})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid ID supplied")
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Pet not found")
			}).
			Response(http.StatusMethodNotAllowed, func(r openapi.Response) {
				r.Description("Validation exception")
			}).
			Security("petstore_auth", "read:pets", "write:pets")
	}).
	Doc()

// UpdatePet handler
func UpdatePet(c *gin.Context) {
	var pet Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID supplied"})
		return
	}
	c.JSON(http.StatusOK, pet)
}

// FindByStatus swagger documentation
var _ = swagger.Swagger().Path("/pet/findByStatus").
	Get(func(op openapi.Operation) {
		op.Summary("Finds Pets by status").
			Description("Multiple status values can be provided with comma separated strings").
			OperationID("findPetsByStatus").
			Tag("pet").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			QueryParameter("status", func(p openapi.Parameter) {
				p.Description("Status values that need to be considered for filter").
					Required(true).
					Type("array").
					CollectionFormat("multi").
					Items(func(item openapi.Schema) {
						item.Type("string").
							Enum("available", "pending", "sold").
							Default("available")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Schema(openapi_spec.SchemaEntity{
						Type:  "array",
						Items: &openapi_spec.SchemaEntity{Ref: "#/definitions/Pet"},
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid status value")
			}).
			Security("petstore_auth", "read:pets", "write:pets")
	}).
	Doc()

// FindByStatus handler
func FindByStatus(c *gin.Context) {
	status := c.QueryArray("status")
	if len(status) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}
	// Dummy response
	pets := []Pet{
		{ID: 1, Name: "doggie", PhotoUrls: []string{"http://example.com/photo1.jpg"}, Status: "available"},
	}
	c.JSON(http.StatusOK, pets)
}

// GetPetByID swagger documentation
var _ = swagger.Swagger().Path("/pet/{petId}").
	Get(func(op openapi.Operation) {
		op.Summary("Find pet by ID").
			Description("Returns a single pet").
			OperationID("getPetById").
			Tag("pet").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			PathParameter("petId", func(p openapi.Parameter) {
				p.Description("ID of pet to return").Type("integer").Format("int64")
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").SchemaFromDTO(&Pet{})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid ID supplied")
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Pet not found")
			}).
			Security("api_key").
			Security("petstore_auth", "read:pets", "write:pets")
	}).
	Doc()

// GetPetByID handler
func GetPetByID(c *gin.Context) {
	petIDstr := c.Param("petId")
	petID, _ := strconv.ParseInt(petIDstr, 10, 64)
	// Dummy response
	pet := Pet{ID: petID, Name: "doggie", PhotoUrls: []string{"http://example.com/photo1.jpg"}, Status: "available"}
	c.JSON(http.StatusOK, pet)
}

// UpdatePetWithForm swagger documentation
var _ = swagger.Swagger().Path("/pet/{petId}").
	Post(func(op openapi.Operation) {
		op.Summary("Updates a pet in the store with form data").
			OperationID("updatePetWithForm").
			Tag("pet").
			Consumes("application/x-www-form-urlencoded").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			PathParameter("petId", func(p openapi.Parameter) {
				p.Description("ID of pet that needs to be updated").Type("integer").Format("int64")
			}).
			FormParameter("name", func(p openapi.Parameter) {
				p.Description("Updated name of the pet").Required(false).Type("string")
			}).
			FormParameter("status", func(p openapi.Parameter) {
				p.Description("Updated status of the pet").Required(false).Type("string")
			}).
			Response(http.StatusMethodNotAllowed, func(r openapi.Response) {
				r.Description("Invalid input")
			}).
			Security("petstore_auth", "read:pets", "write:pets")
	}).
	Doc()

// UpdatePetWithForm handler
func UpdatePetWithForm(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Pet updated"})
}

// DeletePet swagger documentation
var _ = swagger.Swagger().Path("/pet/{petId}").
	Delete(func(op openapi.Operation) {
		op.Summary("Deletes a pet").
			OperationID("deletePet").
			Tag("pet").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			HeaderParameter("api_key", func(p openapi.Parameter) {
				p.Description("").Required(false).Type("string")
			}).
			PathParameter("petId", func(p openapi.Parameter) {
				p.Description("Pet id to delete").Type("integer").Format("int64")
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid ID supplied")
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Pet not found")
			}).
			Security("petstore_auth", "read:pets", "write:pets")
	}).
	Doc()

// DeletePet handler
func DeletePet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Pet deleted"})
}
