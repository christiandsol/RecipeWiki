package model

type Step struct {
	StepID     int     `json:"step_id"`
	RecipeID   int     `json:"recipe_id"`
	StepNumber float64 `json:"step_number"`
	StepText   string  `json:"step_text"`
}
