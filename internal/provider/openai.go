package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenAPIProvider struct {
	APIKey string // Provider'a auth işlemi yapması için APIKey'i saklar.
	URL    string // API Requestlerinin gönderileceği URL'i tutar.
}

func NewOpenAPIProvider(apiKey string) *OpenAPIProvider {
	return &OpenAPIProvider{
		APIKey: apiKey,
		URL:    "https://api.openai.com/v1/chat/completions",
	}
}

// (p *OpenAPIProvider) receiver kısmıdır ve anlamı OpenAPIProvider struct'ına ait bir metot olduğunu belirtir. GetName() metodu, bu sağlayıcının adını döndürür.
func (p *OpenAPIProvider) GetName() string {
	return "openai"
}

/*
Chat() çağrılır.
API'ye istek gönderilir.
Cevap alınır.
Başarılıysa *ChatResponse döner.
Hata varsa error döner.
*/
// (p *OpenAPIProvider) receiver kısmıdır ve OpenAPIProvider struct'ına ait bir metot olduğunu belirtir. Chat() metodu, sohbet isteklerini OpenAI API'sine gönderir ve yanıtları alır. ChatRequest ile bir sohbet isteği gönder; işlemi Context ile yönet; başarılı olursa ChatResponse, hata olursa error döndür.
func (p *OpenAPIProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {

	// OpenAI'ın beklediği payload yapısını oluşturur.. (Alanlar aynı ise direkt ChatRequest'i kullanabiliriz)
	// OpenAI standartlarını kullanacağımız için direkt req'i JSON yapısına çeviririz.
	body, err := json.Marshal(req) // ! ChatRequest yapısı JSON formatına dönüştürülür
	if err != nil {
		return nil, fmt.Errorf("openai marshal error: %v", err)
	}

	// ? HTTP isteği gönderilir (örneğin, net/http veya başka bir HTTP istemcisi kullanarak)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.URL, bytes.NewBuffer(body)) // HTTP isteği oluşturulur
	if err != nil {
		return nil, fmt.Errorf("openai request creation error: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")  // İçerik türü JSON olarak ayarlanır
	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey) // Yetkilendirme başlığı eklenir

	// Request gönderilir
	client := &http.Client{}        // HTTP istemcisi oluşturulur
	resp, err := client.Do(httpReq) // HTTP isteği gönderilir
	if err != nil {
		return nil, fmt.Errorf("OpenAI Request Error: %v", err)
	}
	defer resp.Body.Close() // Yanıtın kapanması sağlanır

	if resp.StatusCode != http.StatusOK {
		// Hata durumunda yanıt okunur ve hata mesajı döndürülür
		return nil, fmt.Errorf("OpenAI Response Error, Status: %d", resp.StatusCode)
	}

	// Response Parse edilir
	var openAIResponse struct {
		ID      string `json:"id"`
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
		Usage Usage `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
		return nil, fmt.Errorf("OpenAI Response Decode Error: %v", err)
	}

	// Kendi formatımıza dönüştürülür
	return &ChatResponse{
		ID:       openAIResponse.ID,                         // OpenAI tarafından döndürülen ID kullanılır
		Provider: "openai",                                  // Sağlayıcı adı "openai" olarak ayarlanır
		Model:    req.Model,                                 // Kullanılan model, istekten alınır
		Content:  openAIResponse.Choices[0].Message.Content, // İlk mesajın içeriği alınır
		Usage:    openAIResponse.Usage,                      // Token kullanımı OpenAI yanıtından alınır
	}, nil
}
