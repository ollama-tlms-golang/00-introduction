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

	fmt.Println("🚀 Ollama API URL:", ollamaRawUrl)

	url, _ := url.Parse(ollamaRawUrl)

	client := api.NewClient(url, http.DefaultClient)

	systemInstructions := "You are a pizzaiolo, a pizza expert. Give brief and structured answers."
	question := "What is the best pizza in the world?"

	// Prompt construction
	messages := []api.Message{
		{Role: "system", Content: systemInstructions},
		{Role: "user", Content: question},
	}

	req := &api.ChatRequest{
		Model:    "qwen2.5:0.5b",
		Messages: messages,
		Options: map[string]interface{}{
			"temperature":   0.8,
			"repeat_last_n": 2,
		},
		Stream: &TRUE,
	}

	err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	})

	if err != nil {
		log.Fatalln("😡", err)
	}
	fmt.Println()

}
