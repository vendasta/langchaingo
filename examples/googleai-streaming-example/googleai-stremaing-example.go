// Set the GOOGLE_API_KEY env var to your API key taken from ai.google.dev
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/vendasta/langchaingo/llms"
	"github.com/vendasta/langchaingo/llms/googleaiv2"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("GOOGLE_API_KEY")
	// See https://ai.google.dev/gemini-api/docs/models/gemini for possible models
	llm, err := googleaiv2.New(ctx, googleaiv2.WithAPIKey(apiKey), googleaiv2.WithDefaultModel("gemini-3-pro-preview"))
	if err != nil {
		log.Fatal(err)
	}

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a company branding design wizard."),
		llms.TextParts(llms.ChatMessageTypeHuman, "What would be a good company name for a company that produces Go-backed LLM tools?"),
	}

	fmt.Println("--- Reasoning ---")
	var inContent bool
	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingReasoningFunc(func(ctx context.Context, reasoningChunk, chunk []byte) error {
		if len(reasoningChunk) > 0 {
			if inContent {
				fmt.Print("\n\n--- Additional Reasoning ---\n")
				inContent = false
			}
			fmt.Print(string(reasoningChunk))
		}
		if len(chunk) > 0 {
			if !inContent {
				fmt.Print("\n\n--- Content ---\n")
				inContent = true
			}
			fmt.Print(string(chunk))
		}
		return nil
	}), llms.WithThinkingMode(llms.ThinkingModeHigh))
	if err != nil {
		log.Fatal(err)
	}
	_ = completion
	fmt.Println()
}
