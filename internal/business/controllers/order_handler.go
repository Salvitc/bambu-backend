package controllers

import (
	"backbu/internal/data"
	"backbu/pkg/database"
	"log"
	"time"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Devuelve el Pedido que coincida con el ID pasado como parámetro de la URL
func GetOrder(c *gin.Context) {
	// El ID del pedido debe ser de tipo ObjectID
	orderId, err := primitive.ObjectIDFromHex(c.Param("orderid"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	// El ID del usuario debe ser de tipo ObjectID
	userid, err := primitive.ObjectIDFromHex(c.Param("userid"))
	if err != nil {
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

// Devuelve todos las entidades de la colección Pedidos
func GetAllOrders(c *gin.Context) {
	/* Obtiene todos los Pedidos */
	var extendedOrders []data.ExtendedOrder

	result, err := db.GetAll[data.Order]("orders")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	/* Construye los Pedidos extendidos con los datos de los Usuarios y Productos asociados */
	extendedOrders, err = buildExtendedOrders(result)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, extendedOrders)
}

// Devuelve los Pedidos que coincidan con el rango de fechas pasado como parámetro de la URL
func GetOrdersByDateRange(c *gin.Context) {
	/* Obtiene los parámetros de la URL */
	startDate := c.Param("startdate")
	endDate := c.Param("enddate")

	/* Obtiene todos los Pedidos que estén en el rango de fechas */
	result, err := db.GetBy[data.Order]("orders", bson.M{"date": bson.M{"$gte": startDate, "$lte": endDate}})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	/* Construye los Pedidos extendidos con los datos de los Usuarios y Productos asociados */
	extendedOrders, err := buildExtendedOrders(result)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, extendedOrders)
}

// Devuelve los Pedidos que coincidan con el ID pasado como parámetro de la URL
func GetAllUserOrders(c *gin.Context) {
	// El ID del usuario debe ser de tipo ObjectID
	userId, err := primitive.ObjectIDFromHex(c.Param("userid"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* obtiene todos los Pedidos que pertenezcan al usuario */
	result, err := db.GetBy[data.Order]("orders", bson.M{"user_id": userId})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	extendedOrders, err := buildExtendedOrders(result)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, extendedOrders)
}

/* Crear un Pedido en la base de datos */
func CreateOrder(c *gin.Context) {
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
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	// El ID debe ser de tipo ObjectID
	userId, err := primitive.ObjectIDFromHex(c.Param("userid"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* Se decodifica el body de la petición en un tipo "Order" */
	var pedido data.Order
	err = json.NewDecoder(c.Request.Body).Decode(&pedido)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* Con el ID y los atributos a modificar, se actualiza la base de datos */
	result, err := db.Update("orders", bson.M{"_id": orderId, "user_id": userId}, pedido)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

func DeleteOrder(c *gin.Context) {
	// El ID debe ser de tipo ObjectID
	orderId, err := primitive.ObjectIDFromHex(c.Param("orderid"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	// El ID debe ser de tipo ObjectID
	userId, err := primitive.ObjectIDFromHex(c.Param("userid"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	result, err := db.Delete("orders", bson.M{"_id": orderId, "user_id": userId})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
	}

	c.IndentedJSON(http.StatusAccepted, result)
}

func buildExtendedOrders(orders []*data.Order) ([]data.ExtendedOrder, error) {
	var extendedOrders []data.ExtendedOrder

	/* Itera sobre los Pedidos para obtener los datos de los Usuarios y Productos asociados */
	for _, order := range orders {
		user, err := db.Get[data.User]("users", bson.M{"_id": order.UserID})
		if err != nil {
			return nil, err
		}

		var products []data.Product
		for _, product := range order.Items {
			p, err := db.Get[data.Product]("products", bson.M{"_id": product.ProductID})
			if err != nil {
				return nil, err
			}
			p.Amount = product.Amount
			products = append(products, *p)
		}

		extendedOrders = append(extendedOrders, data.ExtendedOrder{
			ID:       order.ID,
			OrderID:  order.OrderID,
			User:     *user,
			Products: products,
			Amount:   order.Amount,
		})
	}
	return extendedOrders, nil
}

func DumpCartToOrder(c *gin.Context) {
	/* Obtenemos la secret key de las variables de entorno */
	secret, err := env.MustGet("SECRET_KEY")
	if err != nil {
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
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
	}

	/* Obtenemos el usuario dado el ID */
	usuario, err := db.Get[data.User]("users", bson.M{"_id": objectId})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, data.JsonError{Message: err.Error()})
		return
	}

	/* Insertamos su carrito como un pedido */
	var order data.Order
	order.UserID = usuario.ID
	for _, item := range usuario.Cart {
		var orderItem data.OrderItem
		orderItem.ProductID = item.ProductID
		orderItem.Amount = item.Amount

		order.Items = append(order.Items, orderItem)
		order.Amount += item.Price * float32(item.Amount)
	}

	order.Amount += order.Amount * 0.21 // IVA
	order.Amount += 5.99                // Gastos de envío
	order.Date = primitive.DateTime(time.Now().UTC().UnixNano() / int64(time.Millisecond))

	result, err := db.Create("orders", order)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}

	/* Limpiamos el carrito */
	usuario.Cart = nil
	_, err = db.Update("users", bson.M{"_id": usuario.ID}, usuario)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
