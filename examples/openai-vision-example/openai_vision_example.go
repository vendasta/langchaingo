package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/vendasta/langchaingo/llms"
	"github.com/vendasta/langchaingo/llms/openai"
)

func main() {
	ctx := context.Background()

	// Create an OpenAI LLM with a vision-capable model
	llm, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OpenAI Vision Examples")
	fmt.Println("======================\n")

	// Example 1: Analyze an image from URL
	fmt.Println("Example 1: Analyze image from URL")
	// Use a reliable, publicly accessible image
	imageURL := "https://raw.githubusercontent.com/tmc/langchaingo/main/docs/static/img/parrot-icon.png"
	
	resp1, err := llm.GenerateContent(ctx, []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.ImageURLPart(imageURL),
				llms.TextPart("What do you see in this image? Describe the scene."),
			},
		},
	}, llms.WithMaxTokens(200))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s\n\n", resp1.Choices[0].Content)

	// Example 2: Analyze a base64-encoded image
	fmt.Println("Example 2: Analyze base64-encoded image")
	// For this example, we'll use a small test image
	// In production, you'd read from a file
	imageData := createTestImage()
	
	resp2, err := llm.GenerateContent(ctx, []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.BinaryPart("image/png", imageData),
				llms.TextPart("Describe this image briefly."),
			},
		},
	}, llms.WithMaxTokens(100))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s\n\n", resp2.Choices[0].Content)

	// Example 3: Multi-turn conversation with images
	fmt.Println("Example 3: Multi-turn conversation with image context")
	messages := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.ImageURLPart(imageURL),
				llms.TextPart("What's the dominant color in this image?"),
			},
		},
	}

	resp3, err := llm.GenerateContent(ctx, messages, llms.WithMaxTokens(50))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("First response: %s\n", resp3.Choices[0].Content)

	// Add AI response and follow-up question
	messages = append(messages,
		llms.MessageContent{
			Role:  llms.ChatMessageTypeAI,
			Parts: []llms.ContentPart{llms.TextPart(resp3.Choices[0].Content)},
		},
		llms.MessageContent{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextPart("What mood does this color create?")},
		},
	)

	resp4, err := llm.GenerateContent(ctx, messages, llms.WithMaxTokens(100))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Follow-up response: %s\n\n", resp4.Choices[0].Content)

	// Example 4: Comparing multiple images
	fmt.Println("Example 4: Analyzing multiple images in one request")
	resp5, err := llm.GenerateContent(ctx, []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextPart("Compare these two images:"),
				llms.ImageURLPart(imageURL),
				llms.BinaryPart("image/png", imageData),
				llms.TextPart("What are the main differences?"),
			},
		},
	}, llms.WithMaxTokens(200))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s\n", resp5.Choices[0].Content)
}

// createTestImage creates a small test PNG image (1x1 red pixel)
func createTestImage() []byte {
	// Minimal valid PNG file (1x1 red pixel)
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 dimensions
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D,
		0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, // IEND chunk
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}

// Helper to convert image file to base64 (for reference)
func imageToBase64(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
