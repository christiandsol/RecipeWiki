package controller

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Global struct {
	Conn   *pgxpool.Pool
	ImgDir string
}

func Home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("GOT HOME"))
	if err != nil {
		fmt.Printf("[ERROR] Error writing buffer: %v", err)
	}
}
