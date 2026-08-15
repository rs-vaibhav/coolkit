package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/coolkit-org/coolkit/internal/model"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) FindByEventID(eventID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("event_id = ?", eventID).Order("due_date ASC NULLS LAST").Preload("AssignedTo").Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) FindByID(id uuid.UUID) (*model.Task, error) {
	var task model.Task
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TaskRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Task{}, id).Error
}
