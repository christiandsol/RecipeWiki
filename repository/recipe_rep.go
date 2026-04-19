package repository

import (
	"context"
	"github.com/christiandsol/main/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FindRecipeByID(conn *pgxpool.Pool, id int) (*model.Recipe, error) {
	var r model.Recipe
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
		var ing model.Ingredient
		err := rows.Scan(&ing.Name, &ing.Amount, &ing.Specifier)
		if err != nil {
			return nil, err
		}
		r.Ingredients = append(r.Ingredients, ing)
	}

	return &r, nil
}

func QueryRecipes(conn *pgxpool.Pool) ([]model.Recipe, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, name, description, image_path FROM recipes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []model.Recipe
	for rows.Next() {
		var r model.Recipe
		err := rows.Scan(&r.RecipeID, &r.Name, &r.Description, &r.ImagePath)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, nil
}

func InsertRecipe(conn *pgxpool.Pool, r model.Recipe) (int, error) {
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

func RemoveRecipe(conn *pgxpool.Pool, recipeId int) (string, error) {
	var fileName string

	err := conn.QueryRow(
		context.Background(),
		`DELETE FROM recipes
		 WHERE id = $1
		 RETURNING image_path`,
		recipeId,
	).Scan(&fileName)

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func PatchRecipe(conn *pgxpool.Pool, r model.Recipe) error {
	_, err := conn.Exec(context.Background(), `
		UPDATE recipes SET
			name = $1,
			description = $2,
			image_path = CASE WHEN $3 = '' THEN image_path ELSE $3 END
		WHERE id = $4
    `, r.Name, r.Description, r.ImagePath, r.RecipeID)
	return err
}
