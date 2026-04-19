package controller

import (
	"fmt"
	"github.com/christiandsol/main/model"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

/* ============DEBUG============*/
func printIngredient(i model.Ingredient) {
	fmt.Printf("name: %v\n", i.Name)
	fmt.Printf("amount: %v\n", i.Amount)
	fmt.Printf("Specifier: %v\n", i.Specifier)
}

func printIngredients(recipe model.Recipe) {
	fmt.Println("PRINTING ALL INGREDIENTS")
	for _, i := range recipe.Ingredients {
		printIngredient(i)
	}
}
func printRecipe(recipe model.Recipe) {
	fmt.Printf("Name: %v\n", recipe.Name)
	fmt.Printf("RecipeID: %v\n", recipe.RecipeID)
	fmt.Printf("Description: %v\n", recipe.Description)
	fmt.Printf("Ingredients: %v\n", recipe.Ingredients)
	fmt.Printf("Steps: %v\n", recipe.Steps)
	fmt.Printf("Info: %v\n", recipe.Info)
	fmt.Println()
}

func printRecipes(recipes []model.Recipe) {
	fmt.Println("=======PRINTING RECIPES========")
	for _, recipe := range recipes {
		printRecipe(recipe)
	}
	fmt.Println("===============")
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
