package controller

import (
	"errors"
	"fmt"
)

/* ============DEBUG============*/
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

func printRecipes(recipes []Recipe) {
	fmt.Println("=======PRINTING RECIPES========")
	for _, recipe := range recipes {
		printRecipe(recipe)
	}
	fmt.Println("===============")
}

/* ============Generator============*/
func (s *Store) GenerateID(name string) (int, error) {
	_, ok := s.NameID[name]
	if !ok {
		s.Mu.Lock()
		defer s.Mu.Unlock()
		new_id := s.NextID
		s.NameID[name] = new_id
		s.NextID++
		return new_id, nil
	} else {
		return -1, errors.New("recipe name already exist")
	}
}

func (s *Store) FindRecipeByID(RecipeID int) (*Recipe, int, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var ret *Recipe
	var idx int = -1
	for i := range s.Recipes {
		if RecipeID == s.Recipes[i].RecipeID {
			ret = &s.Recipes[i]
			idx = i
		}
	}

	if ret == nil {
		return ret, -1, fmt.Errorf("Unable to find recipe with ID %v", RecipeID)
	} else {
		return ret, idx, nil
	}
}

func (s *Store) FindIngredientByName(name string, ingredients []Ingredient) (int, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var idx int = -1
	for i := range ingredients {
		if ingredients[i].Name == name {
			return i, nil
		}
	}
	return idx, fmt.Errorf("Unable to find ingredients with name: %v", name)
}
