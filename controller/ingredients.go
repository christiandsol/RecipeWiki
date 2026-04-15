package controller

import (
	"encoding/json"
	"fmt"
	"github.com/christiandsol/main/errUtil"
	"io"
	"net/http"
	"slices"
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
	recipe, _, err := s.FindRecipeByID(body.ID)
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
	fmt.Println("[DEBUG] ADDING INGREDIENT ")
	defer r.Body.Close()
	bytesRead, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[ERROR] Error reading body, err: %v \n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to read request body for request: %v",
			string(bytesRead))))
		return
	}

	var ingredient Ingredient
	err = json.Unmarshal(bytesRead, &ingredient)
	if err != nil {
		fmt.Printf("[ERROR] Error unmarshalling, bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Check that JSON is valid and RecipeID is filled\n")
		return
	}

	recipe, _, err := s.FindRecipeByID(ingredient.RecipeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to find recipe with id: %v\n",
			ingredient.RecipeID)))
		return
	}

	s.Mu.Lock()
	recipe.Ingredients = append(recipe.Ingredients, ingredient)
	s.Mu.Unlock()
	printIngredients(*recipe)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully added ingredient"))
}

func (s *Store) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] DELETING INGREDIENT")
	defer r.Body.Close()
	bytesRead, err := io.ReadAll(r.Body)
	errUtil.CheckErr("Error reading", bytesRead, err)
	fmt.Println("Printing bytes gotten from delete pathway")
	fmt.Println(string(bytesRead))

	var delIng DeleteIngredient
	err = json.Unmarshal(bytesRead, &delIng)
	errUtil.CheckErr("Error Unmarshalling", nil, err)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Check that JSON is valid and RecipeID is filled")
		return
	}

	recipe, _, err := s.FindRecipeByID(delIng.RecipeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	idx, err := s.FindIngredientByName(delIng.Name, recipe.Ingredients)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Couldn't find ingredient by that name"))
		return
	}
	fmt.Printf("Attemping to delete ingredient with recipe id: %v at index %v\n", delIng.RecipeID, idx)
	recipe.Ingredients = slices.Delete(recipe.Ingredients, idx, idx+1)
	printIngredients(*recipe)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted ingredient"))
}

func (s *Store) UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Updating INGREDIENT")
	defer r.Body.Close()
	bytesRead, err := io.ReadAll(r.Body)
	errUtil.CheckErr("Error reading", bytesRead, err)
	fmt.Println("Printing bytes gotten from update pathway")
	fmt.Println(string(bytesRead))

	var upIng UpdateIngredient
	err = json.Unmarshal(bytesRead, &upIng)
	errUtil.CheckErr("Error Unmarshalling", nil, err)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Check that JSON is valid and RecipeID is filled")
		return
	}

	recipe, _, err := s.FindRecipeByID(upIng.RecipeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	idx, err := s.FindIngredientByName(upIng.PrevName, recipe.Ingredients)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Couldn't find ingredient by name: %v", upIng.PrevName)))
		return
	}
	fmt.Printf("Attemping to update ingredient with recipe id: %v at index %v\n", upIng.RecipeID, idx)
	recipe.Ingredients[idx] = upIng.Ingredient
	printIngredients(*recipe)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated ingredient"))
}
