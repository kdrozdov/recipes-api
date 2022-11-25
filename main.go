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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	recipesHandler *handlers.RecipesHandler
	authHandler    *handlers.AuthHandler
)

func init() {
	// Read config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Connect to MongoDB
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DB.URI))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("can't connect to mongodb")
	}
	log.Println("connected to mongodb")
	collectionRecipes := client.Database(cfg.DB.Name).Collection("recipes")
	collectionUsers := client.Database(cfg.DB.Name).Collection("users")

	// Connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DBNumber,
	})
	log.Println("redis", redisClient.Ping(ctx))

	recipesHandler = handlers.NewRecipesHandler(ctx, collectionRecipes, redisClient)
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers, cfg)
}

func main() {
	router := gin.Default()

	router.POST("/signin", authHandler.SignInHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipes)
	router.GET("/recipes", recipesHandler.ListRecipes)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipe)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipe)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipe)
	}

	router.Run()
}
