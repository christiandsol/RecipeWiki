package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/christiandsol/main/errUtil"
)

const (
	AMOUNT = iota
	SERVINGS
	POUNDS
	OZS
	GRAMS
	CUPS
)

type Ingredient struct {
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	Specifier int    `json:"specifier"`
}

func printIngredient(i Ingredient) {
	fmt.Printf("name: %v\n", i.Name)
	fmt.Printf("amount: %v\n", i.Amount)
	fmt.Printf("Specifier: %v\n", i.Specifier)
}

func newIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		panic("[PANIC] Method Doesn't match handler")
	}
	fmt.Println("[DEBUG] CREATING NEW INGREDIENT")

	defer r.Body.Close()
	bytesRead, err := io.ReadAll(r.Body)
	errUtil.PrintErr("Error reading", bytesRead, err)

	var ingredient Ingredient
	err = json.Unmarshal(bytesRead, &ingredient)
	printIngredient(ingredient)
	errUtil.PrintErr("Error Unmarshalling", nil, err)
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("What's up"))
	errUtil.PrintErr("Writing buf", nil, err)
}

func main() {
	fmt.Println("Hello World")
	err := http.ListenAndServe("0.0.0.0:8080", http.HandlerFunc(newIngredient))
	errUtil.PrintErr("Error Starting server", nil, err)
}
