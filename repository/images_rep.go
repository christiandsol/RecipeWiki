package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewImage(conn *pgxpool.Pool, recipeId int) error {
	_, err := conn.Exec(context.Background(),
		`INSERT INTO recipes `)
	if err != nil {
		return err
	}
	return nil
}
