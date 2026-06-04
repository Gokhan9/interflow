package analytics

type UsageEvent struct {
	APIKeyID         int    // API anahtarının benzersiz kimliği
	Provider         string // AI sağlayıcısının adı (örneğin, "openai", "gemini")
	Model            string // Kullanılan AI modelinin adı (örneğin, "gpt-4", "gemini-1")
	PromptTokens     int    // API isteğinde kullanılan prompt token sayısı
	CompletionTokens int    // API isteğinde kullanılan completion token sayısı
	TotalTokens      int    // API isteğinde kullanılan toplam token sayısı (prompt + completion)
	LatencyMs        int64  // API isteğinin yanıt süresi (milisaniye cinsinden)
	StatusCode       int    // API isteğinin HTTP durum kodu

}
