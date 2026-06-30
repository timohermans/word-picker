package main

/*
packages used:
- [gotdotenv](https://github.com/joho/godotenv)
*/

import (
	"context"
	"fmt"
	"hert/gotest/internal/db"
	"hert/gotest/internal/handlers"
	"hert/gotest/internal/html"
	"hert/gotest/internal/server"
	"log"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

const (
	addressDefault = "localhost:4444"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file: " + err.Error())
	}
}

func main() {
	var address string
	if address = os.Getenv("ADDRESS"); address == "" {
		address = addressDefault
	}

	// database setup
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, strings.TrimSpace(os.Getenv("DATABASE_URL")))
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	queries := db.New(pool)

	// auth setup
	useHttpsEnv := os.Getenv("USE_HTTPS")
	useHttps := slices.Contains([]string{"true", "1"}, strings.ToLower(useHttpsEnv))

	store := sessions.NewFilesystemStore(os.TempDir(), []byte(os.Getenv("SECURE_KEY")))
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = useHttps
	store.MaxLength(math.MaxInt64)

	gothic.Store = store
	protocol := "http"
	if useHttps {
		protocol = "https"
	}

	openidConnect, _ := openidConnect.New(
		os.Getenv("OIDC_CLIENT_ID"),
		os.Getenv("OIDC_CLIENT_SECRET"),
		fmt.Sprintf("%s://%s%s", protocol, address, server.UrlAuthCallback),
		os.Getenv("OIDC_DISCOVERY_URL"))

	if openidConnect != nil {
		goth.UseProviders(openidConnect)
	}

	// app setup
	e := echo.New()
	e.Use(server.NewSessionMiddleware(store))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(server.NewAuthenticationMiddleware())

	e.Renderer = &html.GomponentRendererRender{}

	server.RegisterAuthentication(e, handlers.UrlWordListOverview)
	handlers.HandleWordListOverview(e, queries)
	handlers.HandleWordListAdd(e, queries)
	handlers.HandleWordListPick(e, queries)
	handlers.RegisterHealthCheck(e, queries)

	if err := e.Start(":44"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
