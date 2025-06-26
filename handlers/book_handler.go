package handlers

import (
	"database/sql"
	"encoding/json"
	"library-api/db"
	"library-api/models"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type BookHandler struct {
	db *sql.DB
}

func NewBookHandler(database *sql.DB) *BookHandler {
	return &BookHandler{db: database}
}

// GetBooks handles GET /api/v1/books
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))

	// Set defaults
	page := 1
	limit := 10

	// Parse page
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit (max 100)
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	var books []models.Book
	var total int
	var err error

	// Search or get all books
	if searchQuery != "" {
		books, total, err = db.SearchBooks(h.db, searchQuery, page, limit)
	} else {
		books, total, err = db.GetBooks(h.db, page, limit)
	}

	if err != nil {
		logrus.WithError(err).Error("Failed to get books")
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve books")
		return
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := models.PaginatedResponse{
		Success: true,
		Data:    books,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	h.sendJSONResponse(w, http.StatusOK, response)
}

// GetBook handles GET /api/v1/books/{id}
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book, err := db.GetBookByID(h.db, id)
	if err != nil {
		logrus.WithError(err).WithField("book_id", id).Error("Failed to get book")
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve book")
		return
	}

	if book == nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	response := models.APIResponse{
		Success: true,
		Data:    book,
	}

	h.sendJSONResponse(w, http.StatusOK, response)
}

// CreateBook handles POST /api/v1/books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Basic validation
	if strings.TrimSpace(req.Title) == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "Title is required")
		return
	}
	if strings.TrimSpace(req.Author) == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "Author is required")
		return
	}
	if req.PublishedYear < 1000 || req.PublishedYear > 2100 {
		h.sendErrorResponse(w, http.StatusBadRequest, "Published year must be between 1000 and 2100")
		return
	}

	// Trim whitespace
	req.Title = strings.TrimSpace(req.Title)
	req.Author = strings.TrimSpace(req.Author)

	book, err := db.CreateBook(h.db, req)
	if err != nil {
		logrus.WithError(err).Error("Failed to create book")
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to create book")
		return
	}

	response := models.APIResponse{
		Success: true,
		Data:    book,
		Message: "Book created successfully",
	}

	h.sendJSONResponse(w, http.StatusCreated, response)
}

// UpdateBook handles PUT /api/v1/books/{id}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var req models.UpdateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Basic validation
	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed == "" {
			h.sendErrorResponse(w, http.StatusBadRequest, "Title cannot be empty")
			return
		}
		req.Title = &trimmed
	}

	if req.Author != nil {
		trimmed := strings.TrimSpace(*req.Author)
		if trimmed == "" {
			h.sendErrorResponse(w, http.StatusBadRequest, "Author cannot be empty")
			return
		}
		req.Author = &trimmed
	}

	if req.PublishedYear != nil {
		if *req.PublishedYear < 1000 || *req.PublishedYear > 2100 {
			h.sendErrorResponse(w, http.StatusBadRequest, "Published year must be between 1000 and 2100")
			return
		}
	}

	book, err := db.UpdateBook(h.db, id, req)
	if err != nil {
		logrus.WithError(err).WithField("book_id", id).Error("Failed to update book")
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to update book")
		return
	}

	if book == nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	response := models.APIResponse{
		Success: true,
		Data:    book,
		Message: "Book updated successfully",
	}

	h.sendJSONResponse(w, http.StatusOK, response)
}

// DeleteBook handles DELETE /api/v1/books/{id}
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	err = db.DeleteBook(h.db, id)
	if err == sql.ErrNoRows {
		h.sendErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}
	if err != nil {
		logrus.WithError(err).WithField("book_id", id).Error("Failed to delete book")
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete book")
		return
	}

	response := models.APIResponse{
		Success: true,
		Message: "Book deleted successfully",
	}

	h.sendJSONResponse(w, http.StatusOK, response)
}

// Helper methods
func (h *BookHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.WithError(err).Error("Failed to encode JSON response")
	}
}

func (h *BookHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := models.APIResponse{
		Success: false,
		Error:   message,
	}

	h.sendJSONResponse(w, statusCode, response)
}
