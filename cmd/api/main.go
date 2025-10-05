package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/handlers"
	"github.com/mehmettalhairmak/rss-aggregator/internal/middleware"
)

func main() {
	// Load .env file if it exists
	// Continue even if there's an error (production might not have .env)
	_ = godotenv.Load(".env")

	// Check environment variables
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("$DB_URL must be set")
	}

	// Open database connection
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	defer conn.Close()

	// Create database queries and handler configs
	dbQueries := database.New(conn)
	handlerConfig := handlers.NewConfig(dbQueries)
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

	// User endpoints
	v1Router.Post("/users", handlerConfig.HandlerCreateUser)
	v1Router.Get("/users", middlewareConfig.Auth(handlerConfig.HandlerGetUser))

	// Feed endpoints
	v1Router.Post("/feed", middlewareConfig.Auth(handlerConfig.HandlerCreateFeed))
	v1Router.Get("/feed", handlerConfig.HandlerGetFeed)

	// Feed follows endpoints
	v1Router.Post("/feed_follows", middlewareConfig.Auth(handlerConfig.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", middlewareConfig.Auth(handlerConfig.HandlerGetFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}", middlewareConfig.Auth(handlerConfig.HandlerDeleteFeedFollow))

	// Mount v1Router to main router
	router.Mount("/v1", v1Router)

	// Create and start HTTP server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
