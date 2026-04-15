package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/christiandsol/main/errUtil"
)

func (g *Global) GetRecipes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GETTING RECIPES")
	recipes, err := QueryRecipes(g.Conn)
	if err != nil {
		fmt.Printf("[ERROR], error querying recipes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error querying recipes"))
		return
	}
	msg, err := json.Marshal(map[string]any{"recipes": recipes})
	if err != nil {
		fmt.Printf("[ERROR], error marshalling recipes: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request, fix json"))
		return
	}
	fmt.Printf("[DEBUG] get recipes msg: %v\n", string(msg))
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

func (g *Global) AddRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("CREATING NEW RECIPE\n")
	bytesRead, err := io.ReadAll(r.Body)
	var recipe Recipe
	err = json.Unmarshal(bytesRead, &recipe)
	if err != nil {
		fmt.Println(fmt.Sprintf("[ERROR] Error unmarshaling request with body: %v\n",
			string(bytesRead)))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("[ERROR] Error unmarshaling request with body: %v\n",
			string(bytesRead))))
		return
	}

	id, err := InsertRecipe(g.Conn, recipe)
	if err != nil {
		fmt.Println("[ERROR] Unable to insert recipe, err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`[ERROR] Unable to insert recipe`))
		return
	}

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
}
