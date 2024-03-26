package middleware

import (
	"backbu/internal/data"
	"backbu/pkg/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AdminOperation(c *gin.Context){
	if(!checkAuth(c, []string{"ADMIN"})){
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

func UserOperation(c *gin.Context){
	if(!checkAuth(c, []string{"ADMIN", "CUSTOMER"})){
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	
	c.Next()
}


func checkAuth(c *gin.Context, roles []string) bool{
	/* Obtenemos la secret key de las variables de entorno */
	secret, err := env.MustGet("SECRET_KEY"); if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return false
	}

	/* Comprobamos que la cookie contenga la autorización */
	tokenString, err := c.Cookie("Authorization")
	log.Println(c.Cookie("Authorization"))
	if err != nil {
		return false
	}

 	/* 
	 * Decodificamos el token, Para validar que el token es correcto, se implementa una
	 * lambda que lo hace, para más info consultar la documentación de la biblioteca
	 * github.com/golang-jwt/jwt
	 */
 	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
   
		return []byte(secret), nil
   	})

	/* Si el token es válido y se puede obtener el mapa de información, continuamos */
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Comprobamos la fecha de expiración
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return false
		}

		/* Comprobamos que el ID es un ObjectID */
		objectId, err := primitive.ObjectIDFromHex(claims["sub"].(string))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: "Cookie bad formatted. Can not recognize ID"})
			return false
		}

		/* Cmprobamos si el usuario del token coincide con alguno en BD */
		usuario, err := db.Get[data.User]("users", bson.M{"_id": objectId})
		if err != nil {
			return false
		}
		
		/* 
		 * Finalmente, comprobamos si el rol es válido para operar y si no,
		 * devolvemos falso
		 */
		for _, rol := range roles {
			if rol == usuario.Role.Code {
				return true
			}
		}
	}
	return false
}