package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"

	c "github.com/christiandsol/main/controller"
	"github.com/christiandsol/main/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
		case "IMAGE_DIR":
			envs.IMAGE_DIR = value
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
	pool, err := pgxpool.New(context.Background(), URI)
	if err != nil {
		fmt.Printf("Unable to connect:%v", err)
		return
	}
	defer pool.Close()
	imgDirExists, err := exists(envs.IMAGE_DIR)
	if err != nil {
		fmt.Printf("[ERROR] Error opening image: %v", err)
		return
	}
	if !imgDirExists {
		fmt.Printf("Image directory doesn't exist, creating %v ...\n", envs.IMAGE_DIR)
		err := os.MkdirAll(envs.IMAGE_DIR, 0755)
		if err != nil {
			fmt.Printf("[ERROR] Error making directory: %v", err)
			return
		}
	}
	fmt.Println("Successfully Created Image directory")
	global := c.Global{
		Conn:   pool,
		ImgDir: envs.IMAGE_DIR,
	}
	defer pool.Close()
	fmt.Println("Connected!")
	err = file.Close()
	if err != nil {
		fmt.Printf("Error closing file %v", err)
		return
	}
	err = db.RunMigrations(pool)
	if err != nil {
		fmt.Println("Unable to migrate:", err)
		return
	}
	if err != nil {
		fmt.Printf("Unable to migrate %v", err)
	}
	mux := http.NewServeMux()
	// ingredients
	mux.HandleFunc("POST /api/ingredients", global.GetIngredients)
	mux.HandleFunc("POST /api/ingredient", global.AddIngredient)
	mux.HandleFunc("DELETE /api/ingredient", global.DeleteIngredient)
	mux.HandleFunc("PATCH /api/ingredient", global.UpdateIngredient)
	// recipes
	mux.HandleFunc("GET /api/recipes", global.GetRecipes)
	mux.HandleFunc("POST /api/recipe", global.AddRecipe)
	mux.HandleFunc("PATCH /api/recipe", global.UpdateRecipe)
	mux.HandleFunc("DELETE /api/recipe", global.DeleteRecipe)
	mux.HandleFunc("GET /api/recipe/{id}", global.GetRecipe)
	// steps
	mux.HandleFunc("GET /api/steps/{id}", global.GetSteps)
	mux.HandleFunc("POST /api/step", global.AddStep)
	mux.HandleFunc("DELETE /api/step", global.DeleteStep)
	mux.HandleFunc("PATCH /api/step/reorder", global.ReorderStep)
	mux.HandleFunc("PATCH /api/step", global.UpdateStep)
	// fridge
	mux.HandleFunc("GET /api/fridge", global.GetFridge)
	// Frontend fallback
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "./frontend/build" + r.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, "./frontend/build/index.html")
			return
		}
		http.FileServer(http.Dir("./frontend/build")).ServeHTTP(w, r)
	}))
	// Image server
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(global.ImgDir))))
	err = http.ListenAndServe("0.0.0.0:8080", CorsHandler(mux))
	if err != nil {
		fmt.Printf("[ERROR] Error starting server %v", err)
	}
}
