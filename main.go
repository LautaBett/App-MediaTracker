package main

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"proyecto-go-mongo/configs"
	"proyecto-go-mongo/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func main() {
	// 1. Conectar a la BD
	client := configs.ConnectDB()
	collection = configs.GetCollection(client, "media_items")

	// 2. Iniciar el Router
	router := gin.Default()

	// Servir el Frontend en la raíz
	router.StaticFile("/", "./static/index.html")

	// --- RUTA INTELIGENTE: CREAR O ACTUALIZAR (POST) ---
	router.POST("/media", func(c *gin.Context) {
		var newItem models.Media
		if err := c.BindJSON(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Buscamos si ya existe (Duplicado)
		filter := bson.M{
			"title": bson.M{"$regex": primitive.Regex{Pattern: "^" + newItem.Title + "$", Options: "i"}},
			"type":  newItem.Type,
		}

		var existingItem models.Media
		err := collection.FindOne(ctx, filter).Decode(&existingItem)

		if err == nil {
			// YA EXISTE: Comparamos progreso
			oldNum := extractNumber(existingItem.Progress)
			newNum := extractNumber(newItem.Progress)

			if newNum > oldNum {
				// Es más nuevo -> Actualizamos
				update := bson.M{
					"$set": bson.M{
						"progress": newItem.Progress,
						"link":     newItem.Link,
						"title":    newItem.Title,
					},
				}
				_, err := collection.UpdateOne(ctx, bson.M{"_id": existingItem.ID}, update)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"mensaje": "Serie actualizada al nuevo capítulo", "id": existingItem.ID})
			} else {
				// Es más viejo -> No hacemos nada
				c.JSON(http.StatusOK, gin.H{"mensaje": "No se actualizó: El capítulo ya es antiguo o igual."})
			}
		} else {
			// NO EXISTE: Creamos nuevo
			result, err := collection.InsertOne(ctx, newItem)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"mensaje": "Nueva serie guardada!", "id": result.InsertedID})
		}
	})

	// --- RUTA: LEER TODO (GET) ---
	router.GET("/media", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var results []models.Media
		// Ordenamos por ID descendente (lo último creado/editado primero)
		opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})

		cursor, err := collection.Find(ctx, bson.M{}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar"})
			return
		}
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer datos"})
			return
		}
		c.JSON(http.StatusOK, results)
	})

	// --- RUTA: ACTUALIZAR MANUALMENTE (PUT) ---
	// Para cuando usas el botón "Edit" del frontend
	router.PUT("/media/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		var updateData models.Media
		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		update := bson.M{
			"$set": bson.M{
				"title":    updateData.Title,
				"progress": updateData.Progress,
				"link":     updateData.Link,
			},
		}

		result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Actualizado", "modificados": result.ModifiedCount})
	})

	// --- RUTA: BORRAR (DELETE) ---
	router.DELETE("/media/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al borrar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Eliminado", "borrados": result.DeletedCount})
	})

	// 3. Correr servidor
	router.Run("localhost:8080")
}

// --- FUNCIÓN AYUDANTE ---
func extractNumber(text string) float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	match := re.FindString(text)
	if match == "" {
		return 0
	}
	num, _ := strconv.ParseFloat(match, 64)
	return num
}
