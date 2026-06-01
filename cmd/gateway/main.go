package main

import (
	"log"

	"interflow/internal/cache"
	"interflow/internal/config"
	"interflow/internal/middleware"
	"interflow/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {

	// ! Config Yükleme
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error. Config is not Loaded: %v", err)
	}

	// ! Database Bağlantısı (Connection dosyasına taşıdım.)
	err = repository.InitDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	log.Println("Database connection successful - ✅")
	queries := repository.New(repository.DB) // Repository katmanını başlatıyoruz. Bu, veritabanı işlemlerini gerçekleştirmek için kullanılacak sorguları içerir. Repository.New fonksiyonu, pgxpool.Pool türünde bir veritabanı bağlantısı alır ve bu bağlantıyı kullanarak sorguları hazırlar. Bu sayede, uygulamanın diğer bölümlerinde veritabanı işlemleri için bu sorguları kullanabiliriz.

	// ! Redis Bağlantısı
	err = cache.InitRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}
	log.Println("Redis connection successful - ✅")

	gin.SetMode(gin.ReleaseMode) // Production ortamında gereksiz logları kapatmak için

	// ! Gin ile Router Starting
	router := gin.Default() // Logger ve Recovery middleware'leri otomatik olarak ekler, tekrar eklemeye gerek kalmaz.

	/* ! Eğer özel middleware'leri MANUEL eklemek isterseniz, aşağıdaki gibi yapabilirsiniz. Ancak gin.Default() zaten bu middleware'leri içerdiği için, tekrar eklemeye gerek yok.
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	*/

	// ? MIDDLEWARELAR
	router.Use(middleware.AuthMiddleware(queries)) // Authentication middleware'ı tüm route'lara uygular. Her request'te API Key kontrolü yapar ve geçerli değilse 401 Unauthorized döner.
	router.Use(middleware.RateLimitMiddleware())   // Rate Limiting middleware'ı tüm route'lara uygular. Her kullanıcı için belirli bir süre içinde kaç istek attığını takip eder ve limit aşılırsa 429 Too Many Requests döner.

	// ! TEST Endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	log.Printf("Gateway is running on this port: %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server Başlatılamadı : %v", err)
	}
}
