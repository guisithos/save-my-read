CREATE TABLE books (
    id UUID PRIMARY KEY,
    google_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    authors TEXT[] NOT NULL,
    description TEXT,
    categories TEXT[] NOT NULL DEFAULT '{}',
    image_url TEXT,
    status VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for common queries
CREATE INDEX idx_books_user_id ON books(user_id);
CREATE INDEX idx_books_google_id ON books(google_id);
CREATE INDEX idx_books_status ON books(status); 