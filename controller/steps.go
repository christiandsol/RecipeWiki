package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (g *Global) AddStep(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Adding Step")
	type StepJSON struct {
		RecipeID int    `json:"recipe_id"`
		StepText string `json:"step_text"`
	}
	var stepJSON StepJSON
	bytesRead, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[ERROR] Unable to read bytes %v: %v", string(bytesRead), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[ERROR] Error reading body"))
		return
	}
	err = json.Unmarshal(bytesRead, &stepJSON)
	if err != nil {
		fmt.Printf("[ERROR] Unable to unamarshal bytes %v: %v", string(bytesRead), err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[ERROR] Check JSON"))
		return
	}

	step_id, err := InsertStep(g.Conn, stepJSON.RecipeID, stepJSON.StepText)
	if err != nil {
		fmt.Printf("[ERROR] Unable to insert step: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[ERROR] Unable to insert step"))
		return
	}
	data, err := json.Marshal(map[string]interface{}{
		"step_id": step_id,
	})
	if err != nil {
		fmt.Printf("[ERROR] Unable to Marshal step")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[ERROR] Unable to marshal step"))
		return
	}

	_, err = w.Write(data)
	if err != nil {
		fmt.Printf("[ERROR] Unable to Write step")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[ERROR] Unable to write step"))
		return
	}
}

func (g *Global) UpdateStep(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Updating Step")
	type StepJSON struct {
		StepID   int    `json:"id"`
		StepText string `json:"text"`
	}
	var stepJSON StepJSON
	bytesRead, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[ERROR] Unable to read bytes %v: %v", string(bytesRead), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[ERROR] Error reading body"))
		return
	}
	err = json.Unmarshal(bytesRead, &stepJSON)
	if err != nil {
		fmt.Printf("[ERROR] Unable to unamarshal bytes %v: %v", string(bytesRead), err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[ERROR] Check JSON"))
		return
	}
	err = UpdateStepText(g.Conn, stepJSON.StepID, stepJSON.StepText)
	if err != nil {
		fmt.Printf("[ERROR] Unable to update step text\n")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[ERROR] Check JSON"))
		return
	}
	fmt.Println("SUCCESS")
	w.WriteHeader(http.StatusOK)
}

func (g *Global) GetSteps(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GETTING STEPS")
	idStr := r.PathValue("id")
	recipe_id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	steps, err := FindStepsByRecipeID(g.Conn, recipe_id)
	if err != nil {
		fmt.Printf("[ERROR] Unable to Find Steps by Recipe id %v", recipe_id)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to find steps for recipe of given id"))
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"steps": steps,
	})
	if err != nil {
		fmt.Println("[ERROR] Unable to marshal Response")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to Marshal response"))
		return
	}
	fmt.Println("Printing stringified response")
	fmt.Println(string(data))
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(fmt.Sprintf("[ERROR] Unable to write %v\n", string(data)))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to write response"))
		return
	}
}

func (g *Global) DeleteStep(w http.ResponseWriter, r *http.Request) {
	bytesRead, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[ERROR] Error reading bytes %v: %v\n", string(bytesRead), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type Request struct {
		ID int `json:"id"`
	}
	var req Request
	err = json.Unmarshal(bytesRead, &req)
	if err != nil {
		fmt.Printf("[ERROR] Error Unmarshaling bytesRead %v: %v\n", string(bytesRead), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = DeleteStep(g.Conn, req.ID)
	if err != nil {
		fmt.Printf("[ERROR] Error deleting step %v: %v\n", req.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
