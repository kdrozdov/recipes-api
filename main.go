// Recipes API
//
// Description
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
// Produces:
// - application/json
// swagger:meta
package main

import (
	"context"
	"log"
	"recipes-api/config"
	"recipes-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	config.LoadConfig("./config")

	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("db.uri")))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal()
	}
	log.Println("Connected to mongodb")
	collection := client.Database(viper.GetString("db.name")).Collection("recipes")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	log.Println("Redis", redisClient.Ping(ctx))

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipe)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipe)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipe)
	router.GET("/recipes/search", recipesHandler.SearchRecipes)
	router.GET("/recipes", recipesHandler.ListRecipes)
	router.Run()
}
