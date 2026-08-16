package repository

import (
"context"
"github.com/google/uuid"
"github.com/coolkit-org/coolkit/internal/model"
"gorm.io/gorm"
)

type DocumentRepository interface {
Create(ctx context.Context, doc *model.Document) error
GetByID(ctx context.Context, id uuid.UUID) (*model.Document, error)
GetByClubID(ctx context.Context, clubID uuid.UUID) ([]model.Document, error)
GetByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Document, error)
GetByType(ctx context.Context, clubID uuid.UUID, docType model.DocumentType) ([]model.Document, error)
Update(ctx context.Context, doc *model.Document) error
Delete(ctx context.Context, id uuid.UUID) error
}

type documentRepository struct {
db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
return &documentRepository{db: db}
}

func (r *documentRepository) Create(ctx context.Context, doc *model.Document) error {
return r.db.WithContext(ctx).Create(doc).Error
}

func (r *documentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Document, error) {
var doc model.Document
err := r.db.WithContext(ctx).First(&doc, "id = ?", id).Error
if err != nil {
return nil, err
}
return &doc, nil
}

func (r *documentRepository) GetByClubID(ctx context.Context, clubID uuid.UUID) ([]model.Document, error) {
var docs []model.Document
err := r.db.WithContext(ctx).Where("club_id = ?", clubID).Order("created_at DESC").Find(&docs).Error
return docs, err
}

func (r *documentRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Document, error) {
var docs []model.Document
err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Order("created_at DESC").Find(&docs).Error
return docs, err
}

func (r *documentRepository) GetByType(ctx context.Context, clubID uuid.UUID, docType model.DocumentType) ([]model.Document, error) {
var docs []model.Document
err := r.db.WithContext(ctx).Where("club_id = ? AND doc_type = ?", clubID, docType).Order("created_at DESC").Find(&docs).Error
return docs, err
}

func (r *documentRepository) Update(ctx context.Context, doc *model.Document) error {
return r.db.WithContext(ctx).Save(doc).Error
}

func (r *documentRepository) Delete(ctx context.Context, id uuid.UUID) error {
return r.db.WithContext(ctx).Delete(&model.Document{}, "id = ?", id).Error
}
