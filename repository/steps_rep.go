package repository

import (
	"context"
	"github.com/christiandsol/main/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertStep(conn *pgxpool.Pool, recipeID int, text string) (int, float64, error) {
	var stepID int
	var stepNumber float64
	err := conn.QueryRow(context.Background(),
		`INSERT INTO steps (recipe_id, step_number, step_text)
         VALUES ($1, COALESCE((SELECT MAX(step_number) FROM steps WHERE recipe_id = $1), 0) + 1, $2)
         RETURNING step_id, step_number`,
		recipeID, text,
	).Scan(&stepID, &stepNumber)
	return stepID, stepNumber, err
}

func UpdateStepText(conn *pgxpool.Pool, stepID int, text string) error {
	_, err := conn.Exec(context.Background(),
		`UPDATE steps SET step_text = $1 WHERE step_id = $2`,
		text, stepID,
	)
	return err
}

func ReorderStepDB(conn *pgxpool.Pool, stepID int, before, after float64) error {
	_, err := conn.Exec(context.Background(),
		`UPDATE steps SET step_number = $1 WHERE step_id = $2`,
		(before+after)/2, stepID,
	)
	return err
}

func DeleteStep(conn *pgxpool.Pool, stepID int) error {
	_, err := conn.Exec(context.Background(),
		`DELETE FROM steps WHERE step_id = $1`,
		stepID,
	)
	return err
}

func FindSteps(conn *pgxpool.Pool, recipeID int) ([]model.Step, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT step_id, step_text FROM steps WHERE recipe_id = $1 ORDER BY step_number`,
		recipeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []model.Step
	for rows.Next() {
		var s model.Step
		if err := rows.Scan(&s.StepID, &s.StepText); err != nil {
			return nil, err
		}
		steps = append(steps, s)
	}
	return steps, rows.Err()
}

func FindStepsByRecipeID(conn *pgxpool.Pool, recipe_id int) ([]model.Step, error) {
	// Fetch steps for recipe with recipe_id
	rows, err := conn.Query(context.Background(),
		`SELECT step_id, step_text, step_number FROM steps WHERE recipe_id = $1 ORDER BY step_number`,
		recipe_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var steps []model.Step
	for rows.Next() {
		var s model.Step
		if err := rows.Scan(&s.StepID, &s.StepText, &s.StepNumber); err != nil {
			return nil, err
		}
		steps = append(steps, s)
	}

	return steps, nil
}
