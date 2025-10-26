package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/handlers"
	"github.com/mehmettalhairmak/rss-aggregator/internal/logger"
	"github.com/mehmettalhairmak/rss-aggregator/internal/middleware"
	"github.com/mehmettalhairmak/rss-aggregator/internal/scraper"
)

func main() {
	// Initialize logger first
	logger.InitLogger()

	// Load .env file if it exists
	// Continue even if there's an error (production might not have .env)
	_ = godotenv.Load(".env")

	// Check environment variables
	portString := os.Getenv("PORT")
	if portString == "" {
		logger.Fatal("$PORT environment variable must be set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		logger.Fatal("$DB_URL environment variable must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Fatal("$JWT_SECRET environment variable must be set")
	}

	logger.Infof("Starting RSS Aggregator API on port %s", portString)

	// Open database connection
	logger.Info("Connecting to database...")
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.ErrorErr(err, "Failed to connect to database")
		os.Exit(1)
	}
	defer conn.Close()
	logger.Info("Successfully connected to database")

	// Create database queries and handler configs
	dbQueries := database.New(conn)
	handlerConfig := handlers.NewConfig(dbQueries, conn)
	middlewareConfig := middleware.NewConfig(dbQueries)

	// Create Chi router
	router := chi.NewRouter()

	// Add CORS middleware
	// CORS: Cross-Origin Resource Sharing - allows API requests from different domains
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Create v1 API router
	// Using versioning - we can add v2 in the future
	v1Router := chi.NewRouter()

	// Health check endpoints
	v1Router.Get("/ready", handlers.HandlerReadiness)
	v1Router.Get("/error", handlers.HandlerErr)

	// Authentication endpoints (Public - no auth required)
	// POST /v1/auth/register
	// POST /v1/auth/login
	v1Router.Post("/auth/register", handlerConfig.HandlerRegister)
	v1Router.Post("/auth/login", handlerConfig.HandlerLogin)
	v1Router.Post("/auth/refresh", handlerConfig.HandlerRefreshToken)
	v1Router.Get("/auth/logout", middlewareConfig.Auth(handlerConfig.HandlerLogout))

	// User endpoints (Protected - JWT required)
	// GET /v1/users/me - Returns the authenticated user's information
	v1Router.Get("/users/me", middlewareConfig.Auth(handlerConfig.HandlerGetUser))

	// Feed endpoints
	v1Router.Post("/feed", middlewareConfig.Auth(handlerConfig.HandlerCreateFeed))
	v1Router.Get("/feed", handlerConfig.HandlerGetFeed)

	// Feed follows endpoints
	v1Router.Post("/feed_follows", middlewareConfig.Auth(handlerConfig.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", middlewareConfig.Auth(handlerConfig.HandlerGetFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}", middlewareConfig.Auth(handlerConfig.HandlerDeleteFeedFollow))

	// Posts endpoints
	v1Router.Get("/posts", middlewareConfig.Auth(handlerConfig.HandlerGetUserPostsForUser))

	// Mount v1Router to main router
	router.Mount("/v1", v1Router)

	// Start background scraper
	logger.Info("Starting RSS feed scraper...")
	go scraper.StartScraping(dbQueries, 10, time.Minute)

	// Create and start HTTP server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	logger.Infof("Server starting on port %s", portString)
	if err := srv.ListenAndServe(); err != nil {
		logger.ErrorErr(err, "Server failed to start")
		os.Exit(1)
	}
}
