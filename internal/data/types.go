package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* ENTIDADES PRIMARIAS */
type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	Price       float32            `json:"price" bson:"price,omitempty"`
	Category    string             `json:"category" bson:"category,omitempty"`
	InStock     bool               `json:"in_stock" bson:"in_stock,omitempty"`
  Amount      uint8              `json:"amount" bson:"amount,omitempty"`
	Images      []string           `json:"images" bson:"images,omitempty"`
}

type User struct {
	ID       primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string               `json:"username" bson:"username,omitempty"`
	Password string               `json:"password" bson:"password,omitempty"`
	Email    string               `json:"email" bson:"email,omitempty"`
  Phone    string               `json:"phone" bson:"phone,omitempty"`
	Name     string               `json:"name" bson:"name,omitempty"`
	Lastname string               `json:"lastname" bson:"lastname,omitempty"`
	Address  string               `json:"address" bson:"address,omitempty"`
  Role     Role                 `json:"role" bson:"role,omitempty"`
  Cart     []CartItem           `json:"cart" bson:"cart,omitempty"`
  Wishlist []primitive.ObjectID `json:"wishlist" bson:"wishlist,omitempty"`
}

type CartItem struct {
  ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
  Price     float32            `json:"price" bson:"price"` 
  Amount    uint8              `json:"amount" bson:"amount"`
}

type Role struct {
	Code      string              `json:"code" bson:"code"`
}

type Order struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
  OrderID   int32                `json:"order_id" bson:"order_id"`
  Date      primitive.DateTime   `json:"date" bson:"date"`
	UserID    primitive.ObjectID   `json:"user_id" bson:"user_id"`
	Items     []OrderItem          `json:"products" bson:"products"`
	Amount    float32              `json:"amount" bson:"amount"`
}

type OrderItem struct {
  ProductID primitive.ObjectID  `json:"product_id" bson:"product_id"`
  Amount    uint8               `json:"amount" bson:"amount"`
}

type Review struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Rating    uint8              `json:"rating,omitempty" bson:"rating,omitempty"`
	Comment   string             `json:"comment,omitempty" bson:"comment,omitempty"`
}

/* UTILES */
type ExtendedOrder struct {
  ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  OrderID  int32              `json:"order_id" bson:"order_id"`
  Date     primitive.DateTime `json:"date" bson:"date"`
  User     User               `json:"user" bson:"user"`
  Products []Product          `json:"products" bson:"products"`
  Amount   float32            `json:"amount" bson:"amount"`
}

type JsonError struct {
	Message string `json:"error"`
}

// Error implements error.
func (j JsonError) Error() string {
	panic("unimplemented")
}
