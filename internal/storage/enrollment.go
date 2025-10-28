package storage

import (
	"context"
	"fmt"

	"github.com/gigvault/ra/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EnrollmentStorage handles enrollment database operations
type EnrollmentStorage struct {
	db *pgxpool.Pool
}

// NewEnrollmentStorage creates a new enrollment storage
func NewEnrollmentStorage(db *pgxpool.Pool) *EnrollmentStorage {
	return &EnrollmentStorage{db: db}
}

// Create creates a new enrollment
func (s *EnrollmentStorage) Create(ctx context.Context, enrollment *models.Enrollment) error {
	query := `
		INSERT INTO enrollments (common_name, organization, email, csr, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := s.db.QueryRow(ctx, query,
		enrollment.CommonName,
		enrollment.Organization,
		enrollment.Email,
		enrollment.CSR,
		enrollment.Status,
		enrollment.CreatedAt,
		enrollment.UpdatedAt,
	).Scan(&enrollment.ID)

	if err != nil {
		return fmt.Errorf("failed to create enrollment: %w", err)
	}

	return nil
}

// GetByID retrieves an enrollment by ID
func (s *EnrollmentStorage) GetByID(ctx context.Context, id string) (*models.Enrollment, error) {
	query := `
		SELECT id, common_name, organization, email, csr, status, 
		       approved_by, rejected_by, reject_reason, created_at, updated_at
		FROM enrollments
		WHERE id = $1
	`

	var enrollment models.Enrollment
	err := s.db.QueryRow(ctx, query, id).Scan(
		&enrollment.ID,
		&enrollment.CommonName,
		&enrollment.Organization,
		&enrollment.Email,
		&enrollment.CSR,
		&enrollment.Status,
		&enrollment.ApprovedBy,
		&enrollment.RejectedBy,
		&enrollment.RejectReason,
		&enrollment.CreatedAt,
		&enrollment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	return &enrollment, nil
}

// List retrieves all enrollments, optionally filtered by status
func (s *EnrollmentStorage) List(ctx context.Context, status string) ([]*models.Enrollment, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, common_name, organization, email, csr, status,
			       approved_by, rejected_by, reject_reason, created_at, updated_at
			FROM enrollments
			WHERE status = $1
			ORDER BY created_at DESC
			LIMIT 100
		`
		args = append(args, status)
	} else {
		query = `
			SELECT id, common_name, organization, email, csr, status,
			       approved_by, rejected_by, reject_reason, created_at, updated_at
			FROM enrollments
			ORDER BY created_at DESC
			LIMIT 100
		`
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list enrollments: %w", err)
	}
	defer rows.Close()

	var enrollments []*models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		err := rows.Scan(
			&enrollment.ID,
			&enrollment.CommonName,
			&enrollment.Organization,
			&enrollment.Email,
			&enrollment.CSR,
			&enrollment.Status,
			&enrollment.ApprovedBy,
			&enrollment.RejectedBy,
			&enrollment.RejectReason,
			&enrollment.CreatedAt,
			&enrollment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		enrollments = append(enrollments, &enrollment)
	}

	return enrollments, nil
}

// Update updates an enrollment
func (s *EnrollmentStorage) Update(ctx context.Context, enrollment *models.Enrollment) error {
	query := `
		UPDATE enrollments
		SET status = $1, approved_by = $2, rejected_by = $3, 
		    reject_reason = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := s.db.Exec(ctx, query,
		enrollment.Status,
		enrollment.ApprovedBy,
		enrollment.RejectedBy,
		enrollment.RejectReason,
		enrollment.UpdatedAt,
		enrollment.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update enrollment: %w", err)
	}

	return nil
}
