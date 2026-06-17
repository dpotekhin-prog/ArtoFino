package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"ArtoFino/backend/internal/auth"
	"ArtoFino/backend/internal/config"
	pgrepo "ArtoFino/backend/internal/db"
	"ArtoFino/backend/internal/http/handlers"
	"ArtoFino/backend/internal/http/middleware"
	mongorepo "ArtoFino/backend/internal/mongo"
)

// NewRouter initializes all HTTP routes and dependencies.
func NewRouter(
	pg *gorm.DB,
	mongoClient *mongo.Client,
	cfg config.Config,
	kc *auth.KeycloakClient,
	system *handlers.SystemHandler,
) *gin.Engine {

	r := gin.Default()

	// --- Healthcheck ---
	r.GET("/hc", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// --- Initialize Repositories ---
	// 1. MongoDB Setup
	mongoDB := mongoClient.Database("artofino")
	objectsRepo := mongorepo.NewObjectsRepository(mongoDB)

	// 2. PostgreSQL Setup
	transactionsPostgresRepo := pgrepo.NewTransactionsRepository(pg)
	transfersPostgresRepo := pgrepo.NewTransfersRepository(pg)

	// --- Initialize Handlers with Dependencies ---
	objectsHandler := handlers.NewObjectsHandler(objectsRepo)
	usersHandler := handlers.NewUsersHandler()
	transactionsHandler := handlers.NewTransactionsHandler(objectsRepo, transactionsPostgresRepo)
	transfersHandler := handlers.NewTransfersHandler(objectsRepo, transfersPostgresRepo)
	authorsHandler := handlers.NewAuthorsHandler()
	adminHandler := handlers.NewAdminHandler(objectsRepo)

	// --- Public Routes ---
	r.GET("/objects/:id", objectsHandler.GetArtObject)

	// --- Auth Middleware ---
	authMw := middleware.NewAuthMiddleware(kc)

	// --- Protected Routes ---
	protected := r.Group("/")
	protected.Use(authMw.Handle)

	// /users/me available only to "user" role
	protected.GET(
		"/users/me",
		middleware.RequireRole("user", cfg.Keycloak.ClientID),
		usersHandler.Me,
	)

	protected.POST("/transactions/buy", transactionsHandler.BuyShare)
	protected.POST("/transfers/request", transfersHandler.RequestTransfer)
	protected.POST("/transfers/:id/approve", transfersHandler.ApproveTransfer)
	protected.POST("/authors/apply", authorsHandler.ApplyForCreator)

	// --- Admin Section ---
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole("admin", cfg.Keycloak.ClientID))

	admin.GET("/stats", adminHandler.Stats)
	admin.GET("/ping", adminHandler.Ping)
	admin.POST("/applications/:id/approve", adminHandler.ApproveArtistApplication)
	admin.POST("/objects", adminHandler.CreateArtObject)

	return r
}
