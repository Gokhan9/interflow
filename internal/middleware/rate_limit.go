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
		userID := c.GetInt32("user_id")            // Kullanıcı ID'sini context'ten alıyoruz. Bu ID, AuthMiddleware tarafından context'e eklenmiş olmalıdır. Eğer kullanıcı ID'si bulunamazsa, varsayılan olarak 0 alınır.
		key := fmt.Sprintf("Ratelimit:%d", userID) // Redis'te her kullanıcı için benzersiz bir anahtar oluşturuyoruz. Bu anahtar, kullanıcının ID'sine dayanır ve "Ratelimit:" öneki ile başlar. Örneğin, kullanıcı ID'si 123 ise, anahtar "Ratelimit:123" olacaktır

		ctx := c.Request.Context()                              // Gin'in context'i, HTTP isteğiyle ilişkili bir context sağlar. Bu context, isteğin yaşam döngüsü boyunca geçerlidir ve isteğe özgü bilgileri taşıyabilir.
		count, err := cache.RedisClient.Incr(ctx, key).Result() // Redis'te anahtarın değerini "1" arttırmak için Incr kullanıyoruz. Eğer anahtar daha önce oluşturulmamışsa, Redis bu anahtarı "0" olarak başlatır ve ardından "1" ekler, böylece ilk kullanımda değer "1" olur. Incr komutu, anahtarın yeni değerini döndürür. Eğer bir hata oluşursa, bu hata err değişkenine atanır.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
			return
		}

		// Eğer user'ın attığı ilk request ise, Redis'te yer alan anahtarın geçerlilik süresini 1 dakika olarak ayarlıyurz. Bu, kullanıcının 1 dakika boyunca kaç istek attığını takip etmemizi sağlar. Eğer kullanıcı 1 dakika içinde 10'dan fazla istek atarsa, bu middleware tarafından engellenir.
		if count == 1 {
			cache.RedisClient.Expire(ctx, key, time.Minute)
		}

		// Request sayısı 10'u geçerse, user'a "429 Too Many Requests" hatası döner. Request durdurulur. Bu, kullanıcının belirli bir süre içinde çok fazla istek atmasını önlemek için kullanılan bir mekanizmadır. Bu sayede sunucunun aşırı yüklenmesi engellenir ve hizmet kalitesi korunur.
		if count > 10 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"ERROR": "Rate limit exceeded"})
			return
		}

		c.Next()
	}
}
