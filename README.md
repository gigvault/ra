# GigVault RA - Registration Authority

The Registration Authority service handles enrollment validation and approval workflows.

## Features

- Enrollment request management
- CSR validation
- Approval/rejection workflows
- Integration with CA service for certificate issuance
- RESTful API
- PostgreSQL-backed storage

## API Endpoints

### Health Checks
- `GET /health` - Health check
- `GET /ready` - Readiness check

### Enrollment Operations
- `POST /api/v1/enrollments` - Create enrollment request
- `GET /api/v1/enrollments` - List enrollments (optional ?status filter)
- `GET /api/v1/enrollments/{id}` - Get enrollment by ID
- `POST /api/v1/enrollments/{id}/approve` - Approve enrollment
- `POST /api/v1/enrollments/{id}/reject` - Reject enrollment

## Configuration

See `config/example.yaml` for configuration options.

## Development

```bash
# Build
make build

# Run tests
make test

# Run locally
make run-local

# Database migrations
make migrate
```

## Docker

```bash
# Build Docker image
make docker

# Run in Docker
docker run -p 8081:8081 gigvault/ra:local
```

## License

Copyright © 2025 GigVault

