package middleware

import (
	"fmt"
	"interflow/internal/cache"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
			return
		}

		userID, ok := val.(int32)
		if !ok {
			// Eğer int32 değilse, belki int gelmiştir (bazı durumlarda otomatik dönüşebiliyor)
			if intVal, ok := val.(int); ok {
				userID = int32(intVal)
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type in context"})
				return
			}
		}

		key := fmt.Sprintf("Ratelimit:%d", userID)
		ctx := c.Request.Context()
		count, err := cache.RDB.Incr(ctx, key).Result()
		
		// DEBUG LOG
		fmt.Printf("[DEBUG] RateLimit - UserID: %d, Key: %s, Count: %d\n", userID, key, count)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
			return
		}

		// Eğer user'ın attığı ilk request ise, Redis'te yer alan anahtarın geçerlilik süresini 1 dakika olarak ayarlıyurz. Bu, kullanıcının 1 dakika boyunca kaç istek attığını takip etmemizi sağlar. Eğer kullanıcı 1 dakika içinde 10'dan fazla istek atarsa, bu middleware tarafından engellenir.
		if count == 1 {
			cache.RDB.Expire(ctx, key, time.Minute)
		}

		// Request sayısı 10'u geçerse, user'a "429 Too Many Requests" hatası döner. Request durdurulur. Bu, kullanıcının belirli bir süre içinde çok fazla istek atmasını önlemek için kullanılan bir mekanizmadır. Bu sayede sunucunun aşırı yüklenmesi engellenir ve hizmet kalitesi korunur.
		if count > 10 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"ERROR": "Rate limit exceeded"})
			return
		}

		c.Next()
	}
}
