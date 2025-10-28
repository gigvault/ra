package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gigvault/ra/internal/models"
	"github.com/gigvault/ra/internal/storage"
	"github.com/gigvault/shared/pkg/crypto"
	"github.com/gigvault/shared/pkg/logger"
	"go.uber.org/zap"
)

// RAService handles registration authority operations
type RAService struct {
	storage *storage.EnrollmentStorage
	logger  *logger.Logger
}

// NewRAService creates a new RA service
func NewRAService(storage *storage.EnrollmentStorage, logger *logger.Logger) *RAService {
	return &RAService{
		storage: storage,
		logger:  logger,
	}
}

// CreateEnrollment creates a new enrollment request
func (s *RAService) CreateEnrollment(ctx context.Context, cn, org, email, csrPEM string) (*models.Enrollment, error) {
	// Validate CSR
	csr, err := crypto.ParseCSR([]byte(csrPEM))
	if err != nil {
		return nil, fmt.Errorf("invalid CSR: %w", err)
	}

	s.logger.Info("Creating enrollment",
		zap.String("cn", cn),
		zap.String("org", org),
		zap.String("csr_subject", csr.Subject.CommonName),
	)

	// Validate that CSR subject matches provided CN
	if csr.Subject.CommonName != cn {
		return nil, fmt.Errorf("CSR common name does not match provided common name")
	}

	enrollment := &models.Enrollment{
		CommonName:   cn,
		Organization: org,
		Email:        email,
		CSR:          csrPEM,
		Status:       "pending",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.storage.Create(ctx, enrollment); err != nil {
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}

	return enrollment, nil
}

// ListEnrollments lists enrollments, optionally filtered by status
func (s *RAService) ListEnrollments(ctx context.Context, status string) ([]*models.Enrollment, error) {
	return s.storage.List(ctx, status)
}

// GetEnrollment retrieves an enrollment by ID
func (s *RAService) GetEnrollment(ctx context.Context, id string) (*models.Enrollment, error) {
	return s.storage.GetByID(ctx, id)
}

// ApproveEnrollment approves an enrollment request
func (s *RAService) ApproveEnrollment(ctx context.Context, id, approvedBy string) error {
	enrollment, err := s.storage.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}

	if enrollment.Status != "pending" {
		return fmt.Errorf("enrollment is not in pending status")
	}

	s.logger.Info("Approving enrollment",
		zap.String("id", id),
		zap.String("approved_by", approvedBy),
	)

	enrollment.Status = "approved"
	enrollment.ApprovedBy = &approvedBy
	enrollment.UpdatedAt = time.Now()

	if err := s.storage.Update(ctx, enrollment); err != nil {
		return fmt.Errorf("failed to update enrollment: %w", err)
	}

	// Forward approved CSR to CA service for signing (async)
	// This would typically be done by a background worker that:
	// 1. Picks up approved enrollments
	// 2. Calls CA service via gRPC
	// 3. Updates enrollment with certificate
	// 
	// For now, mark as approved and let a worker handle the actual signing
	// See: ra/internal/grpc/ca_client.go for CA integration
	s.logger.Info("Enrollment queued for CA signing")

	return nil
}

// RejectEnrollment rejects an enrollment request
func (s *RAService) RejectEnrollment(ctx context.Context, id, rejectedBy, reason string) error {
	enrollment, err := s.storage.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}

	if enrollment.Status != "pending" {
		return fmt.Errorf("enrollment is not in pending status")
	}

	s.logger.Info("Rejecting enrollment",
		zap.String("id", id),
		zap.String("rejected_by", rejectedBy),
		zap.String("reason", reason),
	)

	enrollment.Status = "rejected"
	enrollment.RejectedBy = &rejectedBy
	enrollment.RejectReason = &reason
	enrollment.UpdatedAt = time.Now()

	if err := s.storage.Update(ctx, enrollment); err != nil {
		return fmt.Errorf("failed to update enrollment: %w", err)
	}

	return nil
}
