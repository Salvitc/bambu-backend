package api

import (
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
	router.GET("/product", )
	router.GET("/product/:id", )
	router.POST("/product", )
	router.PUT("/product/:id", )
	router.DELETE("/product/:id", )

	/* CRUD PROVEEDORES */
	router.GET("/supplier", )
	router.GET("/supplier/:id", )
	router.POST("/supplier", )
	router.PUT("/supplier/:id", )
	router.DELETE("/supplier/:id", )

	/* MANAGE USERS */
	router.POST("/register", )
	router.POST("/login", )
}
