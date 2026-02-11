package openaiclient

import (
	"encoding/json"
	"testing"

	"github.com/vendasta/langchaingo/llms"
)

func TestChatRequest_MarshalJSON(t *testing.T) {
	tests := []struct {
		name                    string
		request                 ChatRequest
		wantMaxTokens           bool
		wantMaxCompletionTokens bool
	}{
		{
			name: "only MaxCompletionTokens set",
			request: ChatRequest{
				Model:               "gpt-4",
				MaxCompletionTokens: 100,
			},
			wantMaxTokens:           false,
			wantMaxCompletionTokens: true,
		},
		{
			name: "only MaxTokens set",
			request: ChatRequest{
				Model:     "gpt-4",
				MaxTokens: 200,
			},
			wantMaxTokens:           true,
			wantMaxCompletionTokens: false,
		},
		{
			name: "both set - only MaxCompletionTokens sent",
			request: ChatRequest{
				Model:               "gpt-4",
				MaxTokens:           300,
				MaxCompletionTokens: 400,
			},
			wantMaxTokens:           false,
			wantMaxCompletionTokens: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			hasMaxTokens := result["max_tokens"] != nil
			hasMaxCompletionTokens := result["max_completion_tokens"] != nil

			if hasMaxTokens != tt.wantMaxTokens {
				t.Errorf("max_tokens presence: got %v, want %v", hasMaxTokens, tt.wantMaxTokens)
			}
			if hasMaxCompletionTokens != tt.wantMaxCompletionTokens {
				t.Errorf("max_completion_tokens presence: got %v, want %v", hasMaxCompletionTokens, tt.wantMaxCompletionTokens)
			}

			// Never both
			if hasMaxTokens && hasMaxCompletionTokens {
				t.Error("Both max_tokens and max_completion_tokens are present - OpenAI API will reject!")
			}
		})
	}
}

func TestChatRequest_TemperatureMarshalJSON(t *testing.T) {
	tests := []struct {
		name            string
		request         ChatRequest
		wantTemperature bool
	}{
		{
			name: "regular model with temperature",
			request: ChatRequest{
				Model:       "gpt-4",
				Temperature: 0.7,
			},
			wantTemperature: true,
		},
		{
			name: "regular model with zero temperature",
			request: ChatRequest{
				Model:       "gpt-3.5-turbo",
				Temperature: 0.0,
			},
			wantTemperature: true,
		},
		{
			name: "gpt-5 model omits temperature",
			request: ChatRequest{
				Model:       "gpt-5-preview",
				Temperature: 0.7,
			},
			wantTemperature: false,
		},
		{
			name: "gpt-5 model omits zero temperature",
			request: ChatRequest{
				Model:       "gpt-5-mini",
				Temperature: 0.0,
			},
			wantTemperature: false,
		},
		{
			name: "o1 model omits temperature",
			request: ChatRequest{
				Model:       "o1-preview",
				Temperature: 0.5,
			},
			wantTemperature: false,
		},
		{
			name: "o1-mini model omits temperature",
			request: ChatRequest{
				Model:       "o1-mini",
				Temperature: 1.0,
			},
			wantTemperature: false,
		},
		{
			name: "o3 model omits temperature",
			request: ChatRequest{
				Model:       "o3-mini",
				Temperature: 0.8,
			},
			wantTemperature: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			hasTemperature := result["temperature"] != nil

			if hasTemperature != tt.wantTemperature {
				t.Errorf("temperature presence: got %v, want %v, JSON: %s", hasTemperature, tt.wantTemperature, string(data))
			}

			// If temperature should be present, verify the value
			if hasTemperature && tt.wantTemperature {
				temp, ok := result["temperature"].(float64)
				if !ok {
					t.Errorf("temperature is not a float64: %T", result["temperature"])
				} else if temp != tt.request.Temperature {
					t.Errorf("temperature value: got %v, want %v", temp, tt.request.Temperature)
				}
			}
		})
	}
}

func TestIsReasoningModel(t *testing.T) {
	tests := []struct {
		model    string
		expected bool
	}{
		// Regular models - should not be reasoning models
		{"gpt-4", false},
		{"gpt-3.5-turbo", false},
		{"gpt-4-turbo", false},
		{"gpt-4o", false},
		{"text-davinci-003", false},

		// GPT-5 models - should be reasoning models
		{"gpt-5", true},
		{"gpt-5-preview", true},
		{"gpt-5-mini", true},
		{"gpt-5-turbo", true},

		// o1 models - should be reasoning models
		{"o1-preview", true},
		{"o1-mini", true},
		{"o1-large", true},

		// o3 models - should be reasoning models
		{"o3", true}, // Base o3 model
		{"o3-mini", true},
		{"o3-preview", true},
		{"o3-large", true},

		// Edge cases
		{"", false},
		{"o10-preview", false}, // Doesn't start with "o1-"
		{"o30-mini", false},    // Doesn't start with "o3-"
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			result := isReasoningModel(tt.model)
			if result != tt.expected {
				t.Errorf("isReasoningModel(%q) = %v, want %v", tt.model, result, tt.expected)
			}
		})
	}
}

func TestChatMessage_BinaryContentTransformation(t *testing.T) {
	t.Parallel()

	testImageData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D,
		0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	msg := ChatMessage{
		Role: "user",
		MultiContent: []llms.ContentPart{
			llms.BinaryPart("image/png", testImageData),
			llms.TextPart("What color is this?"),
		},
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal ChatMessage: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if result["role"] != "user" {
		t.Errorf("Expected role 'user', got %v", result["role"])
	}

	content, ok := result["content"].([]interface{})
	if !ok {
		t.Fatalf("Expected content to be an array, got %T", result["content"])
	}

	if len(content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(content))
	}

	firstPart, ok := content[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected first part to be an object, got %T", content[0])
	}

	if firstPart["type"] != "image_url" {
		t.Errorf("Expected first part type to be 'image_url', got %v", firstPart["type"])
	}

	imageURL, ok := firstPart["image_url"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected image_url to be an object, got %T", firstPart["image_url"])
	}

	url, ok := imageURL["url"].(string)
	if !ok {
		t.Fatalf("Expected url to be a string, got %T", imageURL["url"])
	}

	expectedPrefix := "data:image/png;base64,"
	if len(url) < len(expectedPrefix) || url[:len(expectedPrefix)] != expectedPrefix {
		maxLen := 30
		if len(url) < maxLen {
			maxLen = len(url)
		}
		t.Errorf("Expected data URI starting with '%s', got %s", expectedPrefix, url[:maxLen])
	}

	secondPart, ok := content[1].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected second part to be an object, got %T", content[1])
	}

	if secondPart["type"] != "text" {
		t.Errorf("Expected second part type to be 'text', got %v", secondPart["type"])
	}

	if secondPart["text"] != "What color is this?" {
		t.Errorf("Expected text 'What color is this?', got %v", secondPart["text"])
	}
}
