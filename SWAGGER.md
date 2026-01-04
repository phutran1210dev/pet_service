# Swagger Documentation Guide

## Generate Swagger Documentation

### Using Makefile (Recommended)

```bash
make swagger
```

This command will:

1. Check if `swag` is installed
2. Auto-install `swag` if not found
3. Generate Swagger docs from code annotations
4. Create/update files in `docs/` directory

### Manual Generation

If you prefer to run the command directly:

```bash
# Install swag (one time only)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs (from project root)
swag init -g main.go
# Or use full path if swag not in PATH
~/go/bin/swag init -g main.go
```

## For Team Members / CI/CD

### Option 1: Use Makefile

Just run `make swagger` - it will handle installation automatically.

### Option 2: Install swag globally

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Option 3: Add to CI/CD Pipeline

```yaml
# Example for GitHub Actions
- name: Generate Swagger docs
  run: make swagger

# Or
- name: Install swag
  run: go install github.com/swaggo/swag/cmd/swag@latest
  
- name: Generate docs
  run: swag init -g main.go --parseDependency --parseInternal
```

## Access Swagger UI

Once the API is running, access Swagger documentation at:

- <http://localhost:8001/swagger/index.html>

## API Response Structure

### Pagination Responses

Pagination endpoints return data in this structure:

```json
{
  "data": [...],
  "meta": {
    "total_items": 10,
    "total_pages": 2,
    "page": 1,
    "page_size": 5
  }
}
```

### Error Responses

```json
{
  "code": "ERROR_CODE",
  "message": "Error message",
  "details": []
}
```

### HTTP Status Codes

- 200 OK - Successful GET/UPDATE
- 201 Created - Successful POST creation
- 400 Bad Request - Validation errors
- 401 Unauthorized - Missing/invalid token
- 403 Forbidden - No permission
- 404 Not Found - Resource not found
- 409 Conflict - Resource already exists (e.g., duplicate email)
- 500 Internal Server Error - Server error

## Notes

- Swagger docs are generated from code annotations (comments)
- Files are generated in `docs/` directory
- Should regenerate after changing API endpoints or DTOs
- The generated files (`docs/`) should be committed to git
