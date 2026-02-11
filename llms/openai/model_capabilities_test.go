package openai

import (
	"context"
	"strings"
	"testing"

	"github.com/vendasta/langchaingo/llms"
)

func TestModelCapabilities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		model              string
		wantSupportsSystem bool
	}{
		{
			name:               "GPT Image 1.5 supports system messages",
			model:              "gpt-image-1.5",
			wantSupportsSystem: true,
		},
		{
			name:               "GPT Image 1 supports system messages",
			model:              "gpt-image-1",
			wantSupportsSystem: true,
		},
		{
			name:               "GPT Image case insensitive",
			model:              "GPT-IMAGE-1.5",
			wantSupportsSystem: true,
		},
		{
			name:               "GPT-4 supports system messages",
			model:              "gpt-4",
			wantSupportsSystem: true,
		},
		{
			name:               "GPT-3.5 supports system messages",
			model:              "gpt-3.5-turbo",
			wantSupportsSystem: true,
		},
		{
			name:               "O1 does not support system messages",
			model:              "o1-preview",
			wantSupportsSystem: false,
		},
		{
			name:               "O3 does not support system messages",
			model:              "o3-mini",
			wantSupportsSystem: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			caps := getModelCapabilities(tt.model)
			if caps.SupportsSystem != tt.wantSupportsSystem {
				t.Errorf("getModelCapabilities(%q).SupportsSystem = %v, want %v",
					tt.model, caps.SupportsSystem, tt.wantSupportsSystem)
			}
		})
	}
}

func TestGPTImageModelPattern(t *testing.T) {
	t.Parallel()

	gptImageModels := []string{
		"gpt-image-1",
		"gpt-image-1.5",
		"GPT-IMAGE-2",
		"gpt-image-3.0",
	}

	for _, model := range gptImageModels {
		t.Run(model, func(t *testing.T) {
			caps := getModelCapabilities(model)
			if !caps.SupportsSystem {
				t.Errorf("Model %q should support system messages", model)
			}
			if caps.SupportsThinking {
				t.Errorf("Model %q should not support thinking/reasoning", model)
			}
		})
	}

	nonGPTImageModels := []string{
		"gpt-4",
		"gpt-3.5-turbo",
		"gpt-4o",
	}

	for _, model := range nonGPTImageModels {
		t.Run("not_"+model, func(t *testing.T) {
			caps := getModelCapabilities(model)
			_ = caps.SupportsSystem
		})
	}
}

func TestGenerateImageEmptyPrompt(t *testing.T) {
	t.Parallel()

	llm := &LLM{
		model: "gpt-image-1.5",
	}

	ctx := context.Background()
	_, err := llm.GenerateImage(ctx, "")

	if err == nil {
		t.Fatal("Expected error for empty prompt, got nil")
	}

	expectedMsg := "prompt cannot be empty"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error containing %q, got %q", expectedMsg, err.Error())
	}
}

func TestGenerateContentWithImageModel(t *testing.T) {
	t.Parallel()

	llm := &LLM{
		model: "gpt-image-1.5",
	}

	ctx := context.Background()
	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, "Hello"),
	}

	_, err := llm.GenerateContent(ctx, messages)

	if err == nil {
		t.Fatal("Expected error when using image model for text generation, got nil")
	}

	expectedMsg := "is for image generation only"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error containing %q, got %q", expectedMsg, err.Error())
	}
}
