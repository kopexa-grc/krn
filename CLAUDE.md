# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

KRN (Kopexa Resource Names) is a Go package implementing resource identifiers following Google's Resource Name Design. The format is:

```
//kopexa.com/{service}/{collection}/{resource-id}[/{collection}/{resource-id}][@{version}]
```

## Build & Test Commands

```bash
# Run all tests with race detector
go test -race ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run a single test
go test -run TestParse ./...
go test -run TestParse/simple_KRN ./...

# Run benchmarks
go test -bench=. ./...

# Run linter (requires golangci-lint v2)
golangci-lint run

# Verify go.mod is tidy
go mod tidy
```

## Architecture

This is a single-package Go library with no external dependencies:

- `krn.go` - Core implementation: KRN struct, Parse/MustParse, Builder pattern, child creation, validation
- `krn_test.go` - Table-driven tests with 100% coverage requirement
- `example_test.go` - Runnable examples for godoc

### Key Types

- `KRN` - The main struct representing a Kopexa Resource Name
- `Segment` - A collection/resource-id pair
- `Builder` - Fluent API for constructing KRNs

### Services

Valid services are constants: `ServiceCatalog`, `ServiceISMS`, `ServiceOrg`, `ServiceAudit`, `ServicePolicy`

### Error Types

All errors are sentinel errors for `errors.Is()` compatibility: `ErrEmptyKRN`, `ErrInvalidKRN`, `ErrInvalidDomain`, `ErrInvalidService`, `ErrInvalidResourceID`, `ErrInvalidVersion`, `ErrResourceNotFound`

## Code Quality Requirements

- **100% test coverage** is enforced in CI
- All tests must pass with race detector enabled
- No golangci-lint warnings allowed
- Example tests must produce expected output (verified by `go test`)

## Release Process

Uses Google Release Please for automated releases. Push conventional commits to main:
- `feat:` - New features (minor version bump)
- `fix:` - Bug fixes (patch version bump)
- `chore:`, `docs:`, `test:` - No version bump

Release Please creates PRs that, when merged, trigger GoReleaser.
