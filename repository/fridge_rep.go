package repository

import (
	"context"
	"github.com/christiandsol/main/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetFridge(conn *pgxpool.Pool) ([]model.Ingredient, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, recipe_id, name, amount, specifier, current_amount
		 FROM ingredients`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []model.Ingredient
	for rows.Next() {
		var i model.Ingredient
		err := rows.Scan(
			&i.IngredientId,
			&i.RecipeID,
			&i.Name,
			&i.Amount,
			&i.Specifier,
			&i.CurrentAmount,
		)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}

	return ingredients, nil
}
