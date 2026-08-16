package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClubID          uuid.UUID `gorm:"type:uuid;not null;index" json:"club_id"`
	Title           string    `gorm:"not null" json:"title"`
	Description     string    `json:"description"`
	Date            time.Time `gorm:"not null" json:"date"`
	Location        string    `json:"location"`
	QRCodeURL       string    `gorm:"type:varchar(500)" json:"qr_code_url"`
	FormbricksEnvID   string    `gorm:"type:varchar(100)" json:"formbricks_env_id"`
	FormbricksSurveyID string   `gorm:"type:varchar(100)" json:"formbricks_survey_id"`
	MatrixRoomID    string    `gorm:"type:varchar(100)" json:"matrix_room_id"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Club Club `gorm:"foreignKey:ClubID" json:"club,omitempty"`
}
