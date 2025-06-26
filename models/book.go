package models

import (
	"time"
)

// Book represents a book in the library
type Book struct {
	ID            int       `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Author        string    `json:"author" db:"author"`
	PublishedYear int       `json:"published_year" db:"published_year"`
	Available     bool      `json:"available" db:"available"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// CreateBookRequest represents the request payload for creating a book
type CreateBookRequest struct {
	Title         string `json:"title" validate:"required,min=1,max=255"`
	Author        string `json:"author" validate:"required,min=1,max=255"`
	PublishedYear int    `json:"published_year" validate:"required,min=1000,max=2100"`
	Available     *bool  `json:"available,omitempty"`
}

// UpdateBookRequest represents the request payload for updating a book
type UpdateBookRequest struct {
	Title         *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Author        *string `json:"author,omitempty" validate:"omitempty,min=1,max=255"`
	PublishedYear *int    `json:"published_year,omitempty" validate:"omitempty,min=1000,max=2100"`
	Available     *bool   `json:"available,omitempty"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Error      string      `json:"error,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
