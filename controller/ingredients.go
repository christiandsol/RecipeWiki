package controller

import (
	"encoding/json"
	"fmt"
	"github.com/christiandsol/main/errUtil"
	"io"
	"net/http"
)

func (s *Store) AddIngredient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] CREATING NEW INGREDIENT")
	defer r.Body.Close()
	bytesRead, err := io.ReadAll(r.Body)
	errUtil.CheckErr("Error reading", bytesRead, err)

	var ingredient Ingredient
	err = json.Unmarshal(bytesRead, &ingredient)
	errUtil.CheckErr("Error Unmarshalling", nil, err)
	if err != nil {
		io.WriteString(w, "Check that JSON is valid and RecipeID is filled")
		return
	}

	recipe, err := s.FindRecipeByID(ingredient.RecipeID)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	s.Mu.Lock()
	recipe.Ingredients = append(recipe.Ingredients, ingredient)
	s.Mu.Unlock()
	printIngredients(*recipe)
}
