package api

import (
	"backbu/internal/business/controllers"
	"backbu/internal/business/middleware"

	"github.com/gin-gonic/gin"
)

/* Devuelve un objeto router con los endpoints de la api inicializados */
func Router() *gin.Engine {
	var router = gin.Default()

	inicializarRutas(router)

	return router
}

func inicializarRutas(router *gin.Engine){
	/* CRUD PRODUCTOS */
	router.GET("/product", controllers.GetAllProducts)
	router.GET("/product/:id", controllers.GetProduct)
	router.POST("/product", middleware.AdminOperation, controllers.CreateProduct)
	router.PUT("/product/:id", middleware.AdminOperation, controllers.UpdateProduct)
	router.DELETE("/product/:id", middleware.AdminOperation, controllers.DeleteProduct)

	/* CRUD PEDIDOS */
	router.GET("/order", middleware.AdminOperation, controllers.GetAllOrders)
	router.GET("/order/:userid", middleware.UserOperation, controllers.GetAllUserOrders)
	router.GET("/order/:userid/:orderid", middleware.UserOperation, controllers.GetOrder)
	router.POST("/order", middleware.UserOperation, controllers.CreateOrder)
	router.PUT("/order/:userid/:orderid", middleware.UserOperation, controllers.UpdateOrder)
	router.DELETE("/order/:userid/:orderid", middleware.AdminOperation, controllers.DeleteOrder)

	/* CRUD USUARIOS */
	router.GET("/user", middleware.AdminOperation, controllers.GetAllUsers)
	router.GET("/user/:id", controllers.GetUser)
	router.POST("/user", middleware.AdminOperation, controllers.CreateUser)
	router.PUT("/user/:id", middleware.UserOperation, controllers.UpdateUser)
	router.DELETE("/user/:id", middleware.UserOperation, controllers.DeleteUser)

	/* AUTH USUARIOS */
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)
}
