package service

import (
"context"
"errors"
"github.com/google/uuid"
"github.com/coolkit-org/coolkit/internal/model"
"github.com/coolkit-org/coolkit/internal/repository"
)

var (
ErrDocumentNotFound = errors.New("document not found")
)

type DocumentService interface {
CreateDocument(ctx context.Context, doc *model.Document) error
GetDocument(ctx context.Context, id uuid.UUID) (*model.Document, error)
GetClubDocuments(ctx context.Context, clubID uuid.UUID) ([]model.Document, error)
GetEventDocuments(ctx context.Context, eventID uuid.UUID) ([]model.Document, error)
GetPermissionLetters(ctx context.Context, clubID uuid.UUID) ([]model.Document, error)
DeleteDocument(ctx context.Context, id uuid.UUID) error
}

type documentService struct {
docRepo  repository.DocumentRepository
memberRepo repository.ClubMemberRepository
}

func NewDocumentService(docRepo repository.DocumentRepository, memberRepo repository.ClubMemberRepository) DocumentService {
return &documentService{docRepo: docRepo, memberRepo: memberRepo}
}

func (s *documentService) CreateDocument(ctx context.Context, doc *model.Document) error {
// Verify user is admin/owner of the club
member, err := s.memberRepo.GetByClubAndUser(ctx, doc.ClubID, doc.UploadedBy)
if err != nil {
return errors.New("user not found in club")
}

if member.Role != model.RoleOwner && member.Role != model.RoleAdmin {
return errors.New("only owners and admins can upload documents")
}

return s.docRepo.Create(ctx, doc)
}

func (s *documentService) GetDocument(ctx context.Context, id uuid.UUID) (*model.Document, error) {
doc, err := s.docRepo.GetByID(ctx, id)
if err != nil {
return nil, ErrDocumentNotFound
}
return doc, nil
}

func (s *documentService) GetClubDocuments(ctx context.Context, clubID uuid.UUID) ([]model.Document, error) {
return s.docRepo.GetByClubID(ctx, clubID)
}

func (s *documentService) GetEventDocuments(ctx context.Context, eventID uuid.UUID) ([]model.Document, error) {
return s.docRepo.GetByEventID(ctx, eventID)
}

func (s *documentService) GetPermissionLetters(ctx context.Context, clubID uuid.UUID) ([]model.Document, error) {
return s.docRepo.GetByType(ctx, clubID, model.DocumentPermissionLetter)
}

func (s *documentService) DeleteDocument(ctx context.Context, id uuid.UUID) error {
_, err := s.docRepo.GetByID(ctx, id)
if err != nil {
return ErrDocumentNotFound
}
return s.docRepo.Delete(ctx, id)
}
