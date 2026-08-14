package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/coolkit-org/coolkit/internal/config"
	"github.com/coolkit-org/coolkit/internal/database"
	"github.com/coolkit-org/coolkit/internal/model"
)

func main() {
	fmt.Println("🌱 Starting database seeding...")

	// 1. Load Config
	config.Load()
	cfg := config.Get()

	// 2. Connect to Database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Auto Migrate (ensure tables exist)
	fmt.Println("📦 Running auto-migrations...")
	err = db.AutoMigrate(&model.User{}, &model.Club{}, &model.ClubMember{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	// 4. Seed Data
	seedData(db)
}

func seedData(db *gorm.DB) {
	fmt.Println("👤 Seeding users...")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	users := []model.User{
		{
			ID:       uuid.New(),
			Name:     "Alice Johnson",
			Email:    "alice@example.com",
			PasswordHash: string(passwordHash),
		},
		{
			ID:       uuid.New(),
			Name:     "Bob Smith",
			Email:    "bob@example.com",
			PasswordHash: string(passwordHash),
		},
		{
			ID:       uuid.New(),
			Name:     "Charlie Brown",
			Email:    "charlie@example.com",
			PasswordHash: string(passwordHash),
		},
	}

	for i, u := range users {
		var existing model.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&u).Error; err != nil {
					log.Printf("Failed to create user %s: %v", u.Email, err)
				} else {
					fmt.Printf("  ✅ Created user: %s\n", u.Email)
				}
			} else {
				log.Printf("Database error checking user %s: %v", u.Email, err)
			}
		} else {
			fmt.Printf("  ⏭️ User already exists: %s\n", u.Email)
			users[i] = existing // keep correct ID
		}
	}

	fmt.Println("🏢 Seeding clubs...")

	clubs := []model.Club{
		{
			ID:          uuid.New(),
			Name:        "Tech Innovators",
			Description: "A club for technology enthusiasts",
			OwnerID:     users[0].ID, // Alice
		},
		{
			ID:          uuid.New(),
			Name:        "Design Studio",
			Description: "Creative design and UI/UX club",
			OwnerID:     users[1].ID, // Bob
		},
	}

	for i, c := range clubs {
		var existing model.Club
		if err := db.Where("name = ?", c.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&c).Error; err != nil {
					log.Printf("Failed to create club %s: %v", c.Name, err)
				} else {
					fmt.Printf("  ✅ Created club: %s\n", c.Name)
				}
			} else {
				log.Printf("Database error checking club %s: %v", c.Name, err)
			}
		} else {
			fmt.Printf("  ⏭️ Club already exists: %s\n", c.Name)
			clubs[i] = existing
		}
	}

	fmt.Println("🤝 Seeding memberships...")

	memberships := []model.ClubMember{
		{ClubID: clubs[0].ID, UserID: users[0].ID, Role: "owner", JoinedAt: time.Now()},     // Alice -> Tech
		{ClubID: clubs[0].ID, UserID: users[1].ID, Role: "member", JoinedAt: time.Now()},    // Bob -> Tech
		{ClubID: clubs[0].ID, UserID: users[2].ID, Role: "member", JoinedAt: time.Now()},    // Charlie -> Tech
		{ClubID: clubs[1].ID, UserID: users[1].ID, Role: "owner", JoinedAt: time.Now()},     // Bob -> Design
		{ClubID: clubs[1].ID, UserID: users[0].ID, Role: "member", JoinedAt: time.Now()},    // Alice -> Design
		{ClubID: clubs[1].ID, UserID: users[2].ID, Role: "member", JoinedAt: time.Now()},    // Charlie -> Design
	}

	for _, m := range memberships {
		var existing model.ClubMember
		if err := db.Where("club_id = ? AND user_id = ?", m.ClubID, m.UserID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&m).Error; err != nil {
					log.Printf("Failed to create membership for user %v in club %v: %v", m.UserID, m.ClubID, err)
				} else {
					fmt.Printf("  ✅ Created membership for user in club\n")
				}
			} else {
				log.Printf("Database error checking membership: %v", err)
			}
		} else {
			fmt.Printf("  ⏭️ Membership already exists\n")
		}
	}

	fmt.Println("🎉 Seeding complete!")
}
