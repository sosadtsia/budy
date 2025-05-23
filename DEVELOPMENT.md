# Development Guide

This document provides guidelines and information for developers who want to contribute to the Budy project.

## Development Environment Setup

### Prerequisites

- Go 1.18 or later
- [Ollama](https://ollama.ai/) - For local AI model execution
- [Task](https://taskfile.dev/) - Task runner for development tasks
- [golangci-lint](https://golangci-lint.run/) - For code linting

### Setting Up Your Development Environment

1. Clone the repository:
   ```bash
   git clone https://github.com/sosadtsia/budy.git
   cd budy
   ```

2. Set up Git hooks (optional, but recommended):
   ```bash
   task hooks
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Build the project:
   ```bash
   task build
   ```

5. Run tests:
   ```bash
   task test
   ```

### Setting Up Ollama for Development

1. Install Ollama from [ollama.ai](https://ollama.ai/)

   **macOS Installation:**
   - Download the macOS app from [ollama.ai](https://ollama.ai/download)
   - Open the downloaded file and drag Ollama to your Applications folder
   - Launch Ollama from your Applications folder
   - You'll see the Ollama icon in your menu bar when it's running
   - Verify installation in terminal: `which ollama` (should show `/usr/local/bin/ollama`)

   **macOS Installation with Homebrew:**
   ```bash
   # Install Ollama
   brew install ollama

   # Start Ollama service
   brew services start ollama

   # Verify installation
   which ollama
   ollama --version
   ```

2. Start the Ollama server:
   ```bash
   # Ollama should start automatically after installation
   # If it's not running, you can start it manually
   ollama serve
   ```

3. Pull the models you want to use for development:
   ```bash
   # Pull the default model (REQUIRED)
   ollama pull llama3

   # Other models you might want to use (OPTIONAL)
   ollama pull mistral
   ollama pull gemma
   ```

4. Test Ollama connection:
   ```bash
   # List available models
   curl http://localhost:11434/api/tags
   ```

## Project Structure

The project follows a standard Go project layout:

```
budy/
├── cmd/                  # Application entry points
│   └── budy/             # Main application
│       └── main.go       # Application entry point
├── internal/             # Private application code
│   ├── ai/               # AI capabilities
│   ├── shell/            # Shell execution and history
│   ├── learning/         # Learning and suggestion algorithms
│   └── storage/          # Data storage
├── pkg/                  # Public libraries
│   └── utils/            # Utility functions
├── .githooks/            # Git hooks
├── .github/              # GitHub workflows and configs
└── docs/                 # Documentation
```

## Development Workflow

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes, following the coding standards.

3. Write appropriate tests for your changes.

4. Run tests locally:
   ```bash
   task test
   ```

5. Run linter to check code quality:
   ```bash
   task lint
   ```

6. Commit your changes using [Conventional Commits](https://www.conventionalcommits.org/) format:
   ```bash
   git commit -m "feat: add new feature"
   git commit -m "fix: fix bug in X"
   git commit -m "docs: update documentation"
   ```

7. Push your branch and create a Pull Request.

## Coding Standards

### Code Style

The project follows standard Go coding style guidelines:

- Use `gofmt` to format your code.
- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines.
- Run `golangci-lint` to ensure code quality.

### Testing

- Write unit tests for all new functionality.
- Aim for high test coverage, especially for core functionality.
- Use table-driven tests where appropriate.
- Mock external dependencies for testing.

## Dependencies

We strive to minimize external dependencies. Prefer using the Go standard library over third-party packages when possible.

### Adding Dependencies

If you need to add a new dependency:

1. Ensure it's absolutely necessary and there isn't a standard library alternative.
2. Check the library's license, maintenance status, and community adoption.
3. Add it using `go get`:
   ```bash
   go get github.com/example/package
   ```
4. Update dependencies:
   ```bash
   go mod tidy
   ```

## Building and Testing

### Building the Application

```bash
task build
```

This will create a `budy` executable in the project root.

### Running Tests

```bash
# Run all tests
task test

# Run tests with race detection
go test -race ./...

# Run tests for a specific package
go test ./internal/ai/...
```

### Running the Linter

```bash
task lint
```

## Release Process

The project follows [Semantic Versioning](https://semver.org/).

Releases are automatically created via GitHub Actions when commits following the Conventional Commits format are pushed to the main branch:

- `feat: ...` - Triggers a minor version bump
- `fix: ...` - Triggers a patch version bump
- `feat!: ...` or including `BREAKING CHANGE:` in commit message - Triggers a major version bump

## Development Tasks

Task automation is handled via [Task](https://taskfile.dev/). Common tasks:

```bash
task build      # Build the application
task test       # Run tests
task lint       # Run linter
task hooks      # Set up Git hooks
task clean      # Clean build artifacts
task run        # Run the application
```

## Troubleshooting

### Common Issues

- **Linter errors**: Run `task lint` to see detailed errors.
- **Test failures**: Run `go test -v ./...` for verbose output.
- **Build errors**: Check Go version and dependencies.
- **Ollama connection issues**:
  - Make sure Ollama is installed and running: `curl http://localhost:11434/api/tags`
  - If Ollama is not running, start it with `ollama serve` or `brew services start ollama`
  - Check that you've pulled the required models: `ollama list`
  - Verify Ollama URL in configuration is correct

### Getting Help

If you encounter issues during development:

1. Check existing GitHub issues first.
2. Feel free to create a new issue if your problem hasn't been reported.
3. Include detailed information about your environment and the problem.
