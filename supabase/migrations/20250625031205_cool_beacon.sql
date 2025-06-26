-- Initialize the library database schema
-- This script creates the initial database structure for the library management system

USE db;

-- Create books table with proper indexing and constraints
CREATE TABLE IF NOT EXISTS books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    published_year INT NOT NULL CHECK (published_year >= 1000 AND published_year <= 2100),
    available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes for better query performance
    INDEX idx_title (title),
    INDEX idx_author (author),
    INDEX idx_published_year (published_year),
    INDEX idx_available (available),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert sample data for testing
INSERT INTO books (title, author, published_year, available) VALUES
('The Go Programming Language', 'Alan Donovan, Brian Kernighan', 2015, true),
('Clean Code', 'Robert C. Martin', 2008, true),
('Design Patterns', 'Gang of Four', 1994, true),
('The Pragmatic Programmer', 'Andy Hunt, Dave Thomas', 1999, false),
('Effective Go', 'The Go Team', 2020, true),
('Database Design for Mere Mortals', 'Michael J. Hernandez', 2013, true),
('RESTful Web APIs', 'Leonard Richardson, Mike Amundsen', 2013, true),
('Docker Deep Dive', 'Nigel Poulton', 2020, false),
('Microservices Patterns', 'Chris Richardson', 2018, true),
('Building Microservices', 'Sam Newman', 2015, true);