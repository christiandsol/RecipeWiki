package controller

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func CreateTables(conn *pgx.Conn) error {
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

func QueryRecipes(conn *pgx.Conn) ([]Recipe, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, name, description FROM recipes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err := rows.Scan(&r.RecipeID, &r.Name, &r.Description)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, nil
}

func InsertRecipe(conn *pgx.Conn, r Recipe) (int, error) {
	// Insert recipe first, get the generated ID back
	var recipeID int
	err := conn.QueryRow(context.Background(),
		`INSERT INTO recipes (name, description, steps, info)
         VALUES ($1, $2, $3, $4) RETURNING id`,
		r.Name, r.Description, r.Steps, r.Info,
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
