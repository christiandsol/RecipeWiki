package controller

import (
	"encoding/json"
	"fmt"
	"github.com/christiandsol/main/errUtil"
	"io"
	"net/http"
)

type Response struct {
	Ingredients []Ingredient `json:"ingredients"`
}

func (s *Store) GetIngredients(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID int `json:"id"`
	}
	bytesRead, err := io.ReadAll(r.Body)
	fmt.Printf("Bytes Read: %v\n", string(bytesRead))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[ERROR] Unable to read body"))
		return
	}
	err = json.Unmarshal(bytesRead, &body)
	if err != nil {
		fmt.Println("RETURNING 1")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[ERROR] Invalid json sent"))
		return
	}
	if len(s.Recipes[body.ID].Ingredients) > 0 {
		fmt.Printf("First ingredient for the body id: %v\n", s.Recipes[body.ID].Ingredients[0].Name)
	}
	recipe, err := s.FindRecipeByID(body.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to find recipe by id"))
		return
	}
	response := Response{
		Ingredients: recipe.Ingredients,
	}
	res, err := json.Marshal(response)
	if err != nil {
		fmt.Println("RETURNING 2")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to Marshal response"))
		return
	}
	fmt.Printf("res: %v\n", string(res))
	_, err = w.Write(res)
	if err != nil {
		fmt.Println("RETURNING 3")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to write response"))
		return
	}
}

func (s *Store) AddIngredient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] CREATING NEW INGREDIENT")
	defer r.Body.Close()
	bytesRead, err := io.ReadAll(r.Body)
	errUtil.CheckErr("Error reading", bytesRead, err)
	fmt.Println(string(bytesRead))

	var ingredient Ingredient
	err = json.Unmarshal(bytesRead, &ingredient)
	errUtil.CheckErr("Error Unmarshalling", nil, err)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Check that JSON is valid and RecipeID is filled")
		return
	}

	recipe, err := s.FindRecipeByID(ingredient.RecipeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.Mu.Lock()
	recipe.Ingredients = append(recipe.Ingredients, ingredient)
	s.Mu.Unlock()
	printIngredients(*recipe)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully added ingredient"))
}
