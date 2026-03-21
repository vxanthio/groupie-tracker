# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Live search via `GET /api/search?q=` with debounced client-side fetch
- Artist detail page with members, locations, and concert dates
- 404 and 500 error pages

## [0.1.0] - 2026-03-21

### Added
- Project bootstrap: Go module, folder structure, HTTP server
- API client fetching artists, locations, dates, and relations
- In-memory data store with `Store` interface and `MockStore` for testing
- `GET /` handler rendering the artist list
- `GET /artist/{id}` handler rendering artist detail
- Base layout, home template, artist template, error templates
- Dark-theme CSS design system
- Static file server under `/static/`

[Unreleased]: https://github.com/your-team/groupie-tracker/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/your-team/groupie-tracker/releases/tag/v0.1.0
