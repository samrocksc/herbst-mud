# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- OpenAPI client generation capability for frontend type safety
- New `dev-frontend` Makefile target to automatically generate TypeScript client from backend API
- Integration with `@hey-api/openapi-ts` for generating frontend client code
- Automated client generation during development startup
- Database setup with ent ORM and PostgreSQL
- Cross-shaped rooms initialization with "The Hole" as the central starting room
- Database integration in both SSH server and web API server
- Users, characters, and rooms entities with relationships

### Changed
- Updated Makefile to include new development workflow
- Enhanced admin frontend with API client generation capabilities

## [0.1.0] - 2026-01-25

### Added
- Initial server scaffolding with Gin framework
- Basic health check and OpenAPI specification endpoints
- SSH server implementation
- Admin frontend with React and TanStack Router
- Development tooling with Vite and TypeScript

### Changed
- None

### Deprecated
- None

### Removed
- None

### Fixed
- None

### Security
- None