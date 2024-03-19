/*
 *
 * Cliente Mongo que lee de las variables de entorno del sistema y genera una conexión con la BD 
 *
 */

package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gofor-little/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Cerrojo para controlar el acceso a la creación del cliente */
var lock = &sync.Mutex{}

/* variable para hacer el cliente singleton y no poder crear más de una instancia */
var instance *mongo.Database

func Connect() (*mongo.Database, error) {

	if instance == nil {
		/* Bloqueamos el cerrojo y nos aseguramos que se libera al salir de la función */
		lock.Lock()
		defer lock.Unlock()

		/* Buscamos las variables de entorno y si no, devuelve error */
		host, err := env.MustGet("MONGODB_HOST"); if err != nil {
			return nil, err
		}
		port, err := env.MustGet("MONGODB_PORT"); if err != nil {
			return nil, err
		}
		db, err := env.MustGet("MONGODB_DB"); if err != nil {
			return nil, err
		}

		/* Creamos la url de conexión y la seteamos a cliente*/
		uri := fmt.Sprintf("mongodb://%s:%s/%s", host, port, db)
		log.Println(uri)
		clientOpt := options.Client().ApplyURI(uri)

		cliente, err := mongo.Connect(context.Background(), clientOpt); if err != nil {
			return nil, err
		}

		/* Probamos la conectividad del cliente */
		err = cliente.Ping(context.Background(), nil); if err != nil {
			return nil, err
		}

		instance = cliente.Database(db)
	}

	/* devolvemos el cliente de bd funcionando */
	return instance, nil
}
