package api

import (
	"backbu/internal/business"

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
	router.GET("/product", business.GetAllProducts)
	router.GET("/product/:id", business.GetProduct)
	router.POST("/product", business.CreateProduct)
	router.PUT("/product/:id", business.UpdateProduct)
	router.DELETE("/product/:id", business.DeleteProduct)

	/* CRUD PROVEEDORES */
	router.GET("/supplier", )
	router.GET("/supplier/:id", )
	router.POST("/supplier", )
	router.PUT("/supplier/:id", )
	router.DELETE("/supplier/:id", )

	/* CRUD PEDIDOS */
	router.GET("/order", )
	router.GET("/order/:id", )
	router.POST("/order", )
	router.PUT("/order/:id", )
	router.DELETE("/order/:id", )

	/* MANAGE USERS */
	router.POST("/register", )
	router.POST("/login", )
}
