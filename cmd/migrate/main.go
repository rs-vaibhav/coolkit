package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coolkit-org/coolkit/internal/config"
	"github.com/coolkit-org/coolkit/internal/database"
	"github.com/coolkit-org/coolkit/internal/model"
)

func main() {
	config.Load()
	cfg := config.Get()
	if cfg == nil {
		log.Fatalf("Failed to load config")
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Running database migrations...")

	// Auto-migrate all models
	err = db.AutoMigrate(
		&model.User{},
		&model.Club{},
		&model.Event{},
		&model.Task{},
		&model.Announcement{},
		&model.Booking{},
		&model.Resource{},
		&model.Domain{},
		&model.JoinRequest{},
		&model.EventRole{},
		&model.FinanceEntry{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("✅ Database migration completed successfully!")
	fmt.Println("New fields added:")
	fmt.Println("  - Club: profile_image, banner_image")
	fmt.Println("  - FinanceEntry: proof_image, status, approved_by_id, approved_at, rejection_reason")
	
	// Verify new columns exist
	if db.Migrator().HasColumn(&model.Club{}, "profile_image") {
		fmt.Println("  ✓ Club.profile_image column exists")
	}
	if db.Migrator().HasColumn(&model.Club{}, "banner_image") {
		fmt.Println("  ✓ Club.banner_image column exists")
	}
	if db.Migrator().HasColumn(&model.FinanceEntry{}, "proof_image") {
		fmt.Println("  ✓ FinanceEntry.proof_image column exists")
	}
	if db.Migrator().HasColumn(&model.FinanceEntry{}, "status") {
		fmt.Println("  ✓ FinanceEntry.status column exists")
	}
	if db.Migrator().HasColumn(&model.FinanceEntry{}, "approved_by_id") {
		fmt.Println("  ✓ FinanceEntry.approved_by_id column exists")
	}

	os.Exit(0)
}
