package handler

import (
"github.com/gin-gonic/gin"
"github.com/google/uuid"
"net/http"
"github.com/coolkit-org/coolkit/internal/model"
"github.com/coolkit-org/coolkit/internal/service"
"github.com/coolkit-org/coolkit/pkg/response"
)

type DocumentHandler struct {
docService service.DocumentService
}

func NewDocumentHandler(docService service.DocumentService) *DocumentHandler {
return &DocumentHandler{docService: docService}
}

type CreateDocumentRequest struct {
Title       string `json:"title" binding:"required"`
Description string `json:"description"`
DocType     string `json:"doc_type" binding:"required"`
FileURL     string `json:"file_url" binding:"required"`
EventID     *string `json:"event_id"`
}

func (h *DocumentHandler) CreateDocument(c *gin.Context) {
clubIDStr := c.Param("id")
clubID, err := uuid.Parse(clubIDStr)
if err != nil {
response.BadRequest(c, "Invalid club ID")
return
}

var req CreateDocumentRequest
if err := c.ShouldBindJSON(&req); err != nil {
response.BadRequest(c, err.Error())
return
}

userID, exists := c.Get("user_id")
if !exists {
response.Unauthorized(c, "User not authenticated")
return
}

doc := &model.Document{
ClubID:      clubID,
Title:       req.Title,
Description: req.Description,
DocType:     model.DocumentType(req.DocType),
FileURL:     req.FileURL,
UploadedBy:  userID.(uuid.UUID),
}

if req.EventID != nil {
eventID, err := uuid.Parse(*req.EventID)
if err == nil {
doc.EventID = &eventID
}
}

if err := h.docService.CreateDocument(c.Request.Context(), doc); err != nil {
response.Forbidden(c, err.Error())
return
}

response.Success(c, http.StatusCreated, doc)
}

func (h *DocumentHandler) GetClubDocuments(c *gin.Context) {
clubIDStr := c.Param("id")
clubID, err := uuid.Parse(clubIDStr)
if err != nil {
response.BadRequest(c, "Invalid club ID")
return
}

docs, err := h.docService.GetClubDocuments(c.Request.Context(), clubID)
if err != nil {
response.Error(c, http.StatusInternalServerError, "Failed to fetch documents")
return
}

response.Success(c, http.StatusOK, docs)
}

func (h *DocumentHandler) GetPermissionLetters(c *gin.Context) {
clubIDStr := c.Param("id")
clubID, err := uuid.Parse(clubIDStr)
if err != nil {
response.BadRequest(c, "Invalid club ID")
return
}

docs, err := h.docService.GetPermissionLetters(c.Request.Context(), clubID)
if err != nil {
response.Error(c, http.StatusInternalServerError, "Failed to fetch permission letters")
return
}

response.Success(c, http.StatusOK, docs)
}

func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
docIDStr := c.Param("docId")
docID, err := uuid.Parse(docIDStr)
if err != nil {
response.BadRequest(c, "Invalid document ID")
return
}

if err := h.docService.DeleteDocument(c.Request.Context(), docID); err != nil {
response.Error(c, http.StatusInternalServerError, "Failed to delete document")
return
}

response.Success(c, http.StatusOK, nil)
}
