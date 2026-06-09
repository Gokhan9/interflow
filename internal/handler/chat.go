package handler

import (
	"interflow/internal/analytics"
	"interflow/internal/provider"
	"interflow/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	providerMgr *provider.ProviderManager
	analytics   *service.AnalyticsService
}

func NewChatHandler(pm *provider.ProviderManager, as *service.AnalyticsService) *ChatHandler {
	return &ChatHandler{
		providerMgr: pm,
		analytics:   as,
	}
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
	//1- CONTEXT'ten API-KEY-ID'yi al. (Middleware'de set ettik.)
	//apiKeyID, _ := c.Get("api_key_id")

	val, exists := c.Get("api_key_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "api_key_id not found in context"})
		return
	}

	apiKeyID, ok := val.(int32)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid api_key_id type"})
		return
	}

	//2- REQUEST parse
	var req provider.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	//3- Provider'ı bul(openai)
	providerName := c.GetHeader("X-Provider")
	if providerName == "" {
		providerName = "openai"
	}

	p, err := h.providerMgr.GetProvider(providerName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	//4- İsteği gönder ve latency ölç.
	start := time.Now()
	resp, err := p.Chat(c.Request.Context(), req)
	latency := time.Since(start).Milliseconds()

	//5- Analytics'e asenkron olarak ilet.
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	usageEvent := analytics.UsageEvent{
		APIKeyID: int(apiKeyID),
		Provider: providerName,
		Model:    req.Model,
		//PromptTokens:     resp.Usage.PromptTokens,
		//CompletionTokens: resp.Usage.CompletionTokens,
		//TotalTokens:      resp.Usage.TotalTokens,
		LatencyMs:  latency,
		StatusCode: statusCode,
	}

	if resp != nil {
		usageEvent.PromptTokens = resp.Usage.PromptTokens
		usageEvent.CompletionTokens = resp.Usage.CompletionTokens
		usageEvent.TotalTokens = resp.Usage.TotalTokens
	}

	h.analytics.Record(usageEvent)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)

}
