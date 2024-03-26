package controllers

import (
	"backbu/internal/data"
	"backbu/pkg/database"
	"log"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Devuelve el Pedido que coincida con el ID pasado como parámetro de la URL
func GetOrder(c *gin.Context){
	// El ID del pedido debe ser de tipo ObjectID
	orderId, err := primitive.ObjectIDFromHex(c.Param("orderid"))
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	// El ID del usuario debe ser de tipo ObjectID
	userid, err := primitive.ObjectIDFromHex(c.Param("userid"))
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* obtiene el Pedido su id y que pertenezcan al usuario */
	result, err := db.Get[data.Order]("orders", bson.M{"_id": orderId, "user_id": userid})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result)	
}

//Devuelve todos las entidades de la colección Pedidos
func GetAllOrders(c *gin.Context){
	/* Obtiene todos los Pedidos */
	result, err := db.GetAll[data.Order]("orders")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result);
}

//Devuelve los Pedidos que coincidan con el ID pasado como parámetro de la URL
func GetAllUserOrders(c *gin.Context){
	// El ID del usuario debe ser de tipo ObjectID
	userId, err := primitive.ObjectIDFromHex(c.Param("userid"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* obtiene todos los Pedidos que pertenezcan al usuario */
	result, err := db.Get[data.Order]("orders", bson.M{"user_id": userId})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result)	
}

/* Crear un Pedido en la base de datos */
func CreateOrder(c *gin.Context){
	var pedido data.Order

	/* Intenta decodificar el Body y rellenar el objeto Pedido */
	err := json.NewDecoder(c.Request.Body).Decode(&pedido)
	log.Println(pedido)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
		return
	}

	/* Crea el Pedido en base de datos */
	result, err := db.Create("orders", pedido)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}

	/* Devuelve el objeto con el ID asociado por Mongo */
	c.IndentedJSON(http.StatusCreated, result)
}

/* Actualizar Pedido en base de datos */
func UpdateOrder(c *gin.Context) {
	// El ID debe ser de tipo ObjectID
	orderId, err := primitive.ObjectIDFromHex(c.Param("orderid"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	// El ID debe ser de tipo ObjectID
	userId, err := primitive.ObjectIDFromHex(c.Param("userid"))
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* Se decodifica el body de la petición en un tipo "Order" */
	var pedido data.Order;
	err = json.NewDecoder(c.Request.Body).Decode(&pedido)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* Con el ID y los atributos a modificar, se actualiza la base de datos */
	result, err := db.Update("orders", bson.M{"_id": orderId, "user_id": userId}, pedido)
	if (err != nil) {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}
	
	c.IndentedJSON(http.StatusOK, result)
}

func DeleteOrder(c *gin.Context) {
	// El ID debe ser de tipo ObjectID
	orderId, err := primitive.ObjectIDFromHex(c.Param("orderid"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	// El ID debe ser de tipo ObjectID
	userId, err := primitive.ObjectIDFromHex(c.Param("userid"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	result, err := db.Delete("orders", bson.M{"_id": orderId, "user_id": userId})
	if(err != nil){
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
	}

	c.IndentedJSON(http.StatusAccepted, result)
}