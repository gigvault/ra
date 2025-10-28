package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// CAClient handles gRPC communication with the CA service
type CAClient struct {
	conn   *grpc.ClientConn
	logger *zap.Logger
}

// CAClientConfig holds configuration for CA client
type CAClientConfig struct {
	Address    string
	TLSEnabled bool
	CACertPath string
	CertPath   string
	KeyPath    string
}

// NewCAClient creates a new CA gRPC client
func NewCAClient(cfg CAClientConfig, logger *zap.Logger) (*CAClient, error) {
	var opts []grpc.DialOption

	if cfg.TLSEnabled {
		// Load CA certificate
		caCert, err := os.ReadFile(cfg.CACertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to add CA certificate to pool")
		}

		// Load client certificate and key
		clientCert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}

		// Create TLS credentials
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs:      certPool,
			MinVersion:   tls.VersionTLS13,
		}

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// Add timeout
	opts = append(opts, grpc.WithTimeout(30*time.Second))

	// Connect to CA service
	conn, err := grpc.Dial(cfg.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to CA service: %w", err)
	}

	logger.Info("Connected to CA service",
		zap.String("address", cfg.Address),
		zap.Bool("tls_enabled", cfg.TLSEnabled),
	)

	return &CAClient{
		conn:   conn,
		logger: logger,
	}, nil
}

// SignCSR sends a CSR to CA for signing
func (c *CAClient) SignCSR(ctx context.Context, csrPEM string, validityDays int) (string, error) {
	// TODO: Define protobuf messages for CA service
	// For now, this is a stub that shows the interface
	c.logger.Info("Signing CSR via CA service",
		zap.Int("validity_days", validityDays),
	)

	/*
		Example implementation:
		
		client := capb.NewCAServiceClient(c.conn)
		req := &capb.SignCSRRequest{
			CsrPem:       csrPEM,
			ValidityDays: int32(validityDays),
		}
		
		resp, err := client.SignCSR(ctx, req)
		if err != nil {
			return "", fmt.Errorf("failed to sign CSR: %w", err)
		}
		
		return resp.CertificatePem, nil
	*/

	// Stub implementation
	return "", fmt.Errorf("CA gRPC client not fully implemented - protobuf definitions needed")
}

// GetCertificate retrieves a certificate by serial from CA
func (c *CAClient) GetCertificate(ctx context.Context, serial string) (string, error) {
	/*
		client := capb.NewCAServiceClient(c.conn)
		req := &capb.GetCertificateRequest{
			Serial: serial,
		}
		
		resp, err := client.GetCertificate(ctx, req)
		if err != nil {
			return "", fmt.Errorf("failed to get certificate: %w", err)
		}
		
		return resp.CertificatePem, nil
	*/

	return "", fmt.Errorf("CA gRPC client not fully implemented - protobuf definitions needed")
}

// RevokeCertificate requests certificate revocation from CA
func (c *CAClient) RevokeCertificate(ctx context.Context, serial string, reason string) error {
	/*
		client := capb.NewCAServiceClient(c.conn)
		req := &capb.RevokeCertificateRequest{
			Serial: serial,
			Reason: reason,
		}
		
		_, err := client.RevokeCertificate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to revoke certificate: %w", err)
		}
		
		return nil
	*/

	return fmt.Errorf("CA gRPC client not fully implemented - protobuf definitions needed")
}

// Close closes the gRPC connection
func (c *CAClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

