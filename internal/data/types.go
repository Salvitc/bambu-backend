package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* ENTIDADES PRIMARIAS */
type Product struct {
	ID			primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	Name		string				`json:"name" bson:"name,omitempty"`
	Description string				`json:"description" bson:"description,omitempty"`
	Price		float32				`json:"price" bson:"price,omitempty"`
	Category	string				`json:"category" bson:"category,omitempty"`
	InStock		bool				`json:"in_stock" bson:"in_stock,omitempty"`
	Images		[]string			`json:"images" bson:"images,omitempty"`
	Reviews		[]review			`json:"reviews" bson:"reviews,omitempty"`
}

type review struct {
	UserID		primitive.ObjectID	`json:"user_id,omitempty" bson:"user_id,omitempty"`
	Rating		uint8				`json:"rating,omitempty" bson:"rating,omitempty"`
	Comment		string				`json:"comment,omitempty" bson:"comment,omitempty"`
}

/* UTILES */
type JsonError struct {
	Message		string				`json:"error"`
}
