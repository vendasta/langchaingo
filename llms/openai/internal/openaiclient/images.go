package openaiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	defaultImageModel = "gpt-image-1.5"
)

// imageGenerationPayload is the request payload for image generation.
type imageGenerationPayload struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model,omitempty"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	Style          string `json:"style,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

// imageGenerationResponse is the API response structure.
type imageGenerationResponse struct {
	Created int64                `json:"created"`
	Data    []generatedImageData `json:"data"`
}

// generatedImageData represents a single generated image.
type generatedImageData struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// ImageGenerationRequest is the public request structure.
type ImageGenerationRequest struct {
	Prompt         string
	Model          string
	N              int
	Size           string
	Quality        string
	Style          string
	ResponseFormat string
	User           string
}

// ImageGenerationResponse is the public response structure.
type ImageGenerationResponse struct {
	Created int64
	Images  []GeneratedImage
}

// GeneratedImage represents a single generated image.
type GeneratedImage struct {
	URL           string
	B64JSON       string
	RevisedPrompt string
}

// CreateImageGeneration calls the /v1/images/generations endpoint.
func (c *Client) CreateImageGeneration(ctx context.Context, r *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	payload := &imageGenerationPayload{
		Prompt:         r.Prompt,
		Model:          r.Model,
		N:              r.N,
		Size:           r.Size,
		Quality:        r.Quality,
		Style:          r.Style,
		ResponseFormat: r.ResponseFormat,
		User:           r.User,
	}

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

	return c.createImageGeneration(ctx, payload)
}

func (c *Client) createImageGeneration(ctx context.Context, payload *imageGenerationPayload) (*ImageGenerationResponse, error) {
	url := c.buildURL("/images/generations", payload.Model)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, sanitizeHTTPError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(resp)
	}

	var imageResp imageGenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&imageResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(imageResp.Data) == 0 {
		return nil, ErrEmptyResponse
	}

	images := make([]GeneratedImage, len(imageResp.Data))
	for i, img := range imageResp.Data {
		images[i] = GeneratedImage{
			URL:           img.URL,
			B64JSON:       img.B64JSON,
			RevisedPrompt: img.RevisedPrompt,
		}
	}

	return &ImageGenerationResponse{
		Created: imageResp.Created,
		Images:  images,
	}, nil
}
