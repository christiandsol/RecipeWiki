package controller

import (
	"encoding/json"
	"fmt"
	"github.com/christiandsol/main/model"
	repo "github.com/christiandsol/main/repository"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (g *Global) GetRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	recipe, err := repo.FindRecipeByID(g.Conn, id)
	if err != nil {
		http.Error(w, "recipe not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(recipe)
	if err != nil {
		http.Error(w, "failed to marshal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (g *Global) GetRecipes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GETTING RECIPES")
	recipes, err := repo.QueryRecipes(g.Conn)
	if err != nil {
		fmt.Printf("[ERROR], error querying recipes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error querying recipes"))
		return
	}
	msg, err := json.Marshal(map[string]any{"recipes": recipes})
	if err != nil {
		fmt.Printf("[ERROR], error marshalling recipes: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request, fix json"))
		return
	}
	fmt.Printf("[DEBUG] get recipes msg: %v\n", string(msg))
	w.Write(msg)
}

func (g *Global) AddRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CREATING NEW RECIPE")

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "[ERROR] Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	// Image is optional — only save if one was attached
	var imageFilename string
	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageFilename, err = g.saveImage(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
			return
		}
	}

	recipe := model.Recipe{
		Name:        name,
		Description: description,
		ImagePath:   imageFilename,
	}

	id, err := repo.InsertRecipe(g.Conn, recipe)
	if err != nil {
		fmt.Printf("[ERROR] Unable to insert recipe, err: %v\n", err)
		http.Error(w, "[ERROR] Unable to insert recipe", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"id":        id,
		"image_url": imageFilename,
	})
	if err != nil {
		http.Error(w, "[ERROR] Unable to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (g *Global) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	bytesRead, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[ERROR] Unable to read request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`[ERROR] Unable to read request`))
		return
	}
	type DeleteRec struct {
		RecipeId int `json:"recipe_id"`
	}
	var deleteRec DeleteRec
	err = json.Unmarshal(bytesRead, &deleteRec)
	if err != nil {
		fmt.Printf("[ERROR] Unable to unmarshal bytes read: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`[ERROR] Unable to marshal bytes read`))
		return
	}

	fileName, err := repo.RemoveRecipe(g.Conn, deleteRec.RecipeId)
	if err != nil {
		fmt.Printf("[ERROR] Unable to remove recipe", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`[ERROR] Unable to remove recipe`))
		return
	}

	err = g.removeImage(fileName)
	if err != nil {
		fmt.Printf("[ERROR] Unable to remove file image", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`[ERROR] Unable to remove file image`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted recipe"))
}

func (g *Global) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Update recipe")
	var recipe model.Recipe

	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		// Image upload path
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			http.Error(w, "[ERROR] Failed to parse form", http.StatusBadRequest)
			return
		}

		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "[ERROR] Invalid id", http.StatusBadRequest)
			return
		}
		recipe.RecipeID = id
		recipe.Name = r.FormValue("name")
		recipe.Description = r.FormValue("description")

		file, _, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			filename, err := g.saveImage(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
				return
			}
			recipe.ImagePath = filename
			oldRecipe, err := repo.FindRecipeByID(g.Conn, recipe.RecipeID)
			if err != nil {
				fmt.Printf("[ERROR] Unable to find original recipe %v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = g.removeImage(oldRecipe.ImagePath)
			if err != nil {
				fmt.Printf("[ERROR] Unable to remove old image%v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	} else {
		// JSON path (name/description only)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "[ERROR] Unable to read body", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &recipe); err != nil {
			http.Error(w, "[ERROR] Unable to unmarshal body", http.StatusBadRequest)
			return
		}
	}

	if err := repo.PatchRecipe(g.Conn, recipe); err != nil {
		fmt.Printf("[ERROR] Unable to update recipe: %v\n", err)
		http.Error(w, "[ERROR] Unable to update recipe", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "updated"}`))
}

var allowedTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/heic": ".heic",
	"image/heif": ".heif", // HEIC is a container, HEIF is the format — iPhones use both
	"image/webp": ".webp", // common from web browsers
}
