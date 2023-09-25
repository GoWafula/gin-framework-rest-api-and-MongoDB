package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevin/webgin/database"
	"github.com/kevin/webgin/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var AlbCollection *mongo.Collection = database.OpenCollection(database.Client, "albums")

func GetAlbums(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cursor, err := AlbCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Error finding albums:", err)
		return
	}
	defer cursor.Close(ctx)

	var albums []models.Album
	for cursor.Next(ctx) {
		var album models.Album
		if err := cursor.Decode(&album); err != nil {
			log.Println("Error decoding album:", err)
			continue
		}

		albums = append(albums, album)
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func PostAlbums(c *gin.Context) {
	var newAlbum models.Album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	newAlbum.ID = primitive.NewObjectID()

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	_, err := AlbCollection.InsertOne(ctx, newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Error inserting album:", err)
		return
	}

}

func GetAlbumByID(c *gin.Context) {
	albumID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid album id"})
		return
	}

	var album models.Album

	err = AlbCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&album)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})

		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			log.Println("Error finding album:", err)
		}
		return
	}
	c.IndentedJSON(http.StatusOK, album)

}

func UpdateAlbum(c *gin.Context) {
	albumID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updatedAlbum models.Album
	if err := c.ShouldBindJSON(&updatedAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid data"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updatedAlbum}

	_, err = AlbCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Error updating album:", err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Album updated"})

}

func DeleteAlbum(c *gin.Context) {
	albumID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	filter := bson.M{"_id": objID}
	result, err := AlbCollection.DeleteOne(ctx, filter)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Error deleting album:", err)
		return
	}
	if result.DeletedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Album deleted"})

}
