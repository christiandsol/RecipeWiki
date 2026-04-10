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
	id, err := s.GenerateID(recipe.Name)
	if err != nil {
		fmt.Println("[ERROR] Generating new id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`[ERROR] Generating new id, 
			recipe name may already be in use`))
		return
	}
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
	data, err := json.Marshal(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		fmt.Println("[ERROR] Unable to marshal id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`[ERROR] Unable to marshal id`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		fmt.Println("[ERROR] Unable to write response")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`[ERROR] Unable write response`))
		return
	}
	printRecipes(s.Recipes)
}

func (s *Store) GetRecipes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GETTING RECIPES")
	var recipes struct {
		Recipes []Recipe `json:"recipes"`
	}
	if s.Recipes == nil {
		recipes.Recipes = []Recipe{}
	} else {
		recipes.Recipes = s.Recipes
	}
	msg, err := json.Marshal(recipes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request, fix json"))
		return
	}
	w.Write(msg)
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
