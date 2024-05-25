package controllers

import (
	"backbu/internal/data"
	"backbu/pkg/auth"
	db "backbu/pkg/database"
	"time"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//AUTENTICACIÓN DE USUARIOS
//Registra un usuario de tipo Customer
func Register(c *gin.Context){
	var usuario data.User

	/* Intenta decodificar el Body y rellenar el objeto Usuario */
	err := json.NewDecoder(c.Request.Body).Decode(&usuario)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
		return
	}

	/* Hasheamos la contraseña */
	usuario.Password, err = auth.HashPassword(usuario.Password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
		return
	}

	/* Los usuarios creados desde el register son de tipo cliente */
	usuario.Role.Code = "CUSTOMER"

	/* Crea el Usuario en base de datos */
	result, err := db.Create("users", usuario)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}

	/* Devuelve el objeto con el ID asociado por Mongo */
	c.IndentedJSON(http.StatusCreated, result)
}

func Login(c *gin.Context){
	/* Obtenemos la secret key de las variables de entorno */
	secret, err := env.MustGet("SECRET_KEY"); if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	/* Cuando un usuario se logea se espera mail y contraseña */
	var loggedUser struct {
		Email string
		Password string
	}

	/* Si se puede decodear el body, se ha recibido lo esperado */
	err = json.NewDecoder(c.Request.Body).Decode(&loggedUser)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
		return
	}

	/* obtenemos el Usuario dado el mail */
	usuario, err := db.Get[data.User]("users", bson.M{"email": loggedUser.Email})

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	/* Comparamos la contraseña introducida con el hash de base de datos */
	if (!auth.CheckPasswordHash(loggedUser.Password, usuario.Password)){
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	/* Generamos token para el usuario */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usuario.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	/* Terminamos de generar el token con la SECRET_KEY */
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}

	/* Le ofrecemos al usuario una Cookie con su token de acceso */
	c.SetSameSite(http.SameSiteLaxMode)
 	c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "", "", false, true)

 	c.JSON(http.StatusOK, gin.H{})
}

func Logout(c *gin.Context){
  /* Borramos la cookie de autenticación */
  c.SetCookie("Authorization", "", 0, "", "", false, true)
  c.JSON(http.StatusOK, gin.H{})
}

func GetUserFromToken(c *gin.Context){
  /* Obtenemos la secret key de las variables de entorno */
  secret, err := env.MustGet("SECRET_KEY"); if err != nil {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  /* Obtenemos el token de la cookie */
  tokenString, err := c.Cookie("Authorization")
  if err != nil {
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }

  /* Parseamos el token */
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(secret), nil
  })
  if err != nil {
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }

  /* Obtenemos el ID del usuario */
  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }

	// El ID debe ser de tipo ObjectID
	objectId, err := primitive.ObjectIDFromHex(claims["sub"].(string))
	if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

  /* Obtenemos el usuario dado el ID */
  usuario, err := db.Get[data.User]("users", bson.M{"_id": objectId})
  if err != nil {
    c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
    return
  }

  c.IndentedJSON(http.StatusOK, usuario)
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

/* Actualiza la wishlist de un usuario */
func UpdateWishlist(c *gin.Context){
  // El ID debe ser de tipo ObjectID
  objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
  if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
  }

  usuario, err := db.Get[data.User]("users", bson.M{"_id": objectId})
  if err != nil {
    c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
    return
  }

  /* Se decodifica el body de la petición en un tipo "User" */
  var wishlist []string;
  err = json.NewDecoder(c.Request.Body).Decode(&wishlist)
  if err != nil {
    c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
  }
  
  /* Convertimos la wishlist a un array de ObjectID */
  var objectID_wishlist []primitive.ObjectID
  for _, id := range wishlist {
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
      c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
      return
    }
    objectID_wishlist = append(objectID_wishlist, objectId)
  }

  usuario.Wishlist = objectID_wishlist

  /* Se añade el producto a la wishlist */
  result, err := db.Update("users", bson.M{"_id": objectId}, usuario)
  if (err != nil) {
    c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
    return
  }
  
  c.IndentedJSON(http.StatusOK, result)
}

/* Actualiza el carrito de un usuario */
func UpdateCart(c *gin.Context){
  // El ID debe ser de tipo ObjectID
  objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
  if err != nil{
    	c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
  }

  usuario, err := db.Get[data.User]("users", bson.M{"_id": objectId})
  if err != nil {
    c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
    return
  }

  /* Se decodifica el body de la petición en un tipo "User" */
  var cart []data.CartItem;
  err = json.NewDecoder(c.Request.Body).Decode(&cart)
  if err != nil {
    c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
  }

  usuario.Cart = cart

  /* Se añade el producto al carrito */
  result, err := db.Update("users", bson.M{"_id": objectId}, usuario)
  if (err != nil) {
    c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
    return
  }
  
  c.IndentedJSON(http.StatusOK, result)
}
