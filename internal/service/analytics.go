package service

import (
	"context"
	"interflow/internal/analytics"
	"interflow/internal/database"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

type AnalyticsService struct {
	repo        *database.Queries         // AnalyticsService, veritabanı sorgularını gerçekleştirmek için database.Queries türünde bir repo alanına sahiptir.
	eventChan   chan analytics.UsageEvent // API eventlarını asenkron olarak işlemek için bir kanal (channel) kullanır.
	WorkerCount int                       // Olayları işlemek için kaç tane işçi (worker) kullanılacağını belirten bir alan.
	wg          sync.WaitGroup            // workerların bitmesini beklemek için..
}

func NewService(repo *database.Queries, workerCount int) *AnalyticsService {
	s := &AnalyticsService{
		repo:        repo,                                   // Veritabanı sorguları için repo atanır
		eventChan:   make(chan analytics.UsageEvent, 10000), // Olayları asenkron olarak işlemek için bir kanal oluşturulur (10000 kapasite sınırı)
		WorkerCount: workerCount,                            // İşçi sayısı atanır
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
	case s.eventChan <- event: // !goroutine, kanala veriyi gönderiyor. event geliyor, döngü çalışıyor.
		// Queue'ya başarıyla ekler.
	default:
		// Queue doluysa log düşebilir veya veri drop edilir. (HIZ'I Önceliklendiriyoruz.)
		log.Println("Analytics queue full, drop event.")
	}
}

func (s *AnalyticsService) worker() {

	defer s.wg.Done() // worker bitince sayacı azalt.

	// log'lar arka plana kaydedilir.
	for event := range s.eventChan { // channel dinlenir.
		err := s.repo.CreateUsageLog(context.Background(), database.CreateUsageLogParams{
			ApiKeyID:         int32(event.APIKeyID),
			Provider:         event.Provider,
			Model:            event.Model,
			PromptTokens:     pgtype.Int4{Int32: int32(event.PromptTokens), Valid: true},
			CompletionTokens: pgtype.Int4{Int32: int32(event.CompletionTokens), Valid: true},
			TotalTokens:      pgtype.Int4{Int32: int32(event.TotalTokens), Valid: true},
			LatencyMs:        pgtype.Int4{Int32: int32(event.LatencyMs), Valid: true},
		})

		if err != nil {
			log.Printf("failed to save usage log: %v", err)
		}
	}
}

func (s *AnalyticsService) Shutdown() {
	log.Println("Shutting down Analytics service")
	close(s.eventChan) // Kanala yeni veri girişini engeller, döngüleri bitirir.
	s.wg.Wait()        // Workerların "event"ları bitirmelerini bekler.
	log.Println("analytics service gracefully stopped.")
}
