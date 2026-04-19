package controller

import (
	"encoding/json"
	"fmt"
	repo "github.com/christiandsol/main/repository"
	"net/http"
)

func (g *Global) GetFridge(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] GETTING FRIDGE")
	ingredients, err := repo.GetFridge(g.Conn)
	if err != nil {
		fmt.Printf("[ERROR] Unable to get all ingredients: %v", err)
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
