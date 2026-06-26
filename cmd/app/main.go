package main

/*
packages used:
- [gotdotenv](https://github.com/joho/godotenv)
*/

import (
	"context"
	"hert/gotest/internal/db"
	"hert/gotest/internal/handlers"
	"hert/gotest/internal/html"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	addressDefault = ":44"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file: " + err.Error())
	}
}

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, strings.TrimSpace(os.Getenv("DATABASE_URL")))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	queries := db.New(pool)

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Renderer = &html.GomponentRendererRender{}

	handlers.HandleWordListOverview(e, queries)
	handlers.HandleWordListAdd(e, queries)
	handlers.HandleWordListPick(e, queries)
	handlers.RegisterHealthCheck(e, queries)

	if err := e.Start(addressDefault); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
