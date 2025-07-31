# Create Local App

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
    - [Interactive Mode (First Run)](#interactive-mode-first-run)
    - [Auto Mode (Subsequent Runs)](#auto-mode-subsequent-runs)
    - [Example Workflow](#example-workflow)
  - [Contributing](#contributing)
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

Before using `create-local-app`, ensure you have the required versions of these tools. Click on the badges in the [Built With](#built-with) section for installation instructions.

**Required versions:**

- **Go**: v1.23.1 or higher
- **TrueBlocks**: v5.1.0 or higher  
- **Wails**: v2.10.1 or higher

**For contributors only:**
- **Yarn**: v1.22.22 or higher (for development tooling)

**Check your versions:**

```sh
go version
chifra version
wails version
# For contributors: yarn --version
```

### Installation

1. **Build from source**

   ```sh
   git clone https://github.com/TrueBlocks/create-local-app.git
   cd create-local-app
   
   # Build the Go binary
   go build -o bin/create-local-app main.go
   ```

2. **Install to your PATH** (required)

   ```sh
   # Copy to a system location
   sudo cp bin/create-local-app /usr/local/bin/
   
   # OR (temporarily) add the bin directory to your PATH
   export PATH="$PATH:$(pwd)/bin"
   ```

   > **⚠️ Important:** If you rebuild the binary after installing it to your PATH, you must also re-install it. The old version will continue to run until you copy the new binary over the old one.

3. **For contributors only** (optional development tooling)

   ```sh
   # Install linting tools (requires Node.js/Yarn)
   yarn install
   
   # Available development commands:
   yarn build       # Build the Go binary
   yarn lint        # Run Go and Markdown linters
   yarn clean       # Clean build artifacts
   ```

## Usage

When you first run `create-local-app`, it will interactively prompt you for project information and save your preferences for future use. Subsequent runs can use the `--auto` flag to skip prompts if nothing has changed.

### Interactive Mode (First Run)

The first time you run the command, it provides a guided setup experience:

```sh
create-local-app
```

You'll be prompted for:

- **Organization**: Your organization name (e.g., "TrueBlocks, LLC")
- **Project Name**: Name of your project (e.g., "my-awesome-app")
- **Go Import**: For importing Go packages - no spaces allowed (e.g., "github.com/TrueBlocks/my-awesome-app")
- **Domain**: The domain name of your home page (e.g., "trueblocks.io")

> **⚠️ Warning:** Running this command will overwrite existing files in the current directory. If you've made modifications to generated files, they will be lost. Make sure to commit your changes to version control before re-running.

### Auto Mode (Subsequent Runs)

Skip prompts and use previously saved configuration:

```sh
create-local-app --auto
```

*Note: This uses the configuration saved from your previous interactive run and also overwrites existing files.*

### Example Workflow

```sh
# 1. Create a new project interactively
mkdir my-new-app && cd my-new-app
create-local-app
    # > TrueBlocks, LLC
    # > my-new-app
    # > github.com/TrueBlocks/my-new-app
    # > https://trueblocks.io 

# 2. Set up the my-new-app project
yarn install

# 3. Testing
yarn test

# 4. Start the new app
yarn start
```

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

For detailed development information, see [DEVELOPING.md](DEVELOPING.md).

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
