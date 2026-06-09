package main

import (
	"interflow/internal/config"
	"interflow/internal/database"
	"interflow/internal/repository"
)

func main() {
	cfg, _ := config.LoadConfig()
	repository.InitDB(cfg.DBURL)
	queries1 := database.New(repository.DB)

	// Check usage logs
	// Need a query for this. queries.sql didn't have a GetUsageLogs.
	// I'll just use raw sql if possible, but queries.sql.go is generated.
	// I'll check what's in internal/database/queries.sql.go
}
