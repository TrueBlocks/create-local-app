# {{ORG_NAME}} {{PROJECT_PROPER}}

An opinionated version of a desktop-based Wails app with Golang backend and React frontend using {{ORG_NAME}}'s SDK and Unchained Index.

## Features

- Desktop app built with Wails, React, TypeScript, Mantine, and {{ORG_NAME}}
- Golang backend - fast, concurrent, type safe
- ESLint & Prettier configured for frontend
- GolangCI-lint configured for backend
- VSCode automatic formatting on save (Go, TS, JS, YAML, TOML)

## Getting Started

### Prerequisites

- Golang >= 1.23.1
- Wails >= 2.10.1 and < 3.x
- Yarn (no npm)
- Node.js >= 18.x

### Installation

```bash
git clone https://github.com/{{ORG_NAME}}/{{SLUG}}.git
cd {{SLUG}}
yarn install
yarn test
```

### Running in Development

```bash
yarn dev
```

### Building for Production

```bash
yarn build
```

### Linting

```bash
yarn lint
```

## Project Structure

```[text]
.
├── app
├── frontend
│   ├── dist
│   │   └── assets
│   └── src
│       ├── components
│       ├── context
│       ├── contexts
│       ├── hooks
│       ├── layout
│       ├── utils
│       ├── views
│       └── wizards
│           └── hooks
└── pkg
    ├── markdown
    ├── msgs
    ├── preferences
    ├── project
    └── validation
```

## Contributing

We love contributors. Please see information about our workflow before proceeding.

- Fork this repository into your own repo.
- Create a branch: `git checkout -b <branch_name>`.
- Make changes to your local branch and commit them to your forked repo:  
  `git commit -m '<commit_message>'`
- Push back to the original branch:  
  `git push origin {{ORG_NAME}}/{{SLUG}}`
- Create the pull request.

## License

[LICENSE](./LICENSE)
