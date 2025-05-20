# budy - Your AI Terminal Assistant

Budy is a Go-based terminal assistant that helps you run commands, learns from your daily habits, and answers questions using AI.

## Features

- **Execute Terminal Commands**: Run commands directly in your terminal
- **Smart Command Suggestions**: Get suggestions based on your usage patterns
- **AI-Powered Help**: Ask questions and get smart responses
- **History Tracking**: Remembers your command history
- **Context-Aware**: Suggests commands based on current directory and time of day
- **Secure & Private**: All data stored locally

## Installation

### Prerequisites

- Go 1.18 or later
- OpenAI API key (set as environment variable `OPENAI_API_KEY`)

### Installation Options

#### From Source

1. Clone the repository
   ```
   git clone https://github.com/sosadtsia/budy.git
   cd budy
   ```

2. Build with Go
   ```
   go build -o budy ./cmd/budy
   ```

3. Install to a directory in your PATH
   ```
   # Using Homebrew location (recommended for M1/M2 Macs)
   cp budy /opt/homebrew/bin/

   # Using traditional location (may require sudo)
   sudo cp budy /usr/local/bin/

   # Using user bin directory
   mkdir -p ~/.local/bin
   cp budy ~/.local/bin/
   # Ensure ~/.local/bin is in your PATH
   ```

#### Using Go Install
```
go install github.com/sosadtsia/budy/cmd/budy@latest
```

## Usage

Start the assistant by running:
```
budy
```

### Commands and Questions

- Run any terminal command normally
  ```
  > ls -la
  ```

- Ask a question by prefixing with `?`
  ```
  > ? how do I find the largest files in a directory
  ```

- Exit the assistant
  ```
  > exit
  ```

## How It Works

Budy consists of several core components:

1. **Command Execution Engine**: Executes terminal commands through the system shell
2. **History Manager**: Tracks and stores command usage history
3. **Suggestion Engine**: Analyzes patterns in command usage to provide helpful suggestions
4. **AI Integration**: Connects to OpenAI API to answer questions

The assistant learns from your command usage patterns and provides increasingly relevant suggestions over time.

## Configuration

Budy stores its configuration and history in `~/.budy/` directory.

### OpenAI API Key

Set your OpenAI API key in one of the following ways:

1. As an environment variable:
   ```bash
   export OPENAI_API_KEY=your_api_key_here
   ```
   For permanent configuration, add this to your shell profile file (`.bashrc`, `.zshrc`, etc.).

2. Directly from within budy:
   ```
   > config set openai_key your_api_key_here
   ```
   This will store your API key in the configuration file.

## Development

### Project Structure

```
budy/
├── cmd/
│   └── budy/
│       └── main.go         # Entry point for the application
│
├── internal/
│   ├── ai/
│   │   └── openai.go       # OpenAI API integration
│   │
│   ├── shell/
│   │   ├── executor.go     # Command execution logic
│   │   └── history.go      # Command history management
│   │
│   ├── learning/
│   │   └── suggestions.go  # Command suggestion algorithms
│   │
│   └── storage/
│       └── file.go         # File-based storage using JSON
│
└── pkg/
    └── utils/
        ├── terminal.go     # Terminal utility functions
        └── paths.go        # Path management helpers
```

### Build and Test

```bash
# Build the application
task build

# Run tests
task test

# Run linter
task lint
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
