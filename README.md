# Create Local App

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![GPL 3.0 License][license-shield]][license-url]

**A powerful Go-based scaffolding tool for TrueBlocks/Wails desktop applications**

[**Explore the docs »**](https://github.com/TrueBlocks/create-local-app)

[Report Bug](https://github.com/TrueBlocks/create-local-app/issues/new?labels=bug&template=bug-report---.md) · [Request Feature](https://github.com/TrueBlocks/create-local-app/issues/new?labels=enhancement&template=feature-request---.md)

## Table of Contents

- [Create Local App](#create-local-app)
  - [Table of Contents](#table-of-contents)
  - [About The Project](#about-the-project)
    - [Built With](#built-with)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
    - [Interactive Mode (Default)](#interactive-mode-default)
    - [Auto Mode](#auto-mode)
    - [Example Workflow](#example-workflow)
  - [Operation Modes](#operation-modes)
    - [📝 Interactive Mode](#-interactive-mode)
    - [⚡ Auto Mode](#-auto-mode)
  - [Template Variables](#template-variables)
  - [Roadmap](#roadmap)
  - [Contributing](#contributing)
    - [Development Setup](#development-setup)
    - [Development Setup](#development-setup-1)
  - [License](#license)
  - [Contact](#contact)
  - [Acknowledgments](#acknowledgments)

## About The Project

`create-local-app` is a sophisticated scaffolding tool designed to streamline the creation of TrueBlocks/Wails desktop applications. With a single command, it generates a complete project structure with Go backend and TypeScript frontend, handling all the boilerplate code and configuration for you.

### Built With

This project is built using modern Go technologies and integrates seamlessly with the TrueBlocks ecosystem:

- [![Go][Go.dev]][Go-url]
- [![Wails][Wails.io]][Wails-url]
- [![TypeScript][TypeScript.org]][TypeScript-url]
- [![TrueBlocks][TrueBlocks.io]][TrueBlocks-url]

## Getting Started

To get a local copy up and running, follow these simple steps.

### Prerequisites

Before using `create-local-app`, ensure you have the required versions of these tools. Click on the badges in the [Built With](#built-with) section below for installation instructions.

**Required versions:**

- **Go**: v1.23.1 or higher
- **TrueBlocks**: v5.1.0 or higher
- **Wails**: v2.10.1 or higher
- **Yarn**: v1.22.22 or higher (for development tooling)

**Check your versions:**

```sh
go version
chifra version
wails version
yarn --version
```

### Installation

1. **Build from source**

   ```sh
   git clone https://github.com/TrueBlocks/create-local-app.git
   cd create-local-app
   
   # Using yarn script (recommended)
   yarn build
   
   # OR build directly with Go
   go build -o bin/create-local-app main.go
   ```

2. **Development setup**

   ```sh
   # Install dependencies (for linting and tooling)
   yarn install
   
   # Run linter
   yarn lint
   
   # Run markdown linter
   yarn lint:md
   
   # Clean build artifacts
   yarn clean
   ```

3. **Add to your PATH** (required)

   ```sh
   # Copy to a system location
   sudo cp bin/create-local-app /usr/local/bin/
   
   # OR add the bin directory to your PATH
   export PATH="$PATH:$(pwd)/bin"
   ```

4. **Set up template directory** (optional)

   ```sh
   # If you have a custom template directory
   export TEMPLATE_SOURCE=/path/to/your/template
   ```

## Usage

`create-local-app` supports three main operation modes, each designed for different use cases:

### Interactive Mode (Default)

The default mode provides a guided setup experience:

```sh
create-local-app
```

You'll be prompted for:

- **Organization**: Your organization name (e.g., "TrueBlocks, LLC")
- **Project Name**: Name of your project (e.g., "my-awesome-app")
- **GitHub**: Your GitHub username or organization
- **Domain**: Your domain name (e.g., "trueblocks.io")

### Auto Mode

Skip prompts and use previously saved configuration:

```sh
create-local-app --auto
```

*Note: Requires running interactive mode first to create the configuration file.*

### Example Workflow

```sh
# 1. Create a new project interactively
mkdir my-new-app && cd my-new-app
create-local-app

# 2. Set up the project
cd frontend && yarn install && cd ..

# 3. Start development
wails dev
```

## Operation Modes

### 📝 Interactive Mode

- **Purpose**: First-time setup or when you want to change configuration
- **Behavior**: Prompts for all required information
- **Config**: Saves settings to `.wails-template.json` for future use
- **Safety**: Warns before overwriting existing files

### ⚡ Auto Mode

- **Purpose**: Quick project creation with saved settings
- **Behavior**: Uses previously saved configuration without prompts
- **Requirements**: Must have run interactive mode first
- **Use Case**: Rapid prototyping or multiple similar projects

## Template Variables

The following variables are automatically replaced during project generation:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{PROJECT_NAME}}` | Project name (lowercase) | `my-awesome-app` |
| `{{PROJECT_PROPER}}` | Project name (title case) | `My-awesome-app` |
| `{{ORGANIZATION}}` | Full organization name | `TrueBlocks, LLC` |
| `{{ORG_NAME}}` | Organization name (first part) | `TrueBlocks` |
| `{{ORG_LOWER}}` | Organization name (lowercase) | `trueblocks` |
| `{{SLUG}}` | URL-friendly project identifier | `trueblocks-my-awesome-app` |
| `{{GITHUB}}` | GitHub username/organization | `TrueBlocks` |
| `{{DOMAIN}}` | Domain name | `trueblocks.io` |
| `{{CHIFRA}}` | TrueBlocks Chifra path | `github.com/TrueBlocks/trueblocks-core/src/apps/chifra` |
| `{{PUBLISHER_NAME}}` | Publisher name | `YourCompany` |
| `{{PUBLISHER_EMAIL}}` | Publisher email | `your_email@your_company.com` |

## Roadmap

- [x] Interactive project creation
- [x] Auto mode for rapid development
- [x] Configurable template variables
- [x] File exclusion system
- [ ] Custom template support
- [ ] Multiple template profiles
- [ ] GUI version
- [ ] Docker integration
- [ ] CI/CD templates
- [ ] Plugin system

See the [open issues](https://github.com/TrueBlocks/create-local-app/issues) for a full list of proposed features and known issues.

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Setup

For contributors, additional development tools are available:

```sh
# Install development dependencies
yarn install

# Available yarn scripts:
yarn build       # Build the Go binary
yarn lint        # Run Go linter (golangci-lint)
yarn lint:md     # Run Markdown linter
yarn clean       # Clean build artifacts
yarn test        # Run Go tests
```

**Quality Assurance:**
- Go code is linted with `golangci-lint`
- Markdown files are linted with `markdownlint`
- All contributions should pass linting checks

### Development Setup

```sh
# Clone the repository
git clone https://github.com/TrueBlocks/create-local-app.git
cd create-local-app

# Run the application
go run main.go

# Build for production
go build -o bin/create-local-app main.go
```

## License

Distributed under the GPL 3.0 License. See `LICENSE` for more information.

## Contact

TrueBlocks - [@trueblocks](https://twitter.com/trueblocks) - <info@trueblocks.io>

Project Link: [https://github.com/TrueBlocks/create-local-app](https://github.com/TrueBlocks/create-local-app)

## Acknowledgments

Resources and inspirations that made this project possible:

- [Wails Framework](https://wails.io/) - For creating amazing desktop apps with Go and web technologies
- [TrueBlocks](https://trueblocks.io/) - The decentralized data infrastructure
- [Go Programming Language](https://golang.org/) - For excellent standard library and tooling
- [Best README Template](https://github.com/othneildrew/Best-README-Template) - For this awesome README structure
- [Shields.io](https://shields.io/) - For the beautiful badges
- [Choose an Open Source License](https://choosealicense.com/) - For license guidance

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/TrueBlocks/create-local-app.svg?style=for-the-badge
[contributors-url]: https://github.com/TrueBlocks/create-local-app/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/TrueBlocks/create-local-app.svg?style=for-the-badge
[forks-url]: https://github.com/TrueBlocks/create-local-app/network/members
[stars-shield]: https://img.shields.io/github/stars/TrueBlocks/create-local-app.svg?style=for-the-badge
[stars-url]: https://github.com/TrueBlocks/create-local-app/stargazers
[issues-shield]: https://img.shields.io/github/issues/TrueBlocks/create-local-app.svg?style=for-the-badge
[issues-url]: https://github.com/TrueBlocks/create-local-app/issues
[license-shield]: https://img.shields.io/github/license/TrueBlocks/create-local-app.svg?style=for-the-badge
[license-url]: https://github.com/TrueBlocks/create-local-app/blob/main/LICENSE
[Go.dev]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://golang.org/
[Wails.io]: https://img.shields.io/badge/Wails-DF0000?style=for-the-badge&logo=wails&logoColor=white
[Wails-url]: https://wails.io/
[TypeScript.org]: https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white
[TypeScript-url]: https://www.typescriptlang.org/
[TrueBlocks.io]: https://img.shields.io/badge/TrueBlocks-4A90E2?style=for-the-badge&logo=ethereum&logoColor=white
[TrueBlocks-url]: https://trueblocks.io/
