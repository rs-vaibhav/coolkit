package service

import (
	"database/sql/driver"
	"testing"
	"time"

	gosqlite "github.com/glebarez/go-sqlite"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
)

func init() {
	gosqlite.MustRegisterDeterministicScalarFunction("uuid_generate_v4", 0, func(ctx *gosqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		return uuid.New().String(), nil
	})
}

func TestPhase3Integration(t *testing.T) {
	// 1. Initialize in-memory SQLite DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}

	// 2. Create tables manually for SQLite
	tables := []string{
		`CREATE TABLE users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			avatar_url TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE clubs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			join_code TEXT UNIQUE NOT NULL,
			owner_id TEXT NOT NULL,
			owner_label TEXT DEFAULT 'Owner',
			admin_label TEXT DEFAULT 'Admin',
			leadership_label TEXT DEFAULT 'Leadership',
			profile_image TEXT,
			banner_image TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE club_members (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'member',
			domain_id TEXT,
			hierarchy_level_id TEXT,
			joined_at DATETIME
		)`,
		`CREATE TABLE events (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			date DATETIME NOT NULL,
			location TEXT,
			qr_code_url TEXT,
			formbricks_env_id TEXT,
			matrix_room_id TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE finance_entries (
			id TEXT PRIMARY KEY,
			event_id TEXT NOT NULL,
			created_by_id TEXT NOT NULL,
			type TEXT NOT NULL,
			category TEXT NOT NULL,
			amount REAL NOT NULL,
			description TEXT,
			date DATETIME,
			proof_image TEXT,
			status TEXT DEFAULT 'approved',
			approved_by_id TEXT,
			approved_at DATETIME,
			rejection_reason TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE tasks (
			id TEXT PRIMARY KEY,
			event_id TEXT NOT NULL,
			assigned_to_id TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL DEFAULT 'todo',
			due_date DATETIME,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE resources (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			quantity INTEGER NOT NULL DEFAULT 1,
			created_by_id TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE bookings (
			id TEXT PRIMARY KEY,
			resource_id TEXT NOT NULL,
			club_id TEXT NOT NULL,
			booked_by_id TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			purpose TEXT,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at DATETIME,
			updated_at DATETIME
		)`,
	}

	for _, sql := range tables {
		if err := db.Exec(sql).Error; err != nil {
			t.Fatalf("Failed to create table: %v", err)
		}
	}

	// 3. Initialize repositories
	clubRepo := repository.NewClubRepository(db)
	eventRepo := repository.NewEventRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	financeRepo := repository.NewFinanceRepository(db)
	resourceRepo := repository.NewResourceRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// 4. Initialize services
	bookingService := NewBookingService(clubRepo, resourceRepo, bookingRepo)
	analyticsService := NewAnalyticsService(db, clubRepo, eventRepo)

	// 5. Create test users
	alice := &model.User{
		ID:           uuid.New(),
		Name:         "Alice Johnson",
		Email:        "alice@example.com",
		PasswordHash: "hashed_password",
	}
	bob := &model.User{
		ID:           uuid.New(),
		Name:         "Bob Smith",
		Email:        "bob@example.com",
		PasswordHash: "hashed_password",
	}
	assert.NoError(t, db.Create(alice).Error)
	assert.NoError(t, db.Create(bob).Error)

	// 6. Create Club (Alice is owner)
	club := &model.Club{
		ID:          uuid.New(),
		Name:        "FOSS Club",
		Description: "SRM FOSS Club",
		JoinCode:    "FOSS123",
		OwnerID:     alice.ID,
	}
	assert.NoError(t, clubRepo.Create(club))

	// Add Alice as owner membership
	aliceMember := &model.ClubMember{
		ID:       uuid.New(),
		ClubID:   club.ID,
		UserID:   alice.ID,
		Role:     model.RoleOwner,
		JoinedAt: time.Now().Add(-48 * time.Hour), // Joined 2 days ago
	}
	assert.NoError(t, clubRepo.AddMember(aliceMember))

	// Add Bob as standard member
	bobMember := &model.ClubMember{
		ID:       uuid.New(),
		ClubID:   club.ID,
		UserID:   bob.ID,
		Role:     model.RoleMember,
		JoinedAt: time.Now().Add(-24 * time.Hour), // Joined 1 day ago
	}
	assert.NoError(t, clubRepo.AddMember(bobMember))

	// 7. Create an Event
	event := &model.Event{
		ID:          uuid.New(),
		ClubID:      club.ID,
		Title:       "FOSS Hackathon 2026",
		Description: "Annual hackathon event",
		Date:        time.Now().Add(24 * time.Hour), // Upcoming event
		Location:    "Tech Park",
	}
	assert.NoError(t, eventRepo.Create(event))

	// 8. Log Finance entries
	income := &model.FinanceEntry{
		ID:          uuid.New(),
		EventID:     event.ID,
		CreatedByID: alice.ID,
		Type:        model.FinanceTypeIncome,
		Category:    "Sponsorship",
		Amount:      10000.0,
		Description: "Sponsorship from Tech Corp",
		Date:        time.Now(),
		Status:      model.FinanceStatusApproved,
	}
	expense := &model.FinanceEntry{
		ID:          uuid.New(),
		EventID:     event.ID,
		CreatedByID: alice.ID,
		Type:        model.FinanceTypeExpense,
		Category:    "Food",
		Amount:      3000.0,
		Description: "Pizza order for participants",
		Date:        time.Now(),
		Status:      model.FinanceStatusApproved,
	}
	assert.NoError(t, financeRepo.Create(income))
	assert.NoError(t, financeRepo.Create(expense))

	// 9. Log Tasks
	task := &model.Task{
		ID:           uuid.New(),
		EventID:      event.ID,
		AssignedToID: bob.ID,
		Title:        "Design banners",
		Description:  "Create marketing assets",
		Status:       model.TaskStatusDone, // Completed task
	}
	assert.NoError(t, taskRepo.Create(task))

	// 10. Test Resource Booking Flow
	// Add resource (AV Projector, quantity = 1)
	resource, err := bookingService.AddResource("AV Projector", "Full HD Sony Projector", 1, club.ID, alice.ID)
	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Book resource (10 AM to 12 PM tomorrow)
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(2 * time.Hour)
	booking, err := bookingService.BookResource(resource.ID, club.ID, bob.ID, start, end, "Presentation session")
	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, model.BookingStatusPending, booking.Status)

	// Approve booking (By Alice, Owner)
	err = bookingService.ApproveBooking(booking.ID, alice.ID)
	assert.NoError(t, err)

	// Verify status is approved
	bookingCheck, err := bookingRepo.FindByID(booking.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.BookingStatusApproved, bookingCheck.Status)

	// Try booking overlapping timeframe (11 AM to 1 PM tomorrow) - Should fail on approval due to overlap constraint
	startOverlap := start.Add(1 * time.Hour)
	endOverlap := startOverlap.Add(2 * time.Hour)
	bookingOverlap, err := bookingService.BookResource(resource.ID, club.ID, bob.ID, startOverlap, endOverlap, "Overlapping booking")
	assert.NoError(t, err)

	err = bookingService.ApproveBooking(bookingOverlap.ID, alice.ID)
	assert.Error(t, err) // Should return "resource is fully booked during this time slot"
	assert.Contains(t, err.Error(), "fully booked")

	// 11. Test Analytics Service
	analytics, err := analyticsService.GetClubAnalytics(club.ID, alice.ID)
	assert.NoError(t, err)
	assert.NotNil(t, analytics)

	// Validate members analytics (Growth cumulative total = 2 members)
	assert.NotEmpty(t, analytics.MemberGrowth)
	assert.Equal(t, int64(2), analytics.MemberGrowth[len(analytics.MemberGrowth)-1].Count)

	// Validate events analytics (Total = 1 event, Upcoming = 1 event)
	assert.Equal(t, int64(1), analytics.EventStats.TotalEvents)
	assert.Equal(t, int64(1), analytics.EventStats.UpcomingEvents)

	// Validate finance analytics (Income = 10000, Expense = 3000, Balance = 7000)
	assert.Equal(t, 10000.0, analytics.FinanceStats.TotalIncome)
	assert.Equal(t, 3000.0, analytics.FinanceStats.TotalExpense)
	assert.Equal(t, 7000.0, analytics.FinanceStats.Balance)
	assert.Len(t, analytics.FinanceStats.Categories, 2)

	// Validate tasks analytics (Total = 1 task, Done = 1 task)
	assert.Equal(t, int64(1), analytics.TaskStats.Total)
	assert.Equal(t, int64(1), analytics.TaskStats.Done)
}
