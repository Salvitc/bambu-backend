package business

import (
	"backbu/internal/data"
	"backbu/pkg/database"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Devuelve el producto que coincida con el ID pasado como parámetro de la URL
func GetProduct(c *gin.Context){
	// El ID debe ser de tipo ObjectID
	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	result, err := db.Get[data.Product]("products", objectId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result)	
}

//Devuelve todos las entidades de la colección productos
func GetAllProducts(c *gin.Context){
	result, err := db.GetAll[data.Product]("products")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result);
}

/* Crear un producto en la base de datos */
func CreateProduct(c *gin.Context){
	var producto data.Product

	/* Intenta decodificar el Body y rellenar el objeto producto */
	err := json.NewDecoder(c.Request.Body).Decode(&producto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
		return
	}

	/* Crea el producto en base de datos */
	result, err := db.Create("products", producto)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}

	/* Devuelve el objeto con el ID asociado por Mongo */
	c.IndentedJSON(http.StatusCreated, result)
}

/* Actualizar producto en base de datos */
func UpdateProduct(c *gin.Context) {
	// El ID debe ser de tipo ObjectID
	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* Se decodifica el body de la petición en un tipo "Product" */
	var producto data.Product;
	err = json.NewDecoder(c.Request.Body).Decode(&producto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* Con el ID y los atributos a modificar, se actualiza la base de datos */
	result, err := db.Update("products", bson.M{"_id": objectId}, producto)
	if (err != nil) {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}
	
	c.IndentedJSON(http.StatusOK, result)
}

func DeleteProduct(c *gin.Context) {
	// El ID debe ser de tipo ObjectID
	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	result, err := db.Delete("products", bson.M{"_id": objectId})
	if(err != nil){
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
	}

	c.IndentedJSON(http.StatusAccepted, result)
}