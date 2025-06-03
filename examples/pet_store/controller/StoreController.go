package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

// Store DTOs
type Order struct {
	ID       int64     `json:"id,omitempty" yaml:"id,omitempty"`
	PetID    int64     `json:"petId,omitempty" yaml:"petId,omitempty"`
	Quantity int32     `json:"quantity,omitempty" yaml:"quantity,omitempty"`
	ShipDate time.Time `json:"shipDate,omitempty" yaml:"shipDate,omitempty"`
	Status   string    `json:"status,omitempty" yaml:"status,omitempty"` // can be "placed", "approved", "delivered"
	Complete bool      `json:"complete,omitempty" yaml:"complete,omitempty"`
}

// StoreTag defines the Swagger API tag for Store
var _ = swagger.Swagger().
	Tag("store", func(tag openapi.Tag) {
		tag.Description("Access to Petstore orders")
	})

// GetInventory swagger documentation
var _ = swagger.Swagger().Path("/store/inventory").
	Get(func(op openapi.Operation) {
		op.Summary("Returns pet inventories by status").
			Description("Returns a map of status codes to quantities").
			OperationID("getInventory").
			Tag("store").
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.Schema(func(s openapi.Schema) {
							s.Type("object").
								AdditionalProperties(true, func(addPropSchema openapi.Schema) {
									addPropSchema.Type("integer").Format("int32")
								})
						})
					})
			})
	}).
	Doc()

// GetInventory handler
func GetInventory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"available": 10, "pending": 5, "sold": 2})
}

// PlaceOrder swagger documentation
var _ = swagger.Swagger().Path("/store/order").
	Post(func(op openapi.Operation) {
		op.Summary("Place an order for a pet").
			OperationID("placeOrder").
			Tag("store").
			RequestBody(func(rb openapi.RequestBody) {
				rb.Description("order placed for purchasing the pet").Required(true).
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Order{})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) { // Added XML support as per original Produces
						mt.SchemaFromDTO(&Order{})
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Order{})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Order{})
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid Order")
			})
	}).
	Doc()

// PlaceOrder handler
func PlaceOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order"})
		return
	}
	if order.ID == 0 { // Simple ID generation
		order.ID = 1
	}
	c.JSON(http.StatusOK, order)
}

// GetOrderByID swagger documentation
var _ = swagger.Swagger().Path("/store/order/{orderId}").
	Get(func(op openapi.Operation) {
		op.Summary("Find purchase order by ID").
			Description("For valid response try integer IDs with value >= 1 and <= 10. Other values will generate exceptions").
			OperationID("getOrderById").
			Tag("store").
			PathParameter("orderId", func(p openapi.Parameter) {
				p.Description("ID of order that needs to be fetched").Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64").Minimum(1, false).Maximum(10, false)
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Order{})
					}).
					Content(mime.ApplicationXML, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&Order{})
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid ID supplied")
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Order not found")
			})
	}).
	Doc()

// GetOrderByID handler
func GetOrderByID(c *gin.Context) {
	orderIDstr := c.Param("orderId")
	orderID, _ := strconv.ParseInt(orderIDstr, 10, 64)
	order := Order{ID: orderID, PetID: 1, Quantity: 1, Status: "placed", Complete: false, ShipDate: time.Now()}
	c.JSON(http.StatusOK, order)
}

// DeleteOrder swagger documentation
var _ = swagger.Swagger().Path("/store/order/{orderId}").
	Delete(func(op openapi.Operation) {
		op.Summary("Delete purchase order by ID").
			Description("For valid response try integer IDs with positive integer value. Negative or non-integer values will generate API errors").
			OperationID("deleteOrder").
			Tag("store").
			PathParameter("orderId", func(p openapi.Parameter) {
				p.Description("ID of the order that needs to be deleted").Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64").Minimum(1, false)
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) { // Added success response
				r.Description("Order deleted successfully")
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Invalid ID supplied")
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Order not found")
			})
	}).
	Doc()

// DeleteOrder handler
func DeleteOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
