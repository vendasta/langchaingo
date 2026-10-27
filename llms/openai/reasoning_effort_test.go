package openai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vendasta/langchaingo/llms"
)

// TestReasoningEffortRequest verifies the reasoning_effort and temperature
// sent in the request body for each thinking mode and model family.
func TestReasoningEffortRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		model           string
		mode            llms.ThinkingMode
		temperature     float64
		wantEffort      string // empty means reasoning_effort is omitted
		wantTemperature *float64
	}{
		{"gpt-6 none sends none and temperature 0", "gpt-6-luna", llms.ThinkingModeNone, 0, "none", ptr(0.0)},
		{"gpt-6 low keeps temperature", "gpt-6-luna", llms.ThinkingModeLow, 1, "low", ptr(1.0)},
		{"gpt-6 auto sends no effort", "gpt-6-luna", llms.ThinkingModeAuto, 1, "", ptr(1.0)},
		{"gpt-5 low omits temperature", "gpt-5", llms.ThinkingModeLow, 0, "low", nil},
		{"gpt-5 none sends minimal", "gpt-5", llms.ThinkingModeNone, 0, "minimal", nil},
		{"gpt-5-mini none sends minimal", "gpt-5-mini", llms.ThinkingModeNone, 0, "minimal", nil},
		{"gpt-5.1 none sends none", "gpt-5.1", llms.ThinkingModeNone, 0, "none", nil},
		{"gpt-4o sends no effort", "gpt-4o", llms.ThinkingModeLow, 0, "", ptr(0.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var body map[string]any
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				raw, err := io.ReadAll(r.Body)
				if err == nil {
					err = json.Unmarshal(raw, &body)
				}
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"chatcmpl-1","object":"chat.completion","created":1,"model":"` + tt.model +
					`","choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`))
			}))
			defer server.Close()

			llm, err := New(
				WithToken("test-key"),
				WithBaseURL(server.URL),
				WithModel(tt.model),
			)
			require.NoError(t, err)

			_, err = llm.GenerateContent(context.Background(),
				[]llms.MessageContent{llms.TextParts(llms.ChatMessageTypeHuman, "hi")},
				llms.WithThinkingMode(tt.mode),
				llms.WithTemperature(tt.temperature),
			)
			require.NoError(t, err)

			if tt.wantEffort == "" {
				assert.NotContains(t, body, "reasoning_effort")
			} else {
				assert.Equal(t, tt.wantEffort, body["reasoning_effort"])
			}
			if tt.wantTemperature == nil {
				assert.NotContains(t, body, "temperature")
			} else {
				assert.Equal(t, *tt.wantTemperature, body["temperature"])
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }
