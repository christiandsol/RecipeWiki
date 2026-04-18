package controller

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

/* ============DEBUG============*/
func printIngredient(i Ingredient) {
	fmt.Printf("name: %v\n", i.Name)
	fmt.Printf("amount: %v\n", i.Amount)
	fmt.Printf("Specifier: %v\n", i.Specifier)
}

func printIngredients(recipe Recipe) {
	fmt.Println("PRINTING ALL INGREDIENTS")
	for _, i := range recipe.Ingredients {
		printIngredient(i)
	}
}
func printRecipe(recipe Recipe) {
	fmt.Printf("Name: %v\n", recipe.Name)
	fmt.Printf("RecipeID: %v\n", recipe.RecipeID)
	fmt.Printf("Description: %v\n", recipe.Description)
	fmt.Printf("Ingredients: %v\n", recipe.Ingredients)
	fmt.Printf("Steps: %v\n", recipe.Steps)
	fmt.Printf("Info: %v\n", recipe.Info)
	fmt.Println()
}

func printRecipes(recipes []Recipe) {
	fmt.Println("=======PRINTING RECIPES========")
	for _, recipe := range recipes {
		printRecipe(recipe)
	}
	fmt.Println("===============")
}

/* ============Generator============*/
func (s *Store) GenerateID(name string) (int, error) {
	_, ok := s.NameID[name]
	if !ok {
		s.Mu.Lock()
		defer s.Mu.Unlock()
		new_id := s.NextID
		s.NameID[name] = new_id
		s.NextID++
		return new_id, nil
	} else {
		return -1, errors.New("recipe name already exist")
	}
}

func (s *Store) FindRecipeByID(RecipeID int) (*Recipe, int, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var ret *Recipe
	var idx int = -1
	for i := range s.Recipes {
		if RecipeID == s.Recipes[i].RecipeID {
			ret = &s.Recipes[i]
			idx = i
		}
	}

	if ret == nil {
		return ret, -1, fmt.Errorf("Unable to find recipe with ID %v", RecipeID)
	} else {
		return ret, idx, nil
	}
}

func (s *Store) FindIngredientByName(name string, ingredients []Ingredient) (int, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var idx int = -1
	for i := range ingredients {
		if ingredients[i].Name == name {
			return i, nil
		}
	}
	return idx, fmt.Errorf("Unable to find ingredients with name: %v", name)
}

// A little bit about implementation
// http.DetectContentType uses magic bytes and doesn't recognise HEIC,
// must sniff manually, HEIC files have "ftyp" at byte 4 and
// "heic"/"heix"/"hevc" etc. at byte 8.
func detectHEIC(data []byte) (string, bool) {
	if len(data) < 12 {
		return "", false
	}
	if string(data[4:8]) != "ftyp" {
		return "", false
	}
	brand := string(data[8:12])
	heicBrands := map[string]bool{
		"heic": true, "heix": true, "hevc": true,
		"hevx": true, "heim": true, "hevm": true,
		"mif1": true, // also used by HEIF
	}
	if heicBrands[brand] {
		return ".heic", true
	}
	return "", false
}

// saveImage is now a plain helper, not an HTTP handler
func (g *Global) saveImage(file multipart.File) (string, error) {
	image, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Unable to read file: %w", err)
	}

	mimeType := http.DetectContentType(image)
	ext, ok := allowedTypes[mimeType]
	if !ok {
		ext, ok = detectHEIC(image)
		if !ok {
			return "", fmt.Errorf("[ERROR] Unsupported image type: %s", mimeType)
		}
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(g.ImgDir, filename)

	if err := os.WriteFile(fullPath, image, 0644); err != nil {
		return "", fmt.Errorf("[ERROR] Unable to save image: %w", err)
	}

	return filename, nil
}

func (g *Global) removeImage(fileName string) error {
	err := os.Remove(filepath.Join(g.ImgDir, fileName))
	return err
}
