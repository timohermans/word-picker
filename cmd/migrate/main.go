package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Print(err)
	}
}

func main() {
	log := slog.Default()

	connectionStringEnv := strings.TrimSpace(os.Getenv("DATABASE_URL"))

	if strings.Contains(connectionStringEnv, "?") {
		panic("Querystring params are not supported here, as we manually add it")
	}

	connectionString := fmt.Sprintf("%s?sslmode=disable&search_path=app", connectionStringEnv)
	m, err := migrate.New("file://internal/migrations", connectionString)

	if err != nil {
		log.Error("failed to load migration config", "error", err.Error())
		panic(err)
	}

	log.Info("=== Migrating... ===")
	err = m.Up()
	log.Info("=== Migrations done! ===")

	if err != nil && !strings.Contains(err.Error(), "no change") {
		log.Error("failed to apply migration config", "error", err.Error())
		panic(err)
	}
}
