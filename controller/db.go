package controller

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func FindRecipeByID(conn *pgx.Conn, id int) (*Recipe, error) {
	var r Recipe
	err := conn.QueryRow(context.Background(),
		`SELECT id, name, description, steps, info FROM recipes WHERE id = $1`, id,
	).Scan(&r.RecipeID, &r.Name, &r.Description, &r.Steps, &r.Info)
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

func InsertIngredient(conn *pgx.Conn, i Ingredient) error {
	fmt.Println(i)
	_, err := conn.Exec(context.Background(),
		`INSERT INTO ingredients (recipe_id, name, amount, specifier)
             VALUES ($1, $2, $3, $4)`, i.RecipeID, i.Name, i.Amount, i.Specifier)
	if err != nil {
		return err
	}
	return nil
}

func RemoveIngredient(conn *pgx.Conn, id int) error {
	_, err := conn.Exec(context.Background(),
		`DELETE FROM ingredients WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateIngredient(conn *pgx.Conn, i Ingredient) error {
	_, err := conn.Exec(context.Background(), `
		UPDATE ingredients
		SET recipe_id       = $2,
				name            = $3,
				amount          = $4,
				specifier       = $5,
				current_amount  = $6,
		WHERE id = $1
		`, i.RecipeID, i.Name, i.Amount, i.Specifier, i.CurrentAmount, i.IngredientId,
	)
	if err != nil {
		return err
	}
	return nil
}

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

func QueryIngredients(conn *pgx.Conn) ([]Ingredient, error) {
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
