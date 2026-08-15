package service

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
)

type AnalyticsService struct {
	db        *gorm.DB
	clubRepo  *repository.ClubRepository
	eventRepo *repository.EventRepository
}

func NewAnalyticsService(
	db *gorm.DB,
	clubRepo *repository.ClubRepository,
	eventRepo *repository.EventRepository,
) *AnalyticsService {
	return &AnalyticsService{
		db:        db,
		clubRepo:  clubRepo,
		eventRepo: eventRepo,
	}
}

func (s *AnalyticsService) GetClubAnalytics(clubID, userID uuid.UUID) (*model.ClubAnalytics, error) {
	// Verify membership
	isMem, err := s.clubRepo.IsMember(clubID, userID)
	if err != nil {
		return nil, err
	}
	if !isMem {
		return nil, ErrNotAuthorized
	}

	// 1. Member Growth (cumulative members over time)
	members, err := s.clubRepo.FindMembers(clubID)
	if err != nil {
		return nil, err
	}

	monthlyMap := make(map[string]int64)
	for _, m := range members {
		month := m.JoinedAt.Format("2006-02") // format as YYYY-MM
		monthlyMap[month]++
	}

	// Sort months chronologically
	var months []string
	for k := range monthlyMap {
		months = append(months, k)
	}
	sort.Strings(months)

	// Build cumulative counts
	var memberGrowth []model.MonthlyMemberCount
	var cumulativeCount int64
	for _, m := range months {
		cumulativeCount += monthlyMap[m]
		memberGrowth = append(memberGrowth, model.MonthlyMemberCount{
			Month: m,
			Count: cumulativeCount,
		})
	}

	// 2. Events Stats
	events, err := s.eventRepo.FindByClubID(clubID)
	if err != nil {
		return nil, err
	}

	var totalEvents, upcomingEvents, completedEvents int64
	now := time.Now()
	for _, e := range events {
		totalEvents++
		if e.Date.After(now) {
			upcomingEvents++
		} else {
			completedEvents++
		}
	}

	eventStats := model.EventAnalytics{
		TotalEvents:     totalEvents,
		UpcomingEvents:  upcomingEvents,
		CompletedEvents: completedEvents,
	}

	// 3. Finances (aggregate from events)
	var entries []model.FinanceEntry
	err = s.db.Joins("JOIN events ON events.id = finance_entries.event_id").
		Where("events.club_id = ?", clubID).
		Find(&entries).Error
	if err != nil {
		return nil, err
	}

	var totalIncome, totalExpense float64
	catMap := make(map[string]float64)
	catType := make(map[string]string)

	for _, e := range entries {
		if e.Type == model.FinanceTypeIncome {
			totalIncome += e.Amount
		} else if e.Type == model.FinanceTypeExpense {
			totalExpense += e.Amount
		}

		key := e.Category + "_" + e.Type
		catMap[key] += e.Amount
		catType[key] = e.Type
	}

	var categories []model.CategoryStat
	for key, amount := range catMap {
		// key is in format Category_Type
		// Let's parse it
		var category, ftype string
		for i := len(key) - 1; i >= 0; i-- {
			if key[i] == '_' {
				category = key[:i]
				ftype = key[i+1:]
				break
			}
		}
		categories = append(categories, model.CategoryStat{
			Category: category,
			Type:     ftype,
			Amount:   amount,
		})
	}

	financeStats := model.FinanceAnalytics{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
		Categories:   categories,
	}

	// 4. Tasks (aggregate from events)
	var tasks []model.Task
	err = s.db.Joins("JOIN events ON events.id = tasks.event_id").
		Where("events.club_id = ?", clubID).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	var totalTasks, todoTasks, inProgressTasks, doneTasks int64
	for _, t := range tasks {
		totalTasks++
		switch t.Status {
		case model.TaskStatusTodo:
			todoTasks++
		case model.TaskStatusInProgress:
			inProgressTasks++
		case model.TaskStatusDone:
			doneTasks++
		}
	}

	taskStats := model.TaskAnalytics{
		Total:      totalTasks,
		Todo:       todoTasks,
		InProgress: inProgressTasks,
		Done:       doneTasks,
	}

	return &model.ClubAnalytics{
		MemberGrowth: memberGrowth,
		EventStats:   eventStats,
		FinanceStats: financeStats,
		TaskStats:    taskStats,
	}, nil
}
