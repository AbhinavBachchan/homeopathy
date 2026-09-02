package main

import (
	"log"

	"homeopathy-platform/internal/config"
	"homeopathy-platform/internal/db"
	"homeopathy-platform/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	// .env is optional in production (env vars set by the platform instead)
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	cfg := config.Load()

	conn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(conn); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	r := router.New(cfg, conn)

	log.Printf("listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
