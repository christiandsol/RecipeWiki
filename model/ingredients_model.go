package model

type Ingredient struct {
	RecipeID      int            `json:"id"` // ID ingredient belongs to
	IngredientId  int            `json:"ingredient_id"`
	Name          string         `json:"name"`
	Amount        int            `json:"amount"`
	Specifier     string         `json:"specifier"`
	CurrentAmount InventoryLevel `json:"current_amount"` // "high", "medium", "low", "out"
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
