package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"interflow/internal/cache"
	"interflow/internal/config"
	"interflow/internal/database"
	"interflow/internal/handler"
	"interflow/internal/middleware"
	"interflow/internal/provider"
	"interflow/internal/repository"
	"interflow/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	//! Config Yükleme
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error. Config is not Loaded: %v", err)
	}

	//! Database Bağlantısı (Connection dosyasına taşıdım.)
	err = repository.InitDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	log.Println("Database connection successful - ✅")
	queries := database.New(repository.DB) // Repository katmanını başlatıyoruz. Bu, veritabanı işlemlerini gerçekleştirmek için kullanılacak sorguları içerir. Repository.New fonksiyonu, "pgxpool.Pool" türünde bir veritabanı bağlantısı alır ve bu bağlantıyı kullanarak sorguları hazırlar. Bu sayede, uygulamanın diğer bölümlerinde veritabanı işlemleri için bu sorguları kullanabiliriz.

	// ? 1- Analytics Service'i 5 WORKER ile başlatıyoruz..
	analyticsService := service.NewService(queries, 5)

	pManager := provider.NewManager() // Provider Manager ve OpenAI Provider'ı
	pManager.RegisterProvider(provider.NewOpenAPIProvider(cfg.OpenAIKey))

	chatHandler := handler.NewChatHandler(pManager, analyticsService) //"pManager, tüm providerları yöneten yapı(OpenAI, Claude)", "analyticsService, loglama ve metrik toplama gibi işler."

	//! Redis Bağlantısı
	err = cache.InitRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}
	log.Println("Redis connection successful - ✅")

	gin.SetMode(gin.ReleaseMode) // Production ortamında gereksiz logları kapatmak için

	//! Gin ile Router Starting
	router := gin.Default()

	// ! TEST Endpoint (Public)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	//! MIDDLEWARELAR
	router.Use(middleware.AuthMiddleware(queries))
	router.Use(middleware.RateLimitMiddleware())
	router.POST("/v1/chat", chatHandler.HandleChat)

	// ? 2- HTTP-SERVER YAPILANDIRMASI
	srv := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// ? 3- Sunucuyu bir goroutine içerisinde başlatmak için.
	go func() {
		log.Printf("Gateway is running on : %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	// ? 4- KAPATMA SINYALLERI (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Sinyal gelene kadar bekler....
	log.Println("Shutdown Server..")

	// ? 5- Graceful Shutdown & İşletim sistemine kapanış için 5-10 saniye süre
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ? 6- Yeni HTTP Requestleri Durdur..
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: ", err)
	}

	// ? 7- AnalyticsService'i kapat (Queue'de ki verilerin db'ye yazılmalarını bekler.)
	analyticsService.Shutdown()

	log.Println("server existing")

}
