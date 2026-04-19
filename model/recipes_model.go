package model

type Recipe struct {
	RecipeID    int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	ImagePath   string       `json:"image_url"`
	Steps       []string     `json:"steps"`
	Info        []string     `json:"info"`
	Ingredients []Ingredient `json:"ingredients"`
}
