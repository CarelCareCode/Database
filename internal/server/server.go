package server

import (
	"context"
	"emergency-response-backend/internal/config"
	"emergency-response-backend/internal/database"
	"emergency-response-backend/internal/handlers"
	"emergency-response-backend/internal/middleware"
	"emergency-response-backend/internal/redis"
	"emergency-response-backend/internal/websocket"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	config   *config.Config
	db       *database.DB
	redis    *redis.Client
	wsHub    *websocket.Hub
	server   *http.Server
	handlers *handlers.Handlers
}

func New(cfg *config.Config, db *database.DB, redisClient *redis.Client, wsHub *websocket.Hub) *Server {
	return &Server{
		config:   cfg,
		db:       db,
		redis:    redisClient,
		wsHub:    wsHub,
		handlers: handlers.New(db, redisClient, wsHub),
	}
}

func (s *Server) Start() error {
	router := s.setupRoutes()

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Configure properly for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *Server) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Add middleware
	router.Use(middleware.Logging)
	router.Use(middleware.Recovery)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Authentication routes
	api.HandleFunc("/register", s.handlers.Register).Methods("POST")
	api.HandleFunc("/login", s.handlers.Login).Methods("POST")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.JWTAuth(s.config.JWT.Secret))

	// User routes
	protected.HandleFunc("/profile", s.handlers.GetProfile).Methods("GET")
	protected.HandleFunc("/medical", s.handlers.CreateMedicalProfile).Methods("POST")
	protected.HandleFunc("/medical/{user_id}", s.handlers.GetMedicalProfile).Methods("GET")

	// Emergency routes
	protected.HandleFunc("/emergency", s.handlers.CreateEmergency).Methods("POST")
	protected.HandleFunc("/incidents", s.handlers.GetIncidents).Methods("GET")
	protected.HandleFunc("/incidents/{id}", s.handlers.GetIncident).Methods("GET")
	protected.HandleFunc("/incidents/{id}/assign", s.handlers.AssignParamedic).Methods("POST")
	protected.HandleFunc("/incidents/{id}/status", s.handlers.UpdateIncidentStatus).Methods("PUT")

	// Chat routes
	protected.HandleFunc("/chat", s.handlers.SendMessage).Methods("POST")
	protected.HandleFunc("/chat/{incident_id}", s.handlers.GetChatHistory).Methods("GET")

	// Paramedic routes
	protected.HandleFunc("/paramedics", s.handlers.GetParamedics).Methods("GET")
	protected.HandleFunc("/paramedics/location", s.handlers.UpdateParamedicLocation).Methods("POST")

	// WebSocket endpoint
	protected.HandleFunc("/ws", s.handlers.HandleWebSocket).Methods("GET")

	// Health check
	router.HandleFunc("/health", s.handlers.HealthCheck).Methods("GET")

	return router
}
