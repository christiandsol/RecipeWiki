package repository

import (
	"context"
	"fmt"
	"github.com/christiandsol/main/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FindIngredients(conn *pgxpool.Pool, recipeId int) ([]model.Ingredient, error) {
	var ingredients []model.Ingredient
	rows, err := conn.Query(context.Background(),
		`SELECT id, recipe_id, name, amount, specifier, current_amount FROM ingredients WHERE recipe_id = $1`, recipeId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i model.Ingredient
		err := rows.Scan(&i.IngredientId, &i.RecipeID, &i.Name, &i.Amount, &i.Specifier, &i.CurrentAmount)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}
	return ingredients, nil
}

func InsertIngredient(conn *pgxpool.Pool, i model.Ingredient) (int, error) {
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

func UpdateIngredient(conn *pgxpool.Pool, i model.Ingredient) error {
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

func QueryIngredients(conn *pgxpool.Pool) ([]model.Ingredient, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, name, amount, specifier FROM ingredients`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []model.Ingredient
	for rows.Next() {
		var i model.Ingredient
		err := rows.Scan(&i.RecipeID, &i.Name, &i.Amount, &i.Specifier)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}

	return ingredients, nil
}
