package main

import (
	"net/http"
)

type Env struct {
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_PORT     string
	POSTGRES_DATA     string
	POSTGRES_VERSION  string
	IMAGE_DIR         string
}

// =========================GLOBALS==============================

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PATCH")
	headers.Set("Access-Control-Allow-Headers", "Content-Type")

}

func CorsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			addCorsHeader(res)
			if req.Method == http.MethodOptions {
				res.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(res, req)
		},
	)
}
