package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"ArtoFino/backend/internal/auth"
	"ArtoFino/backend/internal/config"
	"ArtoFino/backend/internal/http/handlers"
	"ArtoFino/backend/internal/http/middleware"
)

// NewRouter initializes all HTTP routes and dependencies.
func NewRouter(
	pg *gorm.DB,
	mongo *mongo.Client,
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

	objectsHandler := handlers.NewObjectsHandler()
	r.GET("/objects/:id", objectsHandler.GetArtObject)

	// --- Auth middleware ---
	authMw := middleware.NewAuthMiddleware(kc)

	// --- Protected routes ---
	protected := r.Group("/")
	protected.Use(authMw.Handle)

	usersHandler := handlers.NewUsersHandler()
	transactionsHandler := handlers.NewTransactionsHandler()
	transfersHandler := handlers.NewTransfersHandler()
	authorsHandler := handlers.NewAuthorsHandler()

	// /users/me доступен только роли "user"
	protected.GET(
		"/users/me",
		middleware.RequireRole("user", cfg.Keycloak.ClientID),
		usersHandler.Me,
	)

	protected.POST("/transactions/buy", transactionsHandler.BuyShare)
	protected.POST("/transfers/request", transfersHandler.RequestTransfer)
	protected.POST("/transfers/:id/approve", transfersHandler.ApproveTransfer)
	protected.POST("/authors/apply", authorsHandler.ApplyForCreator)

	// --- Admin section ---
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole("admin", cfg.Keycloak.ClientID))

	adminHandler := handlers.NewAdminHandler()

	admin.GET("/stats", adminHandler.Stats)
	admin.GET("/ping", adminHandler.Ping)

	admin.POST("/applications/:id/approve", adminHandler.ApproveArtistApplication)

	return r
}
