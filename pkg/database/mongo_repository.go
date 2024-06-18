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
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Encapsula la obtención y gestión de errores del cliente mongo */
func GetClient() *mongo.Database {
	mongoDb, err := connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return mongoDb
}

/****************************************************************/

/* Recibe uno o más campos de un elemento y lo busca en base de datos */
func Get[T any](collection string, filter bson.M) (*T, error) {
	result := GetClient().Collection(collection).FindOne(context.Background(), filter)

	var data *T
	if err := result.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

/* Dada una colección, obtiene todos los elementos existentes */
func GetAll[T any](collection string) ([]*T, error) {
	/* Obtiene un objeto mongo Cursor con la información de las entidades */
	cursor, err := GetClient().Collection(collection).Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	/* Construye el array de entidades a través del cursor */
	var data []*T
	err = cursor.All(context.Background(), &data)

	return data, err
}

func GetBy[T any](collection string, filter bson.M) ([]*T, error) {
	/* Obtiene un objeto mongo Cursor con la información de las entidades */
	cursor, err := GetClient().Collection(collection).Find(context.Background(), filter)
	if err != nil {

		return nil, err
	}

	/* Construye el array de entidades a través del cursor */
	var data []*T
	err = cursor.All(context.Background(), &data)

	return data, err
}

/* Recibe el nombre de una colección, un elemento a introducir y efectua la creación de la entidad */
func Create[T any](collection string, data T) (*mongo.InsertOneResult, error) {
	return GetClient().Collection(collection).InsertOne(context.Background(), data)
}

/* Recibe un ID y actualiza la entidad con los datos seleccionados */
func Update[T any](collection string, filter bson.M, data T) (*mongo.UpdateResult, error) {
	result, err := GetClient().Collection(collection).UpdateOne(context.Background(), filter, bson.M{"$set": data})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Replace[T any](collection string, filter bson.M, data T) (*mongo.UpdateResult, error) {
	result, err := GetClient().Collection(collection).ReplaceOne(context.Background(), filter, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Delete(collection string, filter bson.M) (*mongo.DeleteResult, error) {
	result, err := GetClient().Collection(collection).DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CounterIncrement() (int, error) {
	filter := bson.M{"_id": "orders"}
	update := bson.M{"$inc": bson.M{"order_id": 1}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		OrderId int `bson:"order_id"`
	}

	collection := GetClient().Collection("counters")
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, options).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.OrderId, nil
}
