package controllers

import (
	"backbu/internal/data"
	"backbu/pkg/auth"
	"backbu/pkg/database"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//AUTENTICACIÓN DE USUARIOS
//Registra un usuario de tipo Customer
func Register(c *gin.Context){
	
}

func Login(c *gin.Context){

}

//OPERACIONES DE ADMINISTRACIÓN DE USUARIOS
//Devuelve el Usuario que coincida con el ID pasado como parámetro de la URL
func GetUser(c *gin.Context){
	// El ID debe ser de tipo ObjectID
	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* obtiene el Usuario dado el ID */
	result, err := db.Get[data.User]("users", bson.M{"_id": objectId})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result)	
}

//Devuelve todos las entidades de la colección Usuarios
func GetAllUsers(c *gin.Context){
	/* Obtiene todos los Usuarios */
	result, err := db.GetAll[data.User]("users")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result);
}

/* Crear un Usuario en la base de datos */
func CreateUser(c *gin.Context){
	var usuario data.User

	/* Intenta decodificar el Body y rellenar el objeto Usuario */
	err := json.NewDecoder(c.Request.Body).Decode(&usuario)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
		return
	}

	usuario.Password, err = auth.HashPassword(usuario.Password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
		return
	}

	/* Crea el Usuario en base de datos */
	result, err := db.Create("users", usuario)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}

	/* Devuelve el objeto con el ID asociado por Mongo */
	c.IndentedJSON(http.StatusCreated, result)
}

/* Actualizar Usuario en base de datos */
func UpdateUser(c *gin.Context) {
	// El ID debe ser de tipo ObjectID
	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* Se decodifica el body de la petición en un tipo "User" */
	var usuario data.User;
	err = json.NewDecoder(c.Request.Body).Decode(&usuario)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	if(usuario.Password != ""){
		usuario.Password, err = auth.HashPassword(usuario.Password)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
			return
		}
	}

	/* Con el ID y los atributos a modificar, se actualiza la base de datos */
	result, err := db.Update("users", bson.M{"_id": objectId}, usuario)
	if (err != nil) {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}
	
	c.IndentedJSON(http.StatusOK, result)
}

/* Elimina un usuario de la base de datos */
func DeleteUser(c *gin.Context) {
	// El ID debe ser de tipo ObjectID
	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	result, err := db.Delete("users", bson.M{"_id": objectId})
	if(err != nil){
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
	}

	c.IndentedJSON(http.StatusAccepted, result)
}

