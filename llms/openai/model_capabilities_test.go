package openai

import (
	"context"
	"testing"
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

	// Test that the gpt-image pattern matches correctly
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

	// Test that non-gpt-image models don't match
	nonGPTImageModels := []string{
		"gpt-4",
		"gpt-3.5-turbo",
		"dall-e-3",
	}

	for _, model := range nonGPTImageModels {
		t.Run("not_"+model, func(t *testing.T) {
			caps := getModelCapabilities(model)
			// These should match other patterns, not gpt-image
			// Just verify they get some capability (not testing specific values here)
			_ = caps.SupportsSystem // Just access to ensure no panic
		})
	}
}

func TestGenerateImageEmptyPrompt(t *testing.T) {
	t.Parallel()

	llm, err := New()
	if err != nil {
		// Skip test if API key is not available
		t.Skip("OpenAI API key not available")
	}

	ctx := context.Background()
	_, err = llm.GenerateImage(ctx, "")
	if err == nil {
		t.Error("Expected error for empty prompt, got nil")
	}
}
