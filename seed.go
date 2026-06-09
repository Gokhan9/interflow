package main

import (
	"context"
	"fmt"
	"interflow/internal/config"
	"interflow/internal/database"
	"interflow/internal/repository"
	"log"
)

func main1() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = repository.InitDB(cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	queries := database.New(repository.DB)

	// Check users
	user, err := queries.GetUserByAPIKey(context.Background(), "test-api-key")
	if err != nil {
		fmt.Printf("Error or not found: %v\n", err)

		// Let's try to create one
		u, err := queries.CreateUser(context.Background(), "testuser_new")
		if err != nil {
			fmt.Printf("Failed to create user: %v\n", err)
		} else {
			fmt.Printf("Created user ID: %d\n", u.ID)
			ak, err := queries.CreateAPIKey(context.Background(), database.CreateAPIKeyParams{
				UserID:   u.ID,
				KeyValue: "test-api-key",
			})
			if err != nil {
				fmt.Printf("Failed to create API key: %v\n", err)
			} else {
				fmt.Printf("Created API key: %s\n", ak.KeyValue)
			}
		}
	} else {
		fmt.Printf("Found user: %s with API Key ID: %d\n", user.Username, user.ApiKeyID)
	}
}
