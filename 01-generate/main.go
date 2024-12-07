package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ollama/ollama/api"
)

var (
	FALSE = false
	TRUE  = true
)

func main() {
	ctx := context.Background()

	var ollamaRawUrl string
	if ollamaRawUrl = os.Getenv("OLLAMA_HOST"); ollamaRawUrl == "" {
		ollamaRawUrl = "http://localhost:11434"
	}

	url, _ := url.Parse(ollamaRawUrl)

	client := api.NewClient(url, http.DefaultClient)

	req := &api.GenerateRequest{
		Model:  "qwen2.5:0.5b",
		Prompt: "The best pizza in the world is",
		Options: map[string]interface{}{
			"temperature":   0.8,
			"repeat_last_n": 2,
		},
		Stream: &TRUE,
	}

	err := client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	})

	if err != nil {
		log.Fatalln("😡", err)
	}
	fmt.Println()

}
