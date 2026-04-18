package controller

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FindIngredients(conn *pgxpool.Pool, recipeId int) ([]Ingredient, error) {
	var ingredients []Ingredient
	rows, err := conn.Query(context.Background(),
		`SELECT id, recipe_id, name, amount, specifier, current_amount FROM ingredients WHERE recipe_id = $1`, recipeId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i Ingredient
		err := rows.Scan(&i.IngredientId, &i.RecipeID, &i.Name, &i.Amount, &i.Specifier, &i.CurrentAmount)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}
	return ingredients, nil
}
func FindRecipeByID(conn *pgxpool.Pool, id int) (*Recipe, error) {
	var r Recipe
	err := conn.QueryRow(context.Background(),
		`SELECT id, name, description, steps, info, image_path FROM recipes WHERE id = $1`, id,
	).Scan(&r.RecipeID, &r.Name, &r.Description, &r.Steps, &r.Info, &r.ImagePath)
	if err != nil {
		return nil, err
	}

	// Fetch ingredients for recipe
	rows, err := conn.Query(context.Background(),
		`SELECT name, amount, specifier FROM ingredients WHERE recipe_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ing Ingredient
		err := rows.Scan(&ing.Name, &ing.Amount, &ing.Specifier)
		if err != nil {
			return nil, err
		}
		r.Ingredients = append(r.Ingredients, ing)
	}

	return &r, nil
}

func InsertIngredient(conn *pgxpool.Pool, i Ingredient) (int, error) {
	fmt.Println(i)
	var ingredientID int
	err := conn.QueryRow(context.Background(),
		`INSERT INTO ingredients (recipe_id, name, amount, specifier)
             VALUES ($1, $2, $3, $4) RETURNING id`,
		i.RecipeID, i.Name, i.Amount, i.Specifier).Scan(&ingredientID)
	if err != nil {
		return -1, err
	}
	return ingredientID, nil
}

func RemoveIngredient(conn *pgxpool.Pool, id int) error {
	_, err := conn.Exec(context.Background(),
		`DELETE FROM ingredients WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateIngredient(conn *pgxpool.Pool, i Ingredient) error {
	_, err := conn.Exec(context.Background(), `
		UPDATE ingredients
		SET recipe_id      = $2,
			name           = $3,
			amount         = $4,
			specifier      = $5,
			current_amount = $6
		WHERE id = $1
	`,
		i.IngredientId,
		i.RecipeID,
		i.Name,
		i.Amount,
		i.Specifier,
		i.CurrentAmount,
	)

	return err
}

func CreateTables(conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS recipes (
            id          SERIAL PRIMARY KEY,
            name        TEXT NOT NULL,
            description TEXT,
            steps       TEXT[],
            info        TEXT[]
        );

        CREATE TABLE IF NOT EXISTS ingredients (
            id         SERIAL PRIMARY KEY,
            recipe_id  INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
            name       TEXT NOT NULL,
            amount     INT,
            specifier  TEXT,
						current_amount TEXT
        );
    `)
	return err
}

func QueryRecipes(conn *pgxpool.Pool) ([]Recipe, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, name, description, image_path FROM recipes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err := rows.Scan(&r.RecipeID, &r.Name, &r.Description, &r.ImagePath)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, nil
}

func InsertRecipe(conn *pgxpool.Pool, r Recipe) (int, error) {
	// Insert recipe first, get the generated ID back
	var recipeID int
	err := conn.QueryRow(context.Background(),
		`INSERT INTO recipes (name, description, steps, info, image_path)
         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		r.Name, r.Description, r.Steps, r.Info, r.ImagePath,
	).Scan(&recipeID)
	if err != nil {
		return -1, err
	}

	// Insert each ingredient linked to the recipe
	for _, ing := range r.Ingredients {
		_, err := conn.Exec(context.Background(),
			`INSERT INTO ingredients (recipe_id, name, amount, specifier, current_amount)
             VALUES ($1, $2, $3, $4)`,
			recipeID, ing.Name, ing.Amount, ing.Specifier, ing.CurrentAmount,
		)
		if err != nil {
			return -1, err
		}
	}
	return recipeID, nil
}

func QueryIngredients(conn *pgxpool.Pool) ([]Ingredient, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, name, amount, specifier FROM ingredients`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []Ingredient
	for rows.Next() {
		var i Ingredient
		err := rows.Scan(&i.RecipeID, &i.Name, &i.Amount, &i.Specifier)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}

	return ingredients, nil
}

func RemoveRecipe(conn *pgxpool.Pool, recipeId int) error {
	_, err := conn.Exec(context.Background(),
		`DELETE FROM recipes WHERE id = $1`, recipeId)
	if err != nil {
		return err
	}
	return nil
}

func NewImage(conn *pgxpool.Pool, recipeId int) error {
	_, err := conn.Exec(context.Background(),
		`INSERT INTO recipes `)
	if err != nil {
		return err
	}
	return nil
}

func PatchRecipe(conn *pgxpool.Pool, r Recipe) error {
	_, err := conn.Exec(context.Background(), `
		UPDATE recipes SET
			name = $1,
			description = $2,
			image_path = CASE WHEN $3 = '' THEN image_path ELSE $3 END
		WHERE id = $4
    `, r.Name, r.Description, r.ImagePath, r.RecipeID)
	return err
}

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

func FindSteps(conn *pgxpool.Pool, recipeID int) ([]Step, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT step_id, step_text FROM steps WHERE recipe_id = $1 ORDER BY step_number`,
		recipeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []Step
	for rows.Next() {
		var s Step
		if err := rows.Scan(&s.StepID, &s.StepText); err != nil {
			return nil, err
		}
		steps = append(steps, s)
	}
	return steps, rows.Err()
}

func FindStepsByRecipeID(conn *pgxpool.Pool, recipe_id int) ([]Step, error) {
	// Fetch steps for recipe with recipe_id
	rows, err := conn.Query(context.Background(),
		`SELECT step_id, step_text, step_number FROM steps WHERE recipe_id = $1 ORDER BY step_number`,
		recipe_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var steps []Step
	for rows.Next() {
		var s Step
		if err := rows.Scan(&s.StepID, &s.StepText, &s.StepNumber); err != nil {
			return nil, err
		}
		steps = append(steps, s)
	}

	return steps, nil
}
