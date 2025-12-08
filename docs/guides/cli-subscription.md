# CLI Subscription Setup Guide

AgentRunner uses your local CLI sessions to execute tasks. This avoids the need for API keys to be stored in the application and allows you to use your existing subscriptions.

## Supported Providers

- **Codex CLI**: `codex`
- **Claude Code**: `claude` / `claude-code`
- **Gemini CLI**: `gemini`
- **Cursor CLI**: `cursor`

## Setup Instructions

### 1. Codex CLI

1. Install Codex CLI.
2. Login to your account:
   ```bash
   codex login
   ```
   This should create a session file at `~/.codex/auth.json`.
3. AgentRunner will automatically mount this file into the sandbox container.

### 2. Claude Code

1. Install Claude Code (`npm install -g @anthropic-ai/claude-code`).
2. Login:
   ```bash
   claude login
   ```
3. Ensure the `claude` command is in your PATH.

### 3. Gemini CLI

1. Install Gemini CLI.
2. Login or setup credentials as per official documentation.

### 4. Cursor CLI

1. Ensure Cursor is installed and the CLI is available in your PATH.

## Configuration in Multiverse IDE

1. Open **Settings** -> **LLM**.
2. Select your desired provider from the list (e.g., `codex-cli`, `claude-code`).
3. Click "Test Connection" to verify that AgentRunner can access your local session.

## Troubleshooting

- **Session not found**: Ensure you have run the login command for the respective CLI.
- **Permission denied**: On macOS, you might need to grant Full Disk Access to Docker or the terminal running AgentRunner if it needs to read strict paths (though usually standard home paths are fine).
