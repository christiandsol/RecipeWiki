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

func (g *Global) GetIngredients(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] GETTING INGREDIENTS")
	var body struct {
		ID int `json:"id"`
	}
	bytesRead, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[ERROR] Unable to read body"))
		return
	}
	err = json.Unmarshal(bytesRead, &body)
	if err != nil {
		fmt.Println("[ERROR] Error unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[ERROR] Invalid json sent"))
		return
	}
	recipe, err := FindRecipeByID(g.Conn, body.ID)
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
		fmt.Println("[ERROR] Unable to marshal Response")
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

func (g *Global) AddIngredient(w http.ResponseWriter, r *http.Request) {
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

	err = InsertIngredient(g.Conn, ingredient)
	if err != nil {
		fmt.Printf("[ERROR] Error inserting ingredient\n")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to insert ingredient with recipe id: %v\n",
			ingredient.RecipeID)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully added ingredient"))
}

func (g *Global) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
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
	err = RemoveIngredient(g.Conn, delIng.IngID)
	if err != nil {
		fmt.Printf("[ERROR] error removing ingredient %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error removing ingredient"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted ingredient"))
}

func (g *Global) UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Updating INGREDIENT")
	defer r.Body.Close()
	bytesRead, err := io.ReadAll(r.Body)
	errUtil.CheckErr("Error reading", bytesRead, err)
	fmt.Println("Printing bytes gotten from update pathway")
	fmt.Println(string(bytesRead))

	var ingredient Ingredient
	err = json.Unmarshal(bytesRead, &ingredient)
	errUtil.CheckErr("Error Unmarshalling", nil, err)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Check that JSON is valid and RecipeID is filled")
		return
	}

	err = UpdateIngredient(g.Conn, ingredient)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated ingredient"))
}
