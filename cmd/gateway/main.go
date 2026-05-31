package main

import (
	"log"

	"interflow/internal/cache"
	"interflow/internal/config"
	"interflow/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {

	// ! Config Yükleme
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error. Config is not Loaded: %v", err)
	}

	// ! Database Bağlantısı
	err = repository.InitDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	log.Println("Database connection successful - ✅")

	// ! Redis Bağlantısı
	err = cache.InitRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}
	log.Println("Redis connection successful - ✅")

	// ! Gin ile Router Starting
	r := gin.Default()

	// ! TEST Endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	log.Printf("Gateway is running on this port: %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server Başlatılamadı : %v", err)
	}
}
