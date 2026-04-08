package main

import (
	"encoding/json"
	"fmt"
	"github.com/christiandsol/main/errUtil"
	"io"
	"net/http"
	"sync"
)

const (
	AMOUNT = iota
	SERVINGS
	POUNDS
	OZS
	GRAMS
	CUPS
)

type Store struct {
	mu      sync.Mutex
	recipes []Recipe
	nextID  int
}

type Ingredient struct {
	RecipeID  int    `json:"id"` // ID ingredient belongs to
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	Specifier int    `json:"specifier"`
}

type Recipe struct {
	Name        string       `json:"name"`
	RecipeID    int          `json:"id"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Info        []string     `json:"info"`
}

// =========================GLOBALS==============================

var store = &Store{}

// var recipes []Recipe
// var NEXT_ID int = 0

// =========================HELPERS==============================
func (s *Store) GenerateID() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	new_id := s.nextID
	s.nextID++
	return new_id
}

func printIngredient(i Ingredient) {
	fmt.Printf("name: %v\n", i.Name)
	fmt.Printf("amount: %v\n", i.Amount)
	fmt.Printf("Specifier: %v\n", i.Specifier)
}

func printIngredients(recipe Recipe) {
	fmt.Println("PRINTING ALL INGREDIENTS")
	for _, i := range recipe.Ingredients {
		printIngredient(i)
	}
}

func printRecipe(recipe Recipe) {
	fmt.Printf("Name: %v\n", recipe.Name)
	fmt.Printf("RecipeID: %v\n", recipe.RecipeID)
	fmt.Printf("Description: %v\n", recipe.Description)
	fmt.Printf("Ingredients: %v\n", recipe.Ingredients)
	fmt.Printf("Steps: %v\n", recipe.Steps)
	fmt.Printf("Info: %v\n", recipe.Info)
	fmt.Println()
}

func printRecipes() {
	fmt.Println("=======PRINTING RECIPES========")
	for _, recipe := range store.recipes {
		printRecipe(recipe)
	}
	fmt.Println("===============")
}

func (s *Store) findRecipeByID(RecipeID int) (*Recipe, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var ret *Recipe
	for i := range s.recipes {
		if RecipeID == s.recipes[i].RecipeID {
			ret = &s.recipes[i]
		}
	}

	if ret == nil {
		return ret, fmt.Errorf("Unable to find recipe with ID %v", RecipeID)
	} else {
		return ret, nil
	}
}

func (s *Store) newIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		io.WriteString(w, "[RESPONSE] Method Doesn't match handler, method expects POST")
		return
	}

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

	recipe, err := s.findRecipeByID(ingredient.RecipeID)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	s.mu.Lock()
	recipe.Ingredients = append(recipe.Ingredients, ingredient)
	s.mu.Unlock()
	printIngredients(*recipe)
}

func (s *Store) newRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("NEW RECIPE CALLED\n")
	bytesRead, err := io.ReadAll(r.Body)
	var recipe Recipe
	err = json.Unmarshal(bytesRead, &recipe)

	fmt.Printf("HERE\n")
	id := s.GenerateID()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recipes = append(s.recipes, Recipe{
		Name:        recipe.Name,
		RecipeID:    id,
		Description: recipe.Description,
		Ingredients: recipe.Ingredients,
		Steps:       recipe.Steps,
		Info:        recipe.Info,
	})
	_, err = w.Write([]byte("SUCCESS"))
	errUtil.CheckErr("Writing buf", nil, err)
	printRecipes()
}

func home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("GOT HOME"))
	errUtil.CheckErr("Writing buf", nil, err)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/new_ingredient", store.newIngredient)
	mux.HandleFunc("/recipe", store.newRecipe)
	err := http.ListenAndServe("0.0.0.0:8080", mux)
	errUtil.CheckErr("Error Starting server", nil, err)
}
