# Developer Guide for create-local-app

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]

This document contains information specifically for developers working on the `create-local-app` tool itself. If you're looking to use the tool to generate TrueBlocks/Wails applications, please see the main [README.md](README.md).

This guide covers the internal workings of the scaffolding tool, including template variables, project roadmap, and development setup instructions.

## Development Setup

For contributors working on the `create-local-app` tool itself:

```sh
# Clone the repository
git clone https://github.com/TrueBlocks/create-local-app.git
cd create-local-app

# Install development dependencies
yarn install

# Available yarn scripts:
yarn build       # Build the Go binary
yarn lint        # Run Go and Markdown linters
yarn clean       # Clean build artifacts
yarn test        # Run Go tests

# Run the application in development mode
go run main.go

# Build for production
go build -o bin/create-local-app main.go
```

## Template Variables

The following variables are automatically replaced during project generation when the tool processes template files:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{PROJECT_NAME}}` | Project name (lowercase) | `my-awesome-app` |
| `{{PROJECT_PROPER}}` | Project name (title case) | `My-awesome-app` |
| `{{ORGANIZATION}}` | Full organization name | `TrueBlocks, LLC` |
| `{{ORG_NAME}}` | Organization name (first part) | `TrueBlocks` |
| `{{ORG_LOWER}}` | Organization name (lowercase) | `trueblocks` |
| `{{SLUG}}` | URL-friendly project identifier | `trueblocks-my-awesome-app` |
| `{{GITHUB}}` | Go import path | `github.com/TrueBlocks/my-awesome-app` |
| `{{DOMAIN}}` | Domain name | `trueblocks.io` |
| `{{CHIFRA}}` | TrueBlocks Chifra path | `github.com/TrueBlocks/trueblocks-core/src/apps/chifra` |
| `{{PUBLISHER_NAME}}` | Publisher name | `YourCompany` |
| `{{PUBLISHER_EMAIL}}` | Publisher email | `your_email@your_company.com` |

## Roadmap

- [x] Interactive project creation
- [x] Auto mode for rapid development
- [x] Configurable template variables
- [x] File exclusion system
- [ ] Custom template support (possible)
- [ ] Multiple template profiles (possible)
- [ ] Plugin system (unlikely)

See the [open issues](https://github.com/TrueBlocks/create-local-app/issues) for a full list of proposed features and known issues.

## Quality Assurance

- Go code is linted with `golangci-lint`
- Markdown files are linted with `markdownlint`
- All contributions should pass linting checks
- Run `yarn lint` before submitting pull requests

<!-- MARKDOWN LINKS & IMAGES -->
[contributors-shield]: https://img.shields.io/github/contributors/TrueBlocks/create-local-app.svg?style=for-the-badge
[contributors-url]: https://github.com/TrueBlocks/create-local-app/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/TrueBlocks/create-local-app.svg?style=for-the-badge
[forks-url]: https://github.com/TrueBlocks/create-local-app/network/members
[stars-shield]: https://img.shields.io/github/stars/TrueBlocks/create-local-app.svg?style=for-the-badge
[stars-url]: https://github.com/TrueBlocks/create-local-app/stargazers
