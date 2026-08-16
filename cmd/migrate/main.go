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

	// Enable uuid-ossp extension for UUID generation
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	fmt.Println("Running database migrations...")

	// Auto-migrate all models (must match cmd/server/main.go)
	err = db.AutoMigrate(
		&model.User{},
		&model.Club{},
		&model.ClubMember{},
		&model.Event{},
		&model.EventRole{},
		&model.Announcement{},
		&model.Task{},
		&model.FinanceEntry{},
		&model.JoinRequest{},
		&model.HierarchyLevel{},
		&model.Domain{},
		&model.Booking{},
		&model.Resource{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Backfill existing finance entries that have no status set
	result := db.Exec("UPDATE finance_entries SET status = 'approved' WHERE status IS NULL OR status = ''")
	if result.Error != nil {
		log.Printf("Warning: failed to backfill finance entries: %v", result.Error)
	} else if result.RowsAffected > 0 {
		fmt.Printf("  Backfilled %d existing finance entries to 'approved' status\n", result.RowsAffected)
	}

	fmt.Println("Database migration completed successfully!")
	fmt.Println("New fields added:")
	fmt.Println("  - Club: profile_image, banner_image")
	fmt.Println("  - FinanceEntry: proof_image, status, approved_by_id, approved_at, rejection_reason")

	// Verify new columns exist
	if db.Migrator().HasColumn(&model.Club{}, "profile_image") {
		fmt.Println("  - Club.profile_image column exists")
	}
	if db.Migrator().HasColumn(&model.Club{}, "banner_image") {
		fmt.Println("  - Club.banner_image column exists")
	}
	if db.Migrator().HasColumn(&model.FinanceEntry{}, "proof_image") {
		fmt.Println("  - FinanceEntry.proof_image column exists")
	}
	if db.Migrator().HasColumn(&model.FinanceEntry{}, "status") {
		fmt.Println("  - FinanceEntry.status column exists")
	}
	if db.Migrator().HasColumn(&model.FinanceEntry{}, "approved_by_id") {
		fmt.Println("  - FinanceEntry.approved_by_id column exists")
	}

	os.Exit(0)
}

