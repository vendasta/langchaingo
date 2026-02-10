package openaiclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateImageGeneration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		request        *ImageGenerationRequest
		serverResponse imageGenerationResponse
		wantErr        bool
		checkDefaults  bool
	}{
		{
			name: "successful generation with all defaults",
			request: &ImageGenerationRequest{
				Prompt: "A test image",
			},
			serverResponse: imageGenerationResponse{
				Created: 1234567890,
				Data: []generatedImageData{
					{
						URL:           "https://example.com/image.png",
						RevisedPrompt: "A test image (revised)",
					},
				},
			},
			wantErr:       false,
			checkDefaults: true,
		},
		{
			name: "successful generation with custom parameters",
			request: &ImageGenerationRequest{
				Prompt:         "A custom image",
				Model:          "dall-e-3",
				N:              2,
				Size:           "1024x1536",
				Quality:        "high",
				Style:          "vivid",
				ResponseFormat: "b64_json",
				User:           "test-user",
			},
			serverResponse: imageGenerationResponse{
				Created: 1234567890,
				Data: []generatedImageData{
					{
						B64JSON:       "base64data1",
						RevisedPrompt: "Custom image 1",
					},
					{
						B64JSON:       "base64data2",
						RevisedPrompt: "Custom image 2",
					},
				},
			},
			wantErr:       false,
			checkDefaults: false,
		},
		{
			name: "multiple images",
			request: &ImageGenerationRequest{
				Prompt: "Generate variations",
				N:      3,
			},
			serverResponse: imageGenerationResponse{
				Created: 1234567890,
				Data: []generatedImageData{
					{URL: "https://example.com/image1.png"},
					{URL: "https://example.com/image2.png"},
					{URL: "https://example.com/image3.png"},
				},
			},
			wantErr:       false,
			checkDefaults: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify method
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				// Verify path
				if r.URL.Path != "/images/generations" {
					t.Errorf("Expected path /images/generations, got %s", r.URL.Path)
				}

				// Read and verify request body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read request body: %v", err)
				}

				var payload imageGenerationPayload
				if err := json.Unmarshal(body, &payload); err != nil {
					t.Fatalf("Failed to unmarshal request: %v", err)
				}

				// Verify prompt
				if payload.Prompt != tt.request.Prompt {
					t.Errorf("Expected prompt %q, got %q", tt.request.Prompt, payload.Prompt)
				}

				// Check defaults were applied
				if tt.checkDefaults {
					if payload.Model == "" || payload.Model != defaultImageModel {
						t.Errorf("Expected default model %q, got %q", defaultImageModel, payload.Model)
					}
					if payload.N == 0 || payload.N != 1 {
						t.Errorf("Expected default N=1, got %d", payload.N)
					}
					if payload.Size == "" || payload.Size != "1024x1024" {
						t.Errorf("Expected default size 1024x1024, got %q", payload.Size)
					}
					if payload.Quality == "" || payload.Quality != "medium" {
						t.Errorf("Expected default quality medium, got %q", payload.Quality)
					}
					if payload.ResponseFormat == "" || payload.ResponseFormat != "url" {
						t.Errorf("Expected default response format url, got %q", payload.ResponseFormat)
					}
				}

				// Send response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(tt.serverResponse); err != nil {
					t.Fatalf("Failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			// Create client
			client, err := New(
				"test-token",
				"",
				server.URL,
				"",
				APITypeOpenAI,
				"",
				http.DefaultClient,
				"",
				nil,
				"",
				false,
			)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Call CreateImageGeneration
			resp, err := client.CreateImageGeneration(context.Background(), tt.request)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateImageGeneration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify response
			if resp.Created != tt.serverResponse.Created {
				t.Errorf("Expected Created %d, got %d", tt.serverResponse.Created, resp.Created)
			}

			if len(resp.Images) != len(tt.serverResponse.Data) {
				t.Errorf("Expected %d images, got %d", len(tt.serverResponse.Data), len(resp.Images))
			}

			for i, img := range resp.Images {
				expected := tt.serverResponse.Data[i]
				if img.URL != expected.URL {
					t.Errorf("Image %d: Expected URL %q, got %q", i, expected.URL, img.URL)
				}
				if img.B64JSON != expected.B64JSON {
					t.Errorf("Image %d: Expected B64JSON %q, got %q", i, expected.B64JSON, img.B64JSON)
				}
				if img.RevisedPrompt != expected.RevisedPrompt {
					t.Errorf("Image %d: Expected RevisedPrompt %q, got %q", i, expected.RevisedPrompt, img.RevisedPrompt)
				}
			}
		})
	}
}

func TestImageGenerationDefaults(t *testing.T) {
	t.Parallel()

	// Test that defaults are properly applied
	req := &ImageGenerationRequest{
		Prompt: "Test",
	}

	payload := &imageGenerationPayload{
		Prompt: req.Prompt,
		Model:  req.Model,
		N:      req.N,
		Size:   req.Size,
	}

	// Apply defaults (mimicking CreateImageGeneration logic)
	if payload.Model == "" {
		payload.Model = defaultImageModel
	}
	if payload.N == 0 {
		payload.N = 1
	}
	if payload.Size == "" {
		payload.Size = "1024x1024"
	}
	if payload.Quality == "" {
		payload.Quality = "medium"
	}
	if payload.ResponseFormat == "" {
		payload.ResponseFormat = "url"
	}

	// Verify defaults
	if payload.Model != "gpt-image-1.5" {
		t.Errorf("Expected default model gpt-image-1.5, got %s", payload.Model)
	}
	if payload.N != 1 {
		t.Errorf("Expected default N=1, got %d", payload.N)
	}
	if payload.Size != "1024x1024" {
		t.Errorf("Expected default size 1024x1024, got %s", payload.Size)
	}
	if payload.Quality != "medium" {
		t.Errorf("Expected default quality medium, got %s", payload.Quality)
	}
	if payload.ResponseFormat != "url" {
		t.Errorf("Expected default response format url, got %s", payload.ResponseFormat)
	}
}

func TestCreateImageGenerationEmptyResponse(t *testing.T) {
	t.Parallel()

	// Test that empty response returns appropriate error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Return empty data array
		resp := imageGenerationResponse{
			Created: 1234567890,
			Data:    []generatedImageData{},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client, err := New(
		"test-token",
		"",
		server.URL,
		"",
		APITypeOpenAI,
		"",
		http.DefaultClient,
		"",
		nil,
		"",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &ImageGenerationRequest{
		Prompt: "Test",
	}

	_, err = client.CreateImageGeneration(context.Background(), req)
	if err == nil {
		t.Error("Expected error for empty response, got nil")
	}
	if err != nil && err != ErrEmptyResponse {
		t.Errorf("Expected ErrEmptyResponse, got: %v", err)
	}
}
