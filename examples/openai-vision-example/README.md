# OpenAI Vision Example

This example demonstrates how to use OpenAI's vision-capable models (like GPT-4o) to analyze images using the LangChain Go library.

## Vision Capabilities

OpenAI's vision models can:
- Analyze images from URLs
- Process base64-encoded images
- Handle multiple images in a single request
- Maintain image context across multi-turn conversations
- Answer questions about image content, objects, colors, text, and more

## Features Demonstrated

1. Image analysis from public URLs
2. Base64-encoded image processing
3. Multi-turn conversations with image context
4. Comparing multiple images in one request

## Usage

```bash
export OPENAI_API_KEY=your-api-key-here
go run openai_vision_example.go
```

## Vision-Capable Models

The following OpenAI models support vision:
- `gpt-4o` (recommended)
- `gpt-4o-mini`
- `gpt-4-turbo`
- `gpt-4-turbo-2024-04-09`
- `gpt-image-1.5` (primarily for image generation, but supports vision input)

## Code Examples

### Analyze Image from URL

```go
llm, _ := openai.New(openai.WithModel("gpt-4o"))

resp, err := llm.GenerateContent(ctx, []llms.MessageContent{
    {
        Role: llms.ChatMessageTypeHuman,
        Parts: []llms.ContentPart{
            llms.ImageURLPart("https://example.com/image.jpg"),
            llms.TextPart("What do you see in this image?"),
        },
    },
})
```

### Analyze Base64-Encoded Image

```go
imageData, _ := os.ReadFile("photo.jpg")

resp, err := llm.GenerateContent(ctx, []llms.MessageContent{
    {
        Role: llms.ChatMessageTypeHuman,
        Parts: []llms.ContentPart{
            llms.BinaryPart("image/jpeg", imageData),
            llms.TextPart("Describe this image."),
        },
    },
})
```

### Multiple Images

```go
resp, err := llm.GenerateContent(ctx, []llms.MessageContent{
    {
        Role: llms.ChatMessageTypeHuman,
        Parts: []llms.ContentPart{
            llms.TextPart("Compare these images:"),
            llms.ImageURLPart(url1),
            llms.ImageURLPart(url2),
            llms.TextPart("What are the differences?"),
        },
    },
})
```

## Supported Image Formats

- PNG (`image/png`)
- JPEG (`image/jpeg`)
- WEBP (`image/webp`)
- GIF (`image/gif`, non-animated)

## Image Size Limits

- Maximum file size: 20MB
- Recommended: Keep images under 4MB for faster processing
- Images are automatically resized if they exceed model limits

## Options

Control token usage and response length:
- `llms.WithMaxTokens(n)` - Limit response length
- `llms.WithTemperature(t)` - Control creativity (0.0-2.0)
- `llms.WithTopP(p)` - Nucleus sampling parameter

## Important Notes

- Image URLs must be publicly accessible or properly authenticated
- Base64 encoding increases token usage
- Vision API calls cost more tokens than text-only calls
- The model may refuse to process certain images (safety filters)
- For production use, implement proper error handling and retries

## Use Cases

- **Content Moderation**: Analyze user-uploaded images
- **E-commerce**: Extract product details from photos
- **Document Processing**: Read text from images (OCR alternative)
- **Accessibility**: Generate image descriptions for visually impaired users
- **Quality Control**: Inspect product images for defects
- **Medical Imaging**: Preliminary analysis (with appropriate disclaimers)

## Learn More

- [OpenAI Vision API Documentation](https://platform.openai.com/docs/guides/vision)
- [GPT-4 Vision Guide](https://platform.openai.com/docs/guides/images-vision)
- [LangChain Go Documentation](https://pkg.go.dev/github.com/vendasta/langchaingo)
