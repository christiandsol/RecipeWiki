package controller

import (
	"github.com/jackc/pgx/v5"
	"sync"
)

func NewStore() *Store {
	return &Store{
		NameID: make(map[string]int),
	}
}

type Global struct {
	Conn *pgx.Conn
}

type Store struct {
	Mu               sync.Mutex
	Recipes          []Recipe
	NextID           int
	NameID           map[string]int
	NextIngredientID int
}

type Ingredient struct {
	RecipeID      int     `json:"id"` // ID ingredient belongs to
	IngredientId  int     `json:"ingredient_id"`
	Name          string  `json:"name"`
	Amount        int     `json:"amount"`
	Specifier     string  `json:"specifier"`
	CurrentAmount *string `json:"current_amount"` // "high", "medium", "low"
}

type UpdateIng struct {
	RecipeID   int        `json:"id"` // ID ingredient belongs to
	PrevName   string     `json:"prev_name"`
	Ingredient Ingredient `json:"ingredient"`
}

type DeleteIngredient struct {
	RecipeID int    `json:"recipe_id"` // ID ingredient belongs to
	IngID    int    `json:"id"`        // ID ingredient belongs to
	Name     string `json:"name"`
}

type Recipe struct {
	Name        string       `json:"name"`
	RecipeID    int          `json:"id"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Info        []string     `json:"info"`
}
