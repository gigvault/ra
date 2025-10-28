package api

import (
	"encoding/json"
	"net/http"

	"github.com/gigvault/ra/internal/service"
	"github.com/gigvault/shared/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// HTTPHandler handles HTTP requests for the RA service
type HTTPHandler struct {
	service *service.RAService
	logger  *logger.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(service *service.RAService, logger *logger.Logger) *HTTPHandler {
	return &HTTPHandler{
		service: service,
		logger:  logger,
	}
}

// Routes returns the HTTP router
func (h *HTTPHandler) Routes() http.Handler {
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", h.Health).Methods("GET")
	r.HandleFunc("/ready", h.Ready).Methods("GET")

	// Enrollment operations
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/enrollments", h.CreateEnrollment).Methods("POST")
	api.HandleFunc("/enrollments", h.ListEnrollments).Methods("GET")
	api.HandleFunc("/enrollments/{id}", h.GetEnrollment).Methods("GET")
	api.HandleFunc("/enrollments/{id}/approve", h.ApproveEnrollment).Methods("POST")
	api.HandleFunc("/enrollments/{id}/reject", h.RejectEnrollment).Methods("POST")

	return h.loggingMiddleware(r)
}

// Health returns the health status
func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// Ready returns the readiness status
func (h *HTTPHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

// CreateEnrollment creates a new enrollment request
func (h *HTTPHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CommonName   string `json:"common_name"`
		Organization string `json:"organization"`
		Email        string `json:"email"`
		CSR          string `json:"csr"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	enrollment, err := h.service.CreateEnrollment(r.Context(), req.CommonName, req.Organization, req.Email, req.CSR)
	if err != nil {
		h.logger.Error("Failed to create enrollment", zap.Error(err))
		http.Error(w, "Failed to create enrollment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(enrollment)
}

// ListEnrollments lists all enrollments
func (h *HTTPHandler) ListEnrollments(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	enrollments, err := h.service.ListEnrollments(r.Context(), status)
	if err != nil {
		h.logger.Error("Failed to list enrollments", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollments)
}

// GetEnrollment retrieves an enrollment by ID
func (h *HTTPHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	enrollment, err := h.service.GetEnrollment(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get enrollment", zap.String("id", id), zap.Error(err))
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollment)
}

// ApproveEnrollment approves an enrollment
func (h *HTTPHandler) ApproveEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		ApprovedBy string `json:"approved_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.ApproveEnrollment(r.Context(), id, req.ApprovedBy); err != nil {
		h.logger.Error("Failed to approve enrollment", zap.String("id", id), zap.Error(err))
		http.Error(w, "Failed to approve enrollment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RejectEnrollment rejects an enrollment
func (h *HTTPHandler) RejectEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		RejectedBy string `json:"rejected_by"`
		Reason     string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.RejectEnrollment(r.Context(), id, req.RejectedBy, req.Reason); err != nil {
		h.logger.Error("Failed to reject enrollment", zap.String("id", id), zap.Error(err))
		http.Error(w, "Failed to reject enrollment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// loggingMiddleware logs HTTP requests
func (h *HTTPHandler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)
		next.ServeHTTP(w, r)
	})
}
