package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	// Definimos la URI de conexión (estándar local)
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	// Contexto con timeout (si tarda más de 10s, cancela)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Conectamos
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Verificamos con un Ping
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("No se pudo conectar a MongoDB: ", err)
	}

	fmt.Println("¡Conectado a MongoDB exitosamente!")
	return client
}

// Helper para obtener una colección rápido
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("MediaTrackerDB").Collection(collectionName)
}
