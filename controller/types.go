package controller

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewStore() *Store {
	return &Store{
		NameID: make(map[string]int),
	}
}

type Global struct {
	Conn   *pgxpool.Pool
	ImgDir string
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
	CurrentAmount *string `json:"current_amount"` // "high", "medium", "low", "out"
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
	RecipeID    int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	ImagePath   string       `json:"image_url"`
	Steps       []string     `json:"steps"`
	Info        []string     `json:"info"`
	Ingredients []Ingredient `json:"ingredients"`
}

type Step struct {
	StepID     int     `json:"step_id"`
	RecipeID   int     `json:"recipe_id"`
	StepNumber float64 `json:"step_number"`
	StepText   string  `json:"step_text"`
}
