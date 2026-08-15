package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coolkit-org/coolkit/internal/config"
	"github.com/coolkit-org/coolkit/internal/database"
	"github.com/coolkit-org/coolkit/internal/handler"
	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
	"github.com/coolkit-org/coolkit/internal/router"
	"github.com/coolkit-org/coolkit/internal/service"
)

func main() {
	// Load config
	config.Load()
	cfg := config.Get()

	// Connect to DB
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Enable uuid-ossp extension for UUID generation
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	// Auto-migrate models
	err = db.AutoMigrate(&model.User{}, &model.Club{}, &model.ClubMember{}, &model.Event{}, &model.EventRole{}, &model.Announcement{}, &model.Task{}, &model.FinanceEntry{}, &model.JoinRequest{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	log.Println("Database migrated successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	clubRepo := repository.NewClubRepository(db)
	eventRepo := repository.NewEventRepository(db)
	eventRoleRepo := repository.NewEventRoleRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	financeRepo := repository.NewFinanceRepository(db)
	joinRequestRepo := repository.NewJoinRequestRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	clubService := service.NewClubService(clubRepo, joinRequestRepo)
	eventService := service.NewEventService(eventRepo, clubRepo, eventRoleRepo)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(authService)
	clubHandler := handler.NewClubHandler(clubService)
	memberHandler := handler.NewMemberHandler(clubService)
	eventHandler := handler.NewEventHandler(eventService)
	eventRoleHandler := handler.NewEventRoleHandler(eventService)
	announcementHandler := handler.NewAnnouncementHandler(announcementRepo, clubRepo)
	taskHandler := handler.NewTaskHandler(taskRepo, eventRepo, clubRepo)
	financeHandler := handler.NewFinanceHandler(financeRepo, eventRepo, clubRepo)

	// Set up router
	r := router.Setup(healthHandler, authHandler, clubHandler, memberHandler, eventHandler, eventRoleHandler, announcementHandler, taskHandler, financeHandler, cfg.JWTSecret)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server
	go func() {
		log.Printf("🚀 CoolKit server starting on http://0.0.0.0:%s", cfg.Port)
		log.Printf("   Mode: %s", cfg.GinMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
