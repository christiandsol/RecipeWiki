package main

import (
	"net/http"

	c "github.com/christiandsol/main/controller"
	"github.com/christiandsol/main/errUtil"
)

const (
	AMOUNT = iota
	SERVINGS
	POUNDS
	OZS
	GRAMS
	CUPS
)

// =========================GLOBALS==============================
var store = &c.Store{}

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
}

func CorsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			addCorsHeader(res)
			next.ServeHTTP(res, req)
		},
	)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.Home)
	mux.HandleFunc("POST /ingredient", store.AddIngredient)
	mux.HandleFunc("POST /recipe", store.AddRecipe)
	mux.HandleFunc("GET /recipe", store.GetRecipe)
	err := http.ListenAndServe("0.0.0.0:8080", CorsHandler(mux))
	errUtil.CheckErr("Error Starting server", nil, err)
}
