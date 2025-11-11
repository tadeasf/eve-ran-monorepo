# Swagger Documentation Guide

This guide explains how to maintain and regenerate Swagger/OpenAPI documentation for the EVE RAN backend API.

## Overview

This project uses [swaggo/swag](https://github.com/swaggo/swag) to automatically generate Swagger 2.0 documentation from Go annotations in the source code. The generated documentation is served at `/swagger/*` endpoint.

## Prerequisites

- Go 1.19 or later
- swag CLI tool

## Installation

Install the swag CLI tool using Go modules:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Verify installation:

```bash
swag --version
```

## Generating Documentation

### Basic Usage

From the backend directory, run:

```bash
cd /home/tadeas/eve-ran-monorepo/backend
swag init -g src/main.go -o docs
```

This will:
- Parse annotations in `src/main.go` and all imported packages
- Generate three files in the `docs/` directory:
  - `docs.go` - Go package with embedded documentation
  - `swagger.json` - JSON format specification
  - `swagger.yaml` - YAML format specification

### Advanced Options

#### Parse Internal Packages

To include structs and types from internal packages:

```bash
swag init -g src/main.go -o docs --parseInternal
```

#### Parse Dependencies

To include types from dependency packages:

```bash
swag init -g src/main.go -o docs --parseDependency
```

#### Parse Both Internal and Dependencies

For complete coverage:

```bash
swag init -g src/main.go -o docs --parseInternal --parseDependency
```

#### Quiet Mode

To reduce output verbosity:

```bash
swag init -g src/main.go -o docs --quiet
```

## Documentation Annotations

### General API Information

Add these annotations to your `main.go` file (already present):

```go
// @title           EVE RAN API
// @version         1.0
// @description     EVE Online kill tracking and analytics API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/tadeasf/eve-ran-monorepo
// @contact.email  your-email@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
```

### Route Annotations

For each API endpoint, add annotations above the handler function:

```go
// GetCompetitionSettings retrieves the current competition settings
// @Summary Get competition settings
// @Description Get the currently active competition settings (returns defaults if none exist)
// @Tags competition
// @Produce json
// @Success 200 {object} models.CompetitionSettings
// @Failure 500 {object} map[string]string
// @Router /competition/settings [get]
func GetCompetitionSettings(c *gin.Context) {
    // handler code
}
```

### Common Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@Summary` | Brief description | `@Summary Get user by ID` |
| `@Description` | Detailed description | `@Description Retrieves user information by ID` |
| `@Tags` | Group endpoints | `@Tags users,admin` |
| `@Accept` | Request content type | `@Accept json` |
| `@Produce` | Response content type | `@Produce json` |
| `@Param` | Parameter definition | `@Param id path int true "User ID"` |
| `@Success` | Success response | `@Success 200 {object} User` |
| `@Failure` | Error response | `@Failure 404 {object} ErrorResponse` |
| `@Router` | Route definition | `@Router /users/{id} [get]` |
| `@Security` | Security requirement | `@Security ApiKeyAuth` |

### Parameter Types

```go
// Path parameter
// @Param id path int true "User ID"

// Query parameter
// @Param search query string false "Search term"

// Body parameter
// @Param user body models.User true "User object"

// Header parameter
// @Param Authorization header string true "Bearer token"
```

## Current API Structure

The EVE RAN API includes the following endpoint groups:

### Core Endpoints
- `/` - Health check
- `/characters/*` - Character management
- `/regions/*` - Region data
- `/systems/*` - Solar system data
- `/kills/*` - Kill data and statistics

### Competition Endpoints (New)
- `/competition/settings` - GET/POST competition configuration
- `/competition/current` - GET current month standings
- `/competition/history` - GET all historical results
- `/competition/history/:month/:year` - GET specific month results
- `/competition/recent-winners` - GET last month's winners
- `/competition/ytd-winners` - GET year-to-date winners
- `/competition/save/:month/:year` - POST manually save monthly results

### Admin Endpoints
- `/admin/*` - Administrative functions

## Workflow for Adding New Endpoints

1. **Add route annotations** in your handler file:
   ```go
   // @Summary Your endpoint summary
   // @Description Detailed description
   // @Tags your-tag
   // @Accept json
   // @Produce json
   // @Param name path type required "description"
   // @Success 200 {object} YourModel
   // @Failure 400 {object} ErrorResponse
   // @Router /your/route [method]
   func YourHandler(c *gin.Context) {
       // implementation
   }
   ```

2. **Regenerate documentation**:
   ```bash
   swag init -g src/main.go -o docs --parseInternal
   ```

3. **Verify changes**:
   - Start the server: `go run src/main.go`
   - Visit: `http://localhost:8080/swagger/index.html`
   - Check that your new endpoint appears correctly

4. **Commit the generated files**:
   ```bash
   git add docs/
   git commit -m "docs: update Swagger documentation"
   ```

## Formatting Swagger Comments

Swag provides a formatting tool to ensure consistent comment style:

```bash
swag fmt
```

To format specific directories:

```bash
swag fmt -d src/ --exclude src/vendor
```

## Troubleshooting

### Documentation Not Updating

1. Ensure you're running `swag init` from the backend directory
2. Check that the `-g` flag points to the correct main file
3. Verify annotations are directly above function definitions
4. Make sure there are no syntax errors in annotations

### Missing Types

If your models aren't showing up:
- Add `--parseInternal` flag
- Ensure models have exported fields (capitalized)
- Check that JSON tags are present on struct fields

### Import Issues

The generated `docs.go` must be imported in your main file:

```go
import _ "github.com/tadeasf/eve-ran/docs"
```

## Docker Usage

If you prefer to use Docker:

```bash
docker run --rm -v $(pwd):/code ghcr.io/swaggo/swag:latest init -g src/main.go -o docs
```

## Accessing Swagger UI

Once the server is running:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **JSON spec**: `http://localhost:8080/swagger/doc.json`
- **YAML spec**: Files in `docs/swagger.yaml`

## Best Practices

1. **Always regenerate** after adding/modifying endpoints
2. **Use consistent tags** to group related endpoints
3. **Provide examples** in model definitions using `example` tags
4. **Document all parameters** including optional ones
5. **Include error responses** for all possible HTTP status codes
6. **Keep descriptions concise** but informative
7. **Commit generated files** so documentation stays in sync with code

## Additional Resources

- [Swaggo GitHub](https://github.com/swaggo/swag)
- [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format)
- [Gin-Swagger Integration](https://github.com/swaggo/gin-swagger)
- [OpenAPI/Swagger Specification](https://swagger.io/specification/)

## Quick Reference

```bash
# Install/update swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs (basic)
swag init -g src/main.go -o docs

# Generate docs (full)
swag init -g src/main.go -o docs --parseInternal --parseDependency

# Format swagger comments
swag fmt

# Show help
swag init -h
```
