package main

import (
	"bufio"
	"context"
	"fmt"
	c "github.com/christiandsol/main/controller"
	"github.com/christiandsol/main/errUtil"
	"github.com/jackc/pgx/v5"
	"net/http"
	"os"
)

type Env struct {
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_PORT     string
	POSTGRES_DATA     string
	POSTGRES_VERSION  string
}

// =========================GLOBALS==============================
var store = c.NewStore()

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
}

func CorsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			addCorsHeader(res)
			next.ServeHTTP(res, req)
		},
	)
}

func isLocal(args []string) bool {
	// fmt.Printf("OS args[0]: %v", args)
	hasArg := false
	for _, arg := range args {
		if arg == "--help" {
			fmt.Printf("\t--help: get information\n")
			fmt.Printf("\t--local: develop local (default)\n")
			fmt.Printf("\t--cloud: run cloud\n")
			hasArg = true
		} else if arg == "--cloud" {
			fmt.Printf("\t Running via cloud")
			hasArg = true
			return false
		}
	}
	if hasArg == false {
		fmt.Printf("\t--help: get information\n")
		fmt.Println()
	}
	return true
}

func parseToEqual(str string) (string, string) {
	key := ""
	value := ""
	hitEquals := false
	for i := range str {
		if string(str[i]) == "=" {
			hitEquals = true
			continue
		}
		if !hitEquals {
			key += string(str[i])
		} else {
			value += string(str[i])
		}
	}
	return key, value
}

func main() {
	var host string
	if isLocal(os.Args) {
		host = "localhost"
	} else {
		host = "db"
	}
	file, err := os.Open(".env")
	if err != nil {
		fmt.Printf("Error opening file .env, error: %v", err)
		return
	}

	scanner := bufio.NewScanner(file)
	var envs = Env{}

	for scanner.Scan() {
		key, value := parseToEqual(scanner.Text())
		switch key {
		case "POSTGRES_PASSWORD":
			envs.POSTGRES_PASSWORD = value
		case "POSTGRES_USER":
			envs.POSTGRES_USER = value
		case "POSTGRES_DB":
			envs.POSTGRES_DB = value
		case "POSTGRES_PORT":
			envs.POSTGRES_PORT = value
		case "POSTGRES_DATA":
			envs.POSTGRES_DATA = value
		case "POSTGRES_VERSION":
			envs.POSTGRES_VERSION = value
		default:
			fmt.Println("Unknown input")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[ERROR] Scanner error: %v", err)
	}

	URI := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		envs.POSTGRES_USER,
		envs.POSTGRES_PASSWORD,
		host,
		envs.POSTGRES_PORT,
		envs.POSTGRES_DB,
	)
	fmt.Printf("Attemping to connect to URI: %v\n", URI)
	conn, err := pgx.Connect(context.Background(), URI)
	if err != nil {
		fmt.Printf("Unable to connect:%v", err)
		return
	}

	global := c.Global{
		Conn: conn,
	}

	defer conn.Close(context.Background())
	fmt.Println("Connected!")

	err = file.Close()
	if err != nil {
		fmt.Printf("Error closing file %v", err)
		return
	}

	err = c.CreateTables(conn)
	if err != nil {
		fmt.Printf("Unable to create tables %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /ingredients", global.GetIngredients)
	mux.HandleFunc("POST /ingredient", global.AddIngredient)
	mux.HandleFunc("DELETE /ingredient", global.DeleteIngredient)
	mux.HandleFunc("PUT /ingredient", global.UpdateIngredient)
	mux.HandleFunc("GET /recipes", global.GetRecipes)
	mux.HandleFunc("POST /recipe", global.AddRecipe)
	mux.HandleFunc("DELETE /recipe", global.DeleteRecipe)
	mux.Handle("/", http.FileServer(http.Dir("./frontend/build")))
	// mux.HandleFunc("GET /recipe", store.GetRecipe)
	err = http.ListenAndServe("0.0.0.0:8080", CorsHandler(mux))
	errUtil.CheckErr("Error Starting server", nil, err)
}
