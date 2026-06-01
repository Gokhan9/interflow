package middleware

import (
	"interflow/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(q *repository.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key") //Her request’te X-API-Key header’ını okuyor
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key is required"}) //API Key yoksa request durur, 401 Unauthorized döner
			return
		}

		user, err := q.GetUserByAPIKey(c.Request.Context(), apiKey) // API Key varmı diye kontrol et.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"}) // API Key bulunamazsa Invalid hatası ve request kesilir.
			return
		}

		// ! Kullanıcı bilgisi şuan için context içine saklanıyor.
		c.Set("user", user) //Eğer user var ise, kullanıcı bilgisi context içine saklanır. Böylece sonraki handler’larda bu bilgiye erişilebilir.
		c.Next()            // Requestlerin devamlılığını sağlar.
	}
}
