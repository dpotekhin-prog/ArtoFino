// @title ArtoFino Backend API
// @version 1.0
// @description Backend service for ArtoFino platform.
// @host localhost:9000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer <token>"
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "ArtoFino/backend/docs" // <-- ВАЖНО! Без этого Swagger НЕ РАБОТАЕТ

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ArtoFino/backend/internal/auth"
	"ArtoFino/backend/internal/config"
	"ArtoFino/backend/internal/db"
	"ArtoFino/backend/internal/http/handlers"

	apphttp "ArtoFino/backend/internal/http"
	mongostore "ArtoFino/backend/internal/mongo"
)

func main() {
	// --- Load .env ---
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println(".env loaded successfully")
	}

	// --- Load config ---
	cfg := config.Load()

	// --- Init Keycloak ---
	kc, err := auth.NewKeycloakClient(cfg.Keycloak)
	if err != nil {
		log.Fatalf("failed to init keycloak: %v", err)
	}

	// --- Connect to Postgres ---
	log.Println("Connecting to Postgres...")
	pg, err := db.ConnectPostgres(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	sqlDB, err := pg.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// --- Run migrations ---
	if err := db.RunMigrations(sqlDB, "internal/db/migrations"); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	// --- Connect to MongoDB ---
	log.Println("Connecting to Mongo...")
	mongoClient, err := db.ConnectMongo(cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	mongoDB := mongoClient.Database("app")

	// --- Ensure Mongo indexes ---
	if err := mongostore.EnsureIndexes(context.Background(), mongoDB); err != nil {
		log.Fatalf("mongo index error: %v", err)
	}

	// --- System handler ---
	systemHandler := handlers.NewSystemHandler()

	// --- Router ---
	router := apphttp.NewRouter(pg, mongoClient, cfg, kc, systemHandler)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- HTTP server ---
	srv := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	// Start server
	go func() {
		log.Println("Server started on :9000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Mark service as ready
	systemHandler.SetReady(true)

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("postgres close error: %v", err)
	}

	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Printf("mongo close error: %v", err)
	}

	log.Println("Server exited gracefully")
}
