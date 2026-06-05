package middleware

import (
	"interflow/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(q *database.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key") //Her request’te X-API-Key header’ını okuyor
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key is required"}) //API Key yoksa request durur, 401 Unauthorized döner
			return
		}

		// GetUserByAPIKey sorgusu güncellendikten sonra bu "resp" içerisinde hem "user" hem de "ApiKeyID" var.
		resp, err := q.GetUserByAPIKey(c.Request.Context(), apiKey) // API Key varmı diye kontrol et.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"}) // API Key bulunamazsa Invalid hatası ve request kesilir.
			return
		}

		// ! Kullanıcı bilgisi şuan için context içine saklanıyor.
		c.Set("user_id", resp.ID)          // Kullanıcı bilgisi request boyunca erişilebilir hale gelir.
		c.Set("api_key_id", resp.ApiKeyID) // analytics için kullanılacak.
		c.Next()                           // Requestlerin devamlılığını sağlar.
	}
}
