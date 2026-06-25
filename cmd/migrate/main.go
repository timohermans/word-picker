package main

import (
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
}

func main() {
	log := slog.Default()

	m, err := migrate.New(
		"file://internal/migrations",
		os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Error("failed to load migration config", "error", err.Error())
		panic(err)
	}

	log.Info("=== Migrating... ===")
	err = m.Up()
	log.Info("=== Migrations done! ===")

	if err != nil {
		log.Error("failed to apply migration config", "error", err.Error())
		panic(err)
	}
}
