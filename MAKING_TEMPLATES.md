# Making Templates for create-local-app

This document provides comprehensive information about creating, managing, and understanding templates in the `create-local-app` system.

## What Are Templates?

Templates are pre-configured project structures that serve as starting points for new applications. They contain:

- Directory structures
- Source code files with placeholder variables
- Configuration files (wails.json, package.json, etc.)
- Build scripts and tooling configuration

Templates use variable substitution with placeholders like `{{PROJECT_NAME}}`, `{{ORGANIZATION}}`, etc., which get replaced with actual values during project creation.

## Template Storage Locations

Templates are stored in a hierarchical structure under the user's configuration directory:

```
~/.create-local-app/
├── config.json                 # User configuration
├── VERSION                     # System templates version
└── templates/
    ├── system/                 # System templates (built-in)
    │   └── default/            # Default Wails project template
    └── contributed/            # User-created templates
        ├── my-custom-template/ # Your custom template
        └── another-template/   # Another custom template
```

### System Templates vs. Contributed Templates

**System Templates:**
- Built into the application binary as compressed archives
- Extracted automatically on first run or version updates
- Cannot be modified or removed by users
- Managed by core developers
- Examples: `default`

**Contributed Templates:**
- Created by users from existing projects
- Stored in `~/.create-local-app/templates/contributed/`
- Can be created, modified, and removed by users
- Persist across application updates
- Take precedence over system templates with the same name

## Template Resolution Order

When resolving a template name, the system searches in this order:

1. **Contributed templates** (`~/.create-local-app/templates/contributed/`)
2. **System templates** (`~/.create-local-app/templates/system/`)
3. **Environment variable path** (`TEMPLATE_SOURCE` as full path)
4. **Default system template** (if no template specified)

## Creating Templates

### From the Command Line

```bash
# Navigate to your customized project
cd my-awesome-project

# Create a template from the current project
create-local-app --create my-awesome-template

# The template is now available for use
create-local-app --template my-awesome-template
```

### What Gets Included

When creating a template, the system:

- **Includes:** All project files and directories
- **Excludes:** Build artifacts, node_modules, .git, .env files, etc.
- **Processes:** Replaces actual values with template variables
- **Preserves:** File permissions and directory structure

### Template Variables

The following variables are available for substitution:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{PROJECT_NAME}}` | Project name | `my-app` |
| `{{PROJECT_PROPER}}` | Capitalized project name | `My-app` |
| `{{ORGANIZATION}}` | Organization name | `TrueBlocks, LLC` |
| `{{ORG_NAME}}` | Organization name (first part) | `TrueBlocks` |
| `{{ORG_LOWER}}` | Lowercase organization | `trueblocks` |
| `{{GITHUB}}` | GitHub import path | `github.com/TrueBlocks/my-app` |
| `{{DOMAIN}}` | Domain name | `trueblocks.io` |
| `{{SLUG}}` | URL-friendly identifier | `trueblocks-my-app` |
| `{{CHIFRA}}` | TrueBlocks chifra import path | `github.com/TrueBlocks/trueblocks-chifra/v6` |

## Managing Templates

### Listing Available Templates

```bash
create-local-app --list
```

### Using a Template

```bash
# Use a specific template
create-local-app --template my-custom-template

# Template choice is saved for future --auto runs
create-local-app --auto
```

### Removing Templates

```bash
# Remove a contributed template (with confirmation)
create-local-app --remove my-custom-template
```

**Note:** Only contributed templates can be removed. System templates are protected.

## Template Best Practices

### For Template Creators

1. **Test thoroughly** - Create projects from your template to ensure it works
2. **Use meaningful variable names** - Stick to the standard template variables
3. **Include documentation** - Add README files explaining template-specific features
4. **Keep it minimal** - Don't include unnecessary dependencies or complexity
5. **Version control** - Consider versioning your template projects

### For Template Structure

```
my-template/
├── README.md                    # Template documentation
├── wails.json                   # Wails configuration (required)
├── go.mod                       # Go module definition
├── main.go                      # Main application entry
├── frontend/                    # Frontend code
│   ├── package.json
│   └── src/
└── build/                       # Build configuration
    └── appicon.png
```

### File Exclusions

The following are automatically excluded when creating templates:

- `.git/` - Git repository data
- `node_modules/` - Node.js dependencies  
- `dist/` - Distribution builds
- `build/` - Build artifacts (except appicon.png)
- `.env` - Environment files
- `.DS_Store` - macOS metadata
- `.create-local-app.json` - Project-local config

## For Core Developers Only

### System Template Management

System templates are embedded in the application binary and managed through the build process:

1. **Template Archives**: System templates are stored as `.tar.gz` files in `templates/system/`
2. **Embedding**: Templates are embedded using Go's `embed` directive
3. **Extraction**: Templates are extracted to user config on first run or version updates
4. **Versioning**: Template versions are tracked via the `VERSION` file

### Promoting Contributed Templates

The `scripts/promote-template.sh` script helps core developers promote contributed templates to system templates:

```bash
# Promote a contributed template to system template
./scripts/promote-template.sh my-contributed-template
```

**What the script does:**
- Copies the template from contributed to system directory
- Creates a compressed archive for embedding
- Prepares the template for inclusion in the next release

### Adding New System Templates

1. Create or obtain the template project
2. Test the template thoroughly
3. Use the promote script or manually create the archive
4. Update the build process to include the new template
5. Update documentation and examples

### Template Archive Format

System templates are stored as gzipped tar archives:

```bash
# Create a template archive
cd ~/.create-local-app/templates/system/
tar -czf my-template.tar.gz my-template/
```

The archives are then embedded in the binary and extracted at runtime.

## Troubleshooting

### Common Issues

**Template not found:**
- Check template name spelling
- Verify template exists in contributed or system directories
- Use `--list` to see available templates

**Permission errors:**
- Ensure write permissions to `~/.create-local-app/`
- Check if template directory is writable

**Template creation fails:**
- Ensure you're in a valid Wails project directory
- Check that `wails.json` exists
- Verify template name follows naming rules (alphanumeric + dashes)

### Debug Information

Enable verbose output by examining the console output, which shows:
- Template resolution path
- Config file locations
- Template directory being used

## Advanced Usage

### Environment Variable Override

For development and testing, you can override template resolution:

```bash
# Use a specific path
export TEMPLATE_SOURCE=/path/to/my/template
create-local-app

# Use a template name  
export TEMPLATE_SOURCE=my-template
create-local-app
```

### Multiple Template Workflows

Different projects can use different templates:

```bash
# Project A uses custom template
cd project-a
create-local-app --template react-advanced

# Project B uses different template  
cd project-b
create-local-app --template vue-basic

# Each project remembers its template choice
cd project-a && create-local-app --auto  # Uses react-advanced
cd project-b && create-local-app --auto  # Uses vue-basic
```

This enables teams to maintain multiple project templates while keeping each project's template preference isolated.
