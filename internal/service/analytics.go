package service

import (
	"interflow/internal/analytics"
	"interflow/internal/database"
	"log"
)

type AnalyticsService struct {
	repo        *database.Queries         // AnalyticsService, veritabanı sorgularını gerçekleştirmek için database.Queries türünde bir repo alanına sahiptir.
	eventChan   chan analytics.UsageEvent // API eventlarını asenkron olarak işlemek için bir kanal (channel) kullanır.
	WorkerCount int                       // Olayları işlemek için kaç tane işçi (worker) kullanılacağını belirten bir alan.
}

func NewService(repo *database.Queries, workerCount int) *AnalyticsService {
	s := &AnalyticsService{
		repo:        repo,                                  // Veritabanı sorguları için repo atanır
		eventChan:   make(chan analytics.UsageEvent, 1000), // Olayları asenkron olarak işlemek için bir kanal oluşturulur (1000 kapasite sınırı)
		WorkerCount: workerCount,                           // İşçi sayısı atanır
	}

	// ? Belirtilen sayıda işçi (worker) başlatılır. Her işçi, eventChan kanalından gelen UsageEvent'leri dinler ve işleme alır.
	for i := 0; i < workerCount; i++ {
		go s.worker()
	}

	return s

}

// "RECORD" bir event'i queue'ya atar.(Bloklamaz)
func (s *AnalyticsService) Record(event analytics.UsageEvent) {
	select {
	case s.eventChan <- event:
		// Queue'ya başarıyla ekler.
	default:
		// Queue doluysa log düşebilir veya veri drop edilir. (HIZ'I Önceliklendiriyoruz.)
		log.Println("Analytics queue full, drop event.")
	}
}

func (s *AnalyticsService) worker() {
	for event := range s.eventChan {
		log.Printf(
			"Provider=%s Model=%s Tokens=%d",
			event.Provider,
			event.Model,
			event.TotalTokens,
		)
		// Burada DB'ye kayıt işlemi yapılacak (s.repo.CreateUsageLog(...))
		// Bir hata olursa burada yakalanır, ana isteği etkilemez.
	}
}
