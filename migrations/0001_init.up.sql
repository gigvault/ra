-- Enrollments table
CREATE TABLE IF NOT EXISTS enrollments (
    id SERIAL PRIMARY KEY,
    common_name VARCHAR(255) NOT NULL,
    organization VARCHAR(255),
    email VARCHAR(255),
    csr TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    approved_by VARCHAR(255),
    rejected_by VARCHAR(255),
    reject_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_enrollments_status ON enrollments(status);
CREATE INDEX idx_enrollments_common_name ON enrollments(common_name);
CREATE INDEX idx_enrollments_email ON enrollments(email);

