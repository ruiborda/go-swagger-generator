package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/examples/pet_store/controller"
)

func setupRoutes(router *gin.Engine) {

	// Pet routes
	router.POST("/v2/pet/:petId/uploadImage", controller.UploadImage)
	router.POST("/v2/pet", controller.AddPet)
	router.PUT("/v2/pet", controller.UpdatePet)
	router.GET("/v2/pet/findByStatus", controller.FindByStatus)
	router.GET("/v2/pet/findByTags", controller.FindByTags)
	router.GET("/v2/pet/:petId", controller.GetPetByID)
	router.POST("/v2/pet/:petId", controller.UpdatePetWithForm)
	router.DELETE("/v2/pet/:petId", controller.DeletePet)

	// Store routes
	router.GET("/v2/store/inventory", controller.GetInventory)
	router.POST("/v2/store/order", controller.PlaceOrder)
	router.GET("/v2/store/order/:orderId", controller.GetOrderByID)
	router.DELETE("/v2/store/order/:orderId", controller.DeleteOrder)

	// User routes
	router.POST("/v2/user", controller.CreateUser)
	router.POST("/v2/user/createWithArray", controller.CreateUsersWithArray)
	router.POST("/v2/user/createWithList", controller.CreateUsersWithList)
	router.GET("/v2/user/login", controller.LoginUser)
	router.GET("/v2/user/logout", controller.LogoutUser)
	router.GET("/v2/user/:username", controller.GetUserByName)
	router.PUT("/v2/user/:username", controller.UpdateUser)
	router.DELETE("/v2/user/:username", controller.DeleteUser)
}
