package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (g *Global) GetRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	recipe, err := FindRecipeByID(g.Conn, id)
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
	recipes, err := QueryRecipes(g.Conn)
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

	recipe := Recipe{
		Name:        name,
		Description: description,
		ImagePath:   imageFilename,
	}

	id, err := InsertRecipe(g.Conn, recipe)
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

	err = RemoveRecipe(g.Conn, deleteRec.RecipeId)
	if err != nil {
		fmt.Printf("[ERROR] Unable to remove recipe", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`[ERROR] Unable to remove recipe`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted recipe"))
}

func (g *Global) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Update recipe")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "[ERROR] Unable to read body", http.StatusBadRequest)
		return
	}

	var recipe Recipe
	if err := json.Unmarshal(body, &recipe); err != nil {
		http.Error(w, "[ERROR] Unable to unmarshal body", http.StatusBadRequest)
		return
	}

	if err := PatchRecipe(g.Conn, recipe); err != nil {
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

func (g *Global) UploadImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Upload image")

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "[ERROR] Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "[ERROR] Missing image field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	image, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "[ERROR] Unable to read file", http.StatusBadRequest)
		return
	}
	fmt.Printf("Size of body read: %v\n", len(image))

	mimeType := http.DetectContentType(image)
	ext, ok := allowedTypes[mimeType]
	if !ok {
		// http.DetectContentType doesn't know HEIC, handle below
		ext, ok = detectHEIC(image)
		if !ok {
			fmt.Printf("[ERROR] Unsupported image type: %s", mimeType)
			http.Error(w, fmt.Sprintf("[ERROR] Unsupported image type: %s", mimeType), http.StatusUnsupportedMediaType)
			return
		}
	}

	filename := fmt.Sprintf("%v.%v", time.Now().UnixNano(), ext)
	filepath := filepath.Join(g.ImgDir, filename)
	if err := os.WriteFile(filepath, image, 0644); err != nil {
		fmt.Printf("[ERROR] Unable to save image, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[ERROR] Unable to save image, err: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"path": "%s"}`, filename)))
}
