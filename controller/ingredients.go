package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/christiandsol/main/errUtil"
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
	ingredients, err := FindIngredients(g.Conn, body.ID)
	if err != nil {
		fmt.Printf("[ERROR] Unable to find recipe by id, err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to find recipe by id"))
		return
	}
	response := Response{
		Ingredients: ingredients,
	}
	res, err := json.Marshal(response)
	if err != nil {
		fmt.Println("[ERROR] Unable to marshal Response")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to Marshal response"))
		return
	}
	fmt.Println("Printing stringified response")
	fmt.Println(string(res))
	_, err = w.Write(res)
	if err != nil {
		fmt.Println(fmt.Sprintf("[ERROR] Unable to write %v\n", string(res)))
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

	fmt.Printf("Printing ingredient: %v", ingredient)

	ingredientId, err := InsertIngredient(g.Conn, ingredient)
	if err != nil {
		fmt.Printf("[ERROR] Error inserting ingredient\n")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to insert ingredient with recipe id: %v\n",
			ingredient.RecipeID)))
		return
	}
	fmt.Printf("Ingredient id generated is: %v", ingredientId)
	w.WriteHeader(http.StatusOK)
	msg, err := json.Marshal(map[string]any{"id": ingredientId})
	if err != nil {
		fmt.Printf("[ERROR] error marshalling ingredientId")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error martialing from server side"))
		return
	}
	w.Write(msg)
}

func (g *Global) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] DELETING INGREDIENT")
	defer r.Body.Close()
	bytesRead, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[ERROR] Unable to read body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Printing bytes gotten from delete pathway")
	fmt.Println(string(bytesRead))

	var delIng DeleteIngredient
	err = json.Unmarshal(bytesRead, &delIng)
	if err != nil {
		fmt.Printf("[ERROR] Unable to unmarshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	fmt.Printf("SUCCESSFULLY DELETED INGREDIENT")
}

func (g *Global) UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Updating INGREDIENT")
	defer r.Body.Close()
	bytesRead, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[ERROR] Unable to read body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ingredient Ingredient
	err = json.Unmarshal(bytesRead, &ingredient)
	errUtil.CheckErr("Error Unmarshalling", nil, err)
	if err != nil {
		fmt.Printf("[ERROR] Invalid JSON for %v", string(bytesRead))
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Check that JSON is valid and RecipeID is filled")
		return
	}

	err = UpdateIngredient(g.Conn, ingredient)
	if err != nil {
		fmt.Printf("[ERROR] Unable to update ingredient from server side: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Unable to upate ingredient")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated ingredient"))
	fmt.Println("DELETED INGREDIENT")
}
