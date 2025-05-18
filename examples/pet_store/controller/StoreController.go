package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

// Store DTOs
type Order struct {
	ID       int64     `json:"id,omitempty"`
	PetID    int64     `json:"petId,omitempty"`
	Quantity int32     `json:"quantity,omitempty"`
	ShipDate time.Time `json:"shipDate,omitempty"`
	Status   string    `json:"status,omitempty"` // can be "placed", "approved", "delivered"
	Complete bool      `json:"complete,omitempty"`
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
			Produces(mime.ApplicationJSON).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Schema(openapi_spec.SchemaEntity{
						Type: "object",
						AdditionalProperties: &openapi_spec.SchemaEntity{
							Type:   "integer",
							Format: "int32",
						},
					})
			}).
			Security("api_key")
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
			Consumes(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			BodyParameter(func(p openapi.Parameter) {
				p.Description("order placed for purchasing the pet").Required(true).SchemaFromDTO(&Order{})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").SchemaFromDTO(&Order{})
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
			Description("For valid response try integer IDs with value >= 1 and <= 10. Other values will generated exceptions").
			OperationID("getOrderById").
			Tag("store").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			PathParameter("orderId", func(p openapi.Parameter) {
				p.Description("ID of pet that needs to be fetched").
					Type("integer").Format("int64").
					Minimum(1, false).Maximum(10, false)
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

// GetOrderByID handler
func GetOrderByID(c *gin.Context) {
	orderIDstr := c.Param("orderId")
	orderID, _ := strconv.ParseInt(orderIDstr, 10, 64)
	// Dummy response
	order := Order{ID: orderID, PetID: 1, Quantity: 1, Status: "placed", Complete: false}
	c.JSON(http.StatusOK, order)
}

// DeleteOrder swagger documentation
var _ = swagger.Swagger().Path("/store/order/{orderId}").
	Delete(func(op openapi.Operation) {
		op.Summary("Delete purchase order by ID").
			Description("For valid response try integer IDs with positive integer value. Negative or non-integer values will generate API errors").
			OperationID("deleteOrder").
			Tag("store").
			Produces(mime.ApplicationJSON, mime.ApplicationXML).
			PathParameter("orderId", func(p openapi.Parameter) {
				p.Description("ID of the order that needs to be deleted").
					Type("integer").Format("int64").Minimum(1, false)
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
