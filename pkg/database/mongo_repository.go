/*
 *
 * Repositorio genérico que realiza un CRUD en una entidad mongo
 *
 */

package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/* Recibe uno o más campos de un elemento y lo busca en base de datos */
func Get[T any](collection string, filter T) (*mongo.SingleResult){
	mongoDb, err := Connect(); if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	
	return mongoDb.Collection(collection).FindOne(context.Background(), filter)
}

/* Dada una colección, obtiene todos los elementos existentes */
func GetAll[T any](collection string) ([]T, error){
	mongoDb, err := Connect(); if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	/* Obtiene un objeto mongo Cursor con la información de las entidades */
	cursor, err := mongoDb.Collection(collection).Find(context.Background(), bson.D{{}}); if err != nil {
		return nil, err
	}

	/* Construye el array de entidades a través del cursor */
	var data []T
	err = cursor.All(context.Background(), &data)

	return data, err
}

/* Recibe un ID y actualiza la entidad con los datos seleccionados */
func Update[T any](collection string, filter bson.D, data T) (*mongo.UpdateResult, error){
	mongoDb, err := Connect(); if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return mongoDb.Collection(collection).UpdateByID(context.Background(), filter, data)
}

/* Recibe el nombre de una colección, un elemento a introducir y efectua la creación de la entidad */
func Create[T any](collection string, data T) (*mongo.InsertOneResult, error) {
	mongoDb, err := Connect(); if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return mongoDb.Collection(collection).InsertOne(context.Background(), data)
}


