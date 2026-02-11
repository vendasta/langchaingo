# OpenAI Image Generation Example

This example demonstrates how to use OpenAI's GPT Image 1.5 model for image generation through LangChain Go.

## Overview

LangChain Go supports image generation with GPT Image 1.5 - the latest image generation model with advanced instruction following and high-quality output.

## Features Demonstrated

1. Simple image generation with default parameters
2. Custom quality and size settings
3. Generating multiple image variations
4. Base64 response format for immediate storage

## Prerequisites

Set your OpenAI API key:

```bash
export OPENAI_API_KEY=your-api-key-here
```

## Running the Example

```bash
go run openai_image_generation_example.go
```

## Usage Examples

### Basic Generation

```go
llm, _ := openai.New(openai.WithModel("gpt-image-1.5"))
resp, err := llm.GenerateImage(ctx, "A serene mountain landscape")
```

### Custom Parameters

```go
resp, err := llm.GenerateImage(ctx,
    "A futuristic city",
    llms.WithImageQuality("high"),
    llms.WithImageSize("1024x1536"),
    llms.WithImageCount(3),
)
```

## Available Options

- `llms.WithImageModel(model)` - Specify model (default: "gpt-image-1.5")
- `llms.WithImageCount(n)` - Number of images (1-10)
- `llms.WithImageSize(size)` - Dimensions: "1024x1024", "1024x1536", "1536x1024"
- `llms.WithImageQuality(quality)` - Quality: "low", "medium", "high"
- `llms.WithImageStyle(style)` - Style: "vivid", "natural" (model-specific)
- `llms.WithImageResponseFormat(format)` - Format: "url", "b64_json"
- `llms.WithImageUser(user)` - User identifier for monitoring

## Default Values

When not specified:
- Model: `gpt-image-1.5`
- Count: `1`
- Size: `1024x1024`
- Quality: `medium`
- Format: `url`

## Pricing

According to OpenAI's pricing:

**Image Generation:**
- Low quality: $0.009-$0.013 per image
- Medium quality: $0.034-$0.05 per image
- High quality: $0.133-$0.20 per image

## Important Notes

- Generated image URLs expire after 1 hour
- Use `b64_json` format if you need to store images immediately
- The model may revise your prompt for safety or clarity (check `RevisedPrompt` field)
- For production use, consider implementing retry logic and error handling

## Learn More

- [OpenAI Image Generation API Documentation](https://platform.openai.com/docs/guides/image-generation)
- [LangChain Go Documentation](https://pkg.go.dev/github.com/vendasta/langchaingo)
