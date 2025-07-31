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

### Version Management

The application version is stored in the `VERSION` file at the project root and is embedded into the binary at build time using Go's `//go:embed` directive. To update the version:

1. Edit the `VERSION` file (contains only the version string, e.g., "0.2.0")
2. Rebuild the binary - Go will automatically detect the change and rebuild
3. The new version will be available via `--version` command

**Note**: Changes to embedded files (like `VERSION`) trigger a rebuild, so you don't need to use `go clean` when updating the version.

## Creating Custom Templates

The `create-local-app` tool supports creating custom templates from existing projects using the `--reverse` mode. This allows you to capture your project structure and configurations as reusable templates.

### Creating a Template from Your Project

1. **Navigate to your completed project directory**:
   ```sh
   cd my-awesome-project
   ```

2. **Create a template using reverse mode**:
   ```sh
   create-local-app --reverse my-custom-template
   ```

3. **Template storage location**:
   Your custom template will be stored in:
   ```
   ~/.create-local-app/templates/contributed/my-custom-template/
   ```

### How Templates Are Processed

When creating a template with `--reverse`, the tool:

- **Excludes build artifacts**: Automatically skips `.git/`, `node_modules/`, `dist/`, and other generated files
- **Reverse template variables**: Converts your actual values back to template variables (e.g., "my-awesome-project" becomes `{{PROJECT_NAME}}`)
- **Preserves structure**: Maintains your directory structure and file permissions
- **Stores contributed templates**: Places custom templates in the `contributed/` folder to distinguish them from system templates

### Using Custom Templates (Future Feature)

> **📝 Note**: This feature is planned but not yet implemented.

In future versions, you'll be able to use your custom templates with:

```sh
# Use a custom template instead of the default
create-local-app --template my-custom-template

# Combine with other options
create-local-app --template my-custom-template --auto --force
```

### Template Directory Structure

```
~/.create-local-app/templates/
├── system/                     # Built-in templates (managed by the tool)
│   └── default/               # Default Wails project template
└── contributed/               # Your custom templates
    ├── my-custom-template/    # Created with --reverse
    ├── react-template/        # Another custom template
    └── minimal-template/      # Minimal project template
```

### Best Practices for Template Creation

1. **Clean your project first**: Remove temporary files, logs, and personal configurations
2. **Test your project**: Ensure everything works before creating the template
3. **Use meaningful names**: Choose descriptive template names (e.g., `dashboard-app`, `cli-tool`)
4. **Document your templates**: Consider adding a README in your template directory
5. **Version control**: Your templates are stored locally - consider backing them up

### Template Variable Replacement

When creating templates, these values are automatically converted to template variables:

| Your Value | Becomes Template Variable |
|------------|---------------------------|
| `my-awesome-project` | `{{PROJECT_NAME}}` |
| `My-awesome-project` | `{{PROJECT_PROPER}}` |
| `TrueBlocks, LLC` | `{{ORGANIZATION}}` |
| `github.com/TrueBlocks/my-awesome-project` | `{{GITHUB}}` |
| `trueblocks.io` | `{{DOMAIN}}` |

This ensures that when someone uses your template, these values will be replaced with their own project-specific information.

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
