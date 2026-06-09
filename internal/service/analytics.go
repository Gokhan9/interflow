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
	repo        *database.Queries
	eventChan   chan analytics.UsageEvent
	WorkerCount int
	wg          sync.WaitGroup
}

func NewService(repo *database.Queries, workerCount int) *AnalyticsService {
	s := &AnalyticsService{
		repo:        repo,
		eventChan:   make(chan analytics.UsageEvent, 10000),
		WorkerCount: workerCount,
	}

	for i := 0; i < workerCount; i++ {
		s.wg.Add(1)
		go s.worker()
	}

	return s
}

func (s *AnalyticsService) Record(event analytics.UsageEvent) {
	select {
	case s.eventChan <- event:
	default:
		log.Println("Analytics queue full, drop event.")
	}
}

func (s *AnalyticsService) worker() {
	defer s.wg.Done()

	for event := range s.eventChan {
		err := s.repo.CreateUsageLog(context.Background(), database.CreateUsageLogParams{
			ApiKeyID:         int32(event.APIKeyID),
			Provider:         event.Provider,
			Model:            event.Model,
			PromptTokens:     pgtype.Int4{Int32: int32(event.PromptTokens), Valid: true},
			CompletionTokens: pgtype.Int4{Int32: int32(event.CompletionTokens), Valid: true},
			TotalTokens:      pgtype.Int4{Int32: int32(event.TotalTokens), Valid: true},
			LatencyMs:        pgtype.Int4{Int32: int32(event.LatencyMs), Valid: true},
			StatusCode:       pgtype.Int4{Int32: int32(event.StatusCode), Valid: true},
		})

		if err != nil {
			log.Printf("failed to save usage log: %v", err)
		}
	}
}

func (s *AnalyticsService) Shutdown() {
	log.Println("Shutting down Analytics service")
	close(s.eventChan)
	s.wg.Wait()
	log.Println("analytics service gracefully stopped.")
}
