package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vendasta/langchaingo/llms"
	"github.com/vendasta/langchaingo/llms/openai"
)

func main() {
	ctx := context.Background()

	// Initialize OpenAI LLM with GPT Image 1.5
	llm, err := openai.New(
		openai.WithModel("gpt-image-1.5"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Example 1: Simple image generation with defaults
	fmt.Println("Example 1: Simple generation with defaults")
	resp1, err := llm.GenerateImage(ctx,
		"A photorealistic golden retriever playing in a sunny park",
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated image URL: %s\n", resp1.Images[0].URL)
	if resp1.Images[0].RevisedPrompt != "" {
		fmt.Printf("Revised prompt: %s\n", resp1.Images[0].RevisedPrompt)
	}
	fmt.Println()

	// Example 2: High-quality image with custom size
	fmt.Println("Example 2: High-quality with custom size")
	resp2, err := llm.GenerateImage(ctx,
		"A futuristic cyberpunk city at sunset with neon lights",
		llms.WithImageQuality("high"),
		llms.WithImageSize("1024x1536"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated image URL: %s\n", resp2.Images[0].URL)
	fmt.Println()

	// Example 3: Multiple images
	fmt.Println("Example 3: Generate 3 variations")
	resp3, err := llm.GenerateImage(ctx,
		"Abstract geometric patterns in blue and gold",
		llms.WithImageCount(3),
		llms.WithImageQuality("medium"),
	)
	if err != nil {
		log.Fatal(err)
	}
	for i, img := range resp3.Images {
		fmt.Printf("Image %d URL: %s\n", i+1, img.URL)
	}
	fmt.Println()

	// Example 4: Base64 response format
	fmt.Println("Example 4: Base64 format")
	resp4, err := llm.GenerateImage(ctx,
		"A cute cartoon robot",
		llms.WithImageResponseFormat("b64_json"),
		llms.WithImageQuality("low"), // Faster, cheaper
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Base64 data length: %d bytes\n", len(resp4.Images[0].B64JSON))
}
