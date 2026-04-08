package controller

import (
	"encoding/json"
	"fmt"
	"github.com/christiandsol/main/errUtil"
	"io"
	"net/http"
)

func (s *Store) AddRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("NEW RECIPE CALLED\n")
	bytesRead, err := io.ReadAll(r.Body)
	var recipe Recipe
	err = json.Unmarshal(bytesRead, &recipe)

	fmt.Printf("HERE\n")
	id := s.GenerateID()
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Recipes = append(s.Recipes, Recipe{
		Name:        recipe.Name,
		RecipeID:    id,
		Description: recipe.Description,
		Ingredients: recipe.Ingredients,
		Steps:       recipe.Steps,
		Info:        recipe.Info,
	})
	_, err = w.Write([]byte("SUCCESS"))
	errUtil.CheckErr("Writing buf", nil, err)
	printRecipes(s.Recipes)
}

func (s *Store) GetRecipe(w http.ResponseWriter, r *http.Request) {
	bytesRead, err := io.ReadAll(r.Body)
	errUtil.CheckErr("Error getting recipe", nil, err)
	if err == nil {
		w.Write([]byte("Error getting recipe"))
		return
	}
	var recipe Recipe
	err = json.Unmarshal(bytesRead, &recipe)
	errUtil.CheckErr("Error Marshalling recipe", nil, err)
	if err == nil {
		w.Write([]byte("Error Marshalling recipe"))
		return
	}
	json, err := json.Marshal(recipe)
	w.Write(json)
}
