package controller

import (
	"sync"
)

func NewStore() *Store {
	return &Store{
		NameID: make(map[string]int),
	}
}

type Store struct {
	Mu      sync.Mutex
	Recipes []Recipe
	NextID  int
	NameID  map[string]int
}

type Ingredient struct {
	RecipeID  int    `json:"id"` // ID ingredient belongs to
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	Specifier string `json:"specifier"`
}

type Recipe struct {
	Name        string       `json:"name"`
	RecipeID    int          `json:"id"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Info        []string     `json:"info"`
}
