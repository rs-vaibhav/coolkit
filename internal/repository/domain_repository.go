package repository

import (
	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DomainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

func (r *DomainRepository) Create(domain *model.Domain) error {
	return r.db.Create(domain).Error
}

func (r *DomainRepository) FindByClubID(clubID uuid.UUID) ([]model.Domain, error) {
	var domains []model.Domain
	err := r.db.Where("club_id = ?", clubID).Order("name ASC").Find(&domains).Error
	return domains, err
}

func (r *DomainRepository) FindByID(id uuid.UUID) (*model.Domain, error) {
	var domain model.Domain
	err := r.db.First(&domain, id).Error
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

func (r *DomainRepository) Update(domain *model.Domain) error {
	return r.db.Model(domain).Select("Name", "Description").Updates(domain).Error
}

func (r *DomainRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Domain{}, id).Error
}
