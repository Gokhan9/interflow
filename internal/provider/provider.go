package provider

import "context"

/*
* Message, ChatRequest, ChatResponse ve Usage struct'ları, AI servisleriyle yapılan sohbet isteklerini ve yanıtlarını temsil eder.
 */
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

/*
ChatRequest, bir sohbet isteğini temsil eder.
Model, kullanılacak AI modelini belirtir.
Messages, sohbet geçmişini içerir.
Temperature, yanıtın rastgeleliğini kontrol eder (daha yüksek değerler daha yaratıcı yanıtlar üretir).
MaxTokens, yanıtın maksimum token sayısını belirler.
*/
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

//
/*
ChatResponse, bir sohbet yanıtını temsil eder.
ID, yanıtın benzersiz kimliğidir.
Provider, yanıtı üreten AI servisinin adıdır.
Model, kullanılan AI modelini belirtir.
Content, AI tarafından üretilen yanıtın içeriğidir.
?Usage, istekte kullanılan token sayısını ve toplam token sayısını içerir.
*/
type ChatResponse struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Content  string `json:"content"`
	Usage    Usage  `json:"usage"`
}

/*
Usage, bir sohbet isteğinde kullanılan token sayısını ve toplam token sayısını temsil eder.
*/
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Provider, farklı AI servislerinin (OpenAI, Gemini vs.) uyması gereken kontrattır.
type Provider interface {
	GetName() string
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}
