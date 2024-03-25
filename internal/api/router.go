package api

import (
	"backbu/internal/business/controllers"

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
	router.POST("/product", controllers.CreateProduct)
	router.PUT("/product/:id", controllers.UpdateProduct)
	router.DELETE("/product/:id", controllers.DeleteProduct)

	/* CRUD PEDIDOS */
	router.GET("/order", controllers.GetAllOrders)
	router.GET("/order/:userid", controllers.GetAllUserOrders)
	router.GET("/order/:orderid", controllers.GetOrder)
	router.POST("/order", controllers.CreateOrder)
	router.PUT("/order/:orderid", controllers.UpdateOrder)
	router.DELETE("/order/:orderid", controllers.DeleteOrder)

	/* CRUD USUARIOS */
	router.GET("/user", controllers.GetAllUsers)
	router.GET("/user/:id", controllers.GetUser)
	router.POST("/user", controllers.CreateUser)
	router.PUT("/user/:id", controllers.UpdateUser)
	router.DELETE("/user/:id", controllers.DeleteUser)

	/* AUTH USUARIOS */
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)
}
