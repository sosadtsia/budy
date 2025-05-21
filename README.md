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
- [Ollama](https://ollama.ai/) for local AI model execution (default)
- OpenAI API key (optional, only if you want to use OpenAI instead of Ollama)

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

### AI Providers

Budy supports two AI providers:

1. **Ollama** (default) - Uses [Ollama](https://ollama.ai/) for local AI model execution
2. **OpenAI** (optional) - Uses OpenAI's API services (requires API key)

To switch between providers:

```
> config set ai_provider openai    # Use OpenAI (requires API key)
> config set ai_provider ollama    # Use Ollama locally (default)
```

### Ollama Configuration

Ollama is the default AI provider and runs locally on your machine. To use it:

1. Install Ollama from [ollama.ai](https://ollama.ai/)

   **macOS Installation:**
   - Download the macOS app from [ollama.ai](https://ollama.ai/download)
   - Open the downloaded file and drag Ollama to your Applications folder
   - Launch Ollama from your Applications folder
   - You'll see the Ollama icon in your menu bar when it's running

   **macOS Installation with Homebrew:**
   ```bash
   # Install Ollama
   brew install ollama

   # Start Ollama service
   brew services start ollama
   ```

2. Run the Ollama server (it automatically starts after installation)

3. Pull the required model:
   ```bash
   # Pull the default model
   ollama pull llama3
   ```

4. Configure Ollama settings (optional):

   ```
   > config set ollama_url http://localhost:11434    # Default URL
   > config set ollama_model llama3                  # Default model
   ```

4. Available models depend on what you've pulled into Ollama. Some examples:
   ```
   > config set ollama_model llama3        # Use llama3
   > config set ollama_model mistral       # Use mistral
   > config set ollama_model gemma         # Use gemma
   ```

### OpenAI API Key (optional)

If you want to use OpenAI instead of Ollama, set your API key:

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

3. Then switch to OpenAI provider:
   ```
   > config set ai_provider openai
   ```

If you try to use OpenAI without setting an API key, budy will automatically fall back to using Ollama.

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
