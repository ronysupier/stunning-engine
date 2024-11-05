package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CanisterRequest struct {
	Method string `json:"method"`
	Params string `json:"params"`
}

type CanisterResponse struct {
	Result string `json:"result"`
}

func callCanister(url string, reqBody CanisterRequest) (CanisterResponse, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return CanisterResponse{}, fmt.Errorf("error encoding request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return CanisterResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return CanisterResponse{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	var canisterResp CanisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&canisterResp); err != nil {
		return CanisterResponse{}, fmt.Errorf("error decoding response: %v", err)
	}

	return canisterResp, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Ganti dengan URL canister yang sesuai
	canisterURL := "https://your-canister-url.com"
	req := CanisterRequest{
		Method: "get_data",
		Params: "{}",
	}

	resp, err := callCanister(canisterURL, req)
	if err != nil {
		http.Error(w, "Failed to call canister: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, resp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)

	http.Handle("/", r)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
