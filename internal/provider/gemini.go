package provider

import (
	"context"
	"fmt"
)

type GeminiProvider struct {
	APIKey string
}

func NewGeminiProvider(apikey string) *GeminiProvider {
	return &GeminiProvider{
		APIKey: apikey,
	}
}

func (p *GeminiProvider) GetName() string {
	return "gemini"
}

func (p *GeminiProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {

	//1.Gemini'nin beklediği URL formatı
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		req.Model, p.APIKey)

	//2.Mesajları, Gemini'nin "content" yapısına maplemek
	var geminiMessages []geminiContent
	for _, msg := range req.Messages {
		role := msg.Role
		if role == "assistant" {
			role = "model" // gemini,   model kullanır.
		}
		geminiMessages = append(geminiMessages, geminiContent{
			Role:  role,
			Parts: []geminiPart{{Text: msg.Content}},
		})
	}

	//3.Request objesi oluşturmak

	geminiReq := geminiRequest{
		Contents: geminiMessages,
		GenerationConfig: geminiConfig{
			Temperature:     req.Temperature,
			MaxOutputTokens: req.MaxTokens,
		},
	}
}
