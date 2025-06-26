package db

import (
	"database/sql"
	"fmt"
	"library-api/models"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// InitDB initializes the database connection
func InitDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "user"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "Password"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "db"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("Successfully connected to database")
	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS books (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			published_year INT NOT NULL,
			available BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_title (title),
			INDEX idx_author (author),
			INDEX idx_published_year (published_year),
			INDEX idx_available (available)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", i+1, err)
		}
	}

	logrus.Info("Database migrations completed successfully")
	return nil
}

// GetBooks retrieves books with pagination
func GetBooks(db *sql.DB, page, limit int) ([]models.Book, int, error) {
	// Get total count
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM books").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get books with pagination
	query := `SELECT id, title, author, published_year, available, created_at, updated_at 
			  FROM books 
			  ORDER BY created_at DESC 
			  LIMIT ? OFFSET ?`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query books: %w", err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear,
			&book.Available, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over rows: %w", err)
	}

	return books, total, nil
}

// GetBookByID retrieves a single book by ID
func GetBookByID(db *sql.DB, id int) (*models.Book, error) {
	query := `SELECT id, title, author, published_year, available, created_at, updated_at 
			  FROM books WHERE id = ?`

	var book models.Book
	err := db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author,
		&book.PublishedYear, &book.Available, &book.CreatedAt, &book.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return &book, nil
}

// CreateBook creates a new book
func CreateBook(db *sql.DB, req models.CreateBookRequest) (*models.Book, error) {
	available := true
	if req.Available != nil {
		available = *req.Available
	}

	query := `INSERT INTO books (title, author, published_year, available) 
			  VALUES (?, ?, ?, ?)`

	result, err := db.Exec(query, req.Title, req.Author, req.PublishedYear, available)
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return GetBookByID(db, int(id))
}

// UpdateBook updates an existing book
func UpdateBook(db *sql.DB, id int, req models.UpdateBookRequest) (*models.Book, error) {
	// Check if book exists
	existing, err := GetBookByID(db, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}

	if req.Title != nil {
		updates = append(updates, "title = ?")
		args = append(args, *req.Title)
	}
	if req.Author != nil {
		updates = append(updates, "author = ?")
		args = append(args, *req.Author)
	}
	if req.PublishedYear != nil {
		updates = append(updates, "published_year = ?")
		args = append(args, *req.PublishedYear)
	}
	if req.Available != nil {
		updates = append(updates, "available = ?")
		args = append(args, *req.Available)
	}

	if len(updates) == 0 {
		return existing, nil // No updates needed
	}

	query := fmt.Sprintf("UPDATE books SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		fmt.Sprintf("%s", updates[0]))
	for i := 1; i < len(updates); i++ {
		query = fmt.Sprintf("UPDATE books SET %s, %s, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			updates[0], updates[i])
	}

	// Rebuild query properly
	updateClause := ""
	for i, update := range updates {
		if i > 0 {
			updateClause += ", "
		}
		updateClause += update
	}
	query = fmt.Sprintf("UPDATE books SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ?", updateClause)

	args = append(args, id)

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return GetBookByID(db, id)
}

// DeleteBook deletes a book by ID
func DeleteBook(db *sql.DB, id int) error {
	// Check if book exists
	existing, err := GetBookByID(db, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return sql.ErrNoRows
	}

	query := "DELETE FROM books WHERE id = ?"
	_, err = db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

// SearchBooks searches for books by title or author
func SearchBooks(db *sql.DB, query string, page, limit int) ([]models.Book, int, error) {
	searchTerm := "%" + query + "%"

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM books WHERE title LIKE ? OR author LIKE ?"
	err := db.QueryRow(countQuery, searchTerm, searchTerm).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get books with search and pagination
	searchQuery := `SELECT id, title, author, published_year, available, created_at, updated_at 
					FROM books 
					WHERE title LIKE ? OR author LIKE ?
					ORDER BY created_at DESC 
					LIMIT ? OFFSET ?`

	rows, err := db.Query(searchQuery, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear,
			&book.Available, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	return books, total, nil
}
