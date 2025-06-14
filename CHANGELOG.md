# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Allow opening failed workflows
- Add command for generating act3's config

### Changed

- Changed CLI interface to use dual flags (-g/--global)
- Improve HTML layout
- Allow fetching results for multiple repos

## [v1.1.4] - Jan 07, 2025

### Changed

- Remove caching for GitHub client

## [v1.1.3] - Aug 20, 2024

### Changed

- Retain workflow order from config

## [v1.1.2] - Aug 20, 2024

### Fixed

- Show correct date for workflows in tabular and HTML output format

## [v1.1.0] - Aug 18, 2024

### Added

- Allow tabular output
- Show check state/conclusion for non-successful run

## [v1.0.0] - Aug 17, 2024

### Added

- Allow fetching results for a specific repository
- Better highlighting for workflow runs, based on run conclusion

### Changed

- act3 now relies on gh for authentication. If gh is not available,
  authenticating via an environment variable is still supported (via GH_TOKEN).
- act3 now uses the repo associated with the current directory by default.
    Results for preset workflows configured via a config file can still be
    fetched using a flag.

### Removed

- Ability to filter for messages with a specific context value

[unreleased]: https://github.com/dhth/act3/compare/v1.1.4...HEAD
[v1.1.4]: https://github.com/dhth/act3/compare/v1.1.3...v1.1.4
[v1.1.3]: https://github.com/dhth/act3/compare/v1.1.2...v1.1.3
[v1.1.2]: https://github.com/dhth/act3/compare/v1.1.1...v1.1.2
[v1.1.0]: https://github.com/dhth/act3/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/dhth/act3/compare/v0.4.0...v1.0.0
