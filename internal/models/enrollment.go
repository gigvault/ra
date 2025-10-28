package models

import "time"

// Enrollment represents an enrollment request
type Enrollment struct {
	ID           string    `json:"id"`
	CommonName   string    `json:"common_name"`
	Organization string    `json:"organization"`
	Email        string    `json:"email"`
	CSR          string    `json:"csr"`
	Status       string    `json:"status"` // pending, approved, rejected
	ApprovedBy   *string   `json:"approved_by,omitempty"`
	RejectedBy   *string   `json:"rejected_by,omitempty"`
	RejectReason *string   `json:"reject_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
