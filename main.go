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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/exp/slices"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	json.Unmarshal([]byte(file), &recipes)
}

// swagger:operation POST /recipes recipes createRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//   '200':
// 	   description: Successful operation
//   '400':
//	   description: Invalid input
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//   '200':
//	   description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation PUT /recipes/:id recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
// - name: id
//   in: query
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//   '200':
// 	   description: Successful operation
//   '400':
//	   description: Invalid input
//   '404':
//	   description: Invalid recipe ID
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	index := slices.IndexFunc(recipes, func(item Recipe) bool {
		return item.ID == id
	})
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/:id recipes deleteRecipe
// Delete an existing recipe
// ---
// parameters:
// - name: id
//   in: query
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//   '200':
// 	   description: Successful operation
//   '404':
//	   description: Invalid recipe ID
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := slices.IndexFunc(recipes, func(item Recipe) bool {
		return item.ID == id
	})
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	recipes = slices.Delete(recipes, index, index+1)
	c.JSON(http.StatusOK, gin.H{"error": "Recipe has been deleted"})
}

// swagger:operation GET /recipes/search recipes searchRecipe
// Search recipes by tag
// ---
// parameters:
// - name: tag
//   in: query
//   description: tag to search by
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//   '200':
// 	   description: Successful operation
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)
	for _, recipe := range recipes {
		contains := slices.Contains(recipe.Tags, tag)
		if contains {
			listOfRecipes = append(listOfRecipes, recipe)
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.Run()
}
