package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CanisterRequest struct {
	Method string `json:"method"`
	Params string `json:"params"`
}

func callCanister(url string, reqBody CanisterRequest) (*http.Response, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error encoding request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	return resp, nil
}

func main() {
	canisterURL := "https://your-canister-url.com"
	req := CanisterRequest{
		Method: "get_data",
		Params: "{}",
	}

	resp, err := callCanister(canisterURL, req)
	if err != nil {
		log.Fatalf("Failed to call canister: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Response Status: %s\n", resp.Status)

	// Jika ingin menampilkan body respons
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("Failed to read response body: %v", err)
	// }
	// fmt.Println("Response Body:", string(body))
}
