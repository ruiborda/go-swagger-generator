package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

// Pet DTOs
type Pet struct {
	ID        int64     `json:"id,omitempty" yaml:"id,omitempty"`
	Category  *Category `json:"category,omitempty" yaml:"category,omitempty"`
	Name      string    `json:"name" yaml:"name"`
	PhotoUrls []string  `json:"photoUrls" yaml:"photoUrls"`
	Tags      []*Tag    `json:"tags,omitempty" yaml:"tags,omitempty"`     // Changed to pointer array for consistency
	Status    string    `json:"status,omitempty" yaml:"status,omitempty"` // can be "available", "pending", "sold"
}

type Category struct {
	ID   int64  `json:"id,omitempty" yaml:"id,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type Tag struct {
	ID   int64  `json:"id,omitempty" yaml:"id,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type ApiResponse struct {
	Code    int32  `json:"code,omitempty" yaml:"code,omitempty"`
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

// PetTag defines the Swagger API tag for Pet
var _ = swagger.Swagger().
	Tag("pet", func(tag openapi.Tag) {
		tag.Description("Everything about your Pets").
			ExternalDocumentation("http://swagger.io", "Find out more")
	})

// UploadImage swagger documentation
var _ = swagger.Swagger().Path("/pet/{petId}/uploadImage"). // Path should not include the base path /v2
								Post(func(op openapi.Operation) {
		op.Summary("uploads an image").
			OperationID("uploadFile").
			Tag("pet").
			PathParameter("petId", func(p openapi.Parameter) {
				p.Description("ID of pet to update").
					Required(true). // Path parameters are always required
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			RequestBody(func(rb openapi.RequestBody) {
				rb.Description("Image and metadata to upload").
					Required(true).
					Content(mime.ApplicationOctetStream, func(mt openapi.MediaType) {
						mt.Schema(func(s openapi.Schema) {
							s.Type("object").
								Property("additionalMetadata", func(prop openapi.Schema) {
									prop.Type("string").Description("Additional data to pass to server")
								}).
								Property("file", func(prop openapi.Schema) {
									prop.Type("string").Format("binary").Description("file to upload")
								})
						})
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&ApiResponse{})
					})
			})
	}).
	Doc()

// UploadImage handler
func UploadImage(c *gin.Context) {
	c.JSON(http.StatusOK, ApiResponse{Code: http.StatusOK, Type: "success", Message: "Image uploaded"})
}

// FindByTags swagger documentation
var _ = swagger.Swagger().Path("/pet/findByTags").
	Get(func(op openapi.Operation) {
		op.Summary("Finds Pets by tags").
			Description("Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.").
			OperationID("findPetsByTags").
			Tag("pet").
			QueryParameter("tags", func(p openapi.Parameter) {
				p.Description("Tags to filter by").
					Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("array").
							Items(func(item openapi.Schema) { item.Type("string") })
					}).
					Style("form").Explode(false) // For comma separated values for array (default for query)
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.Schema(func(s openapi.Schema) {
							s.Type("array").ItemsRef("#/components/schemas/Pet")
						})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.Schema(func(s openapi.Schema) {
							s.Type("array").ItemsRef("#/components/schemas/Pet")
						})
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid tag value")
			}).
			Deprecated(true)
	}).
	Doc()

// FindByTags handler
func FindByTags(c *gin.Context) {
	tags := c.QueryArray("tags")
	if len(tags) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag value"})
		return
	}
	pets := []Pet{
		{ID: 1, Name: "doggie", Tags: []*Tag{{ID: 1, Name: "tag1"}}, PhotoUrls: []string{"http://example.com/photo1.jpg"}, Status: "available"},
	}
	c.JSON(http.StatusOK, pets)
}

// AddPet swagger documentation
var _ = swagger.Swagger().Path("/pet").
	Post(func(op openapi.Operation) {
		op.Summary("Add a new pet to the store").
			OperationID("addPet").
			Tag("pet").
			RequestBody(func(rb openapi.RequestBody) {
				rb.Description("Pet object that needs to be added to the store").
					Required(true).
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					})
			}).
			Response(http.StatusMethodNotAllowed, func(r openapi.Response) { // Should be 200 or 201 for success
				r.Description("Invalid input") // This response seems misplaced for a successful creation
			}).
			Response(http.StatusOK, func(r openapi.Response) { // Added a success response
				r.Description("Successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					})
			})
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
			RequestBody(func(rb openapi.RequestBody) {
				rb.Description("Pet object that needs to be updated in the store").
					Required(true).
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) { // Added success response
				r.Description("Successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid ID supplied")
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Pet not found")
			}).
			Response(http.StatusMethodNotAllowed, func(r openapi.Response) {
				r.Description("Validation exception")
			})
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
			QueryParameter("status", func(p openapi.Parameter) {
				p.Description("Status values that need to be considered for filter").
					Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("array").
							Items(func(item openapi.Schema) {
								item.Type("string").
									Enum("available", "pending", "sold").
									Default("available")
							})
					}).Explode(false).Style("form")
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.Schema(func(s openapi.Schema) {
							s.Type("array").ItemsRef("#/components/schemas/Pet")
						})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.Schema(func(s openapi.Schema) {
							s.Type("array").ItemsRef("#/components/schemas/Pet")
						})
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid status value")
			})
	}).
	Doc()

// FindByStatus handler
func FindByStatus(c *gin.Context) {
	status := c.QueryArray("status")
	if len(status) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}
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
			PathParameter("petId", func(p openapi.Parameter) {
				p.Description("ID of pet to return").Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Pet{})
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid ID supplied")
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Pet not found")
			})
	}).
	Doc()

// GetPetByID handler
func GetPetByID(c *gin.Context) {
	petIDstr := c.Param("petId")
	petID, _ := strconv.ParseInt(petIDstr, 10, 64)
	pet := Pet{ID: petID, Name: "doggie", PhotoUrls: []string{"http://example.com/photo1.jpg"}, Status: "available"}
	c.JSON(http.StatusOK, pet)
}

// UpdatePetWithForm swagger documentation
var _ = swagger.Swagger().Path("/pet/{petId}").
	Post(func(op openapi.Operation) {
		op.Summary("Updates a pet in the store with form data").
			OperationID("updatePetWithForm").
			Tag("pet").
			PathParameter("petId", func(p openapi.Parameter) {
				p.Description("ID of pet that needs to be updated").Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			RequestBody(func(rb openapi.RequestBody) {
				rb.Content(mime.ApplicationOctetStream, func(mt openapi.MediaType) {
					mt.Schema(func(s openapi.Schema) {
						s.Type("object").
							Property("name", func(prop openapi.Schema) {
								prop.Type("string").Description("Updated name of the pet")
							}).
							Property("status", func(prop openapi.Schema) {
								prop.Type("string").Description("Updated status of the pet")
							})
						// For form parameters, 'required' is on the property level within the schema
					})
				})
			}).
			Response(http.StatusMethodNotAllowed, func(r openapi.Response) {
				r.Description("Invalid input")
			})
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
			HeaderParameter("api_key", func(p openapi.Parameter) {
				p.Description("Session token for authentication").Required(false).
					Schema(func(s openapi.Schema) {
						s.Type("string")
					})
			}).
			PathParameter("petId", func(p openapi.Parameter) {
				p.Description("Pet id to delete").Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) { // Added success response
				r.Description("Pet deleted successfully")
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid ID supplied")
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Pet not found")
			})
	}).
	Doc()

// DeletePet handler
func DeletePet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Pet deleted"})
}
