# Linear MCP Integration Setup

## Status
Linear MCP server configuration has been added to Claude Code.

## Setup Steps

### 1. Get Linear API Key

1. Go to https://linear.app/settings/api
2. Click "Personal API keys"
3. Click "Create key"
4. Give it a name (e.g., "MAAT Terminal MCP")
5. Copy the API key (starts with `lin_api_...`)

### 2. Add API Key to Environment

Add the API key to your shell profile:

```bash
# Add to ~/.zshrc (or ~/.bashrc)
echo 'export LINEAR_API_KEY="your_api_key_here"' >> ~/.zshrc

# Reload shell
source ~/.zshrc
```

**CRITICAL**: Never paste the API key in chat. Add it directly to the file manually.

### 3. Restart Claude Code

After setting the environment variable:
1. Quit Claude Code completely (Cmd+Q)
2. Relaunch Claude Code
3. The Linear MCP server will initialize on startup

### 4. Verify Connection

After restart, verify Linear MCP tools are available by attempting to:
- List Linear teams
- Search for issues
- Create test issue

## Configuration Location

MCP configuration: `~/Library/Application Support/Claude/claude_desktop_config.json`

Linear server added as:
```json
"linear": {
  "command": "npx",
  "args": ["-y", "@linear/mcp-server"],
  "env": {
    "LINEAR_API_KEY": "${LINEAR_API_KEY}"
  }
}
```

## MAAT Project Setup

Once Linear MCP is working, we'll:
1. Create or identify MAAT team in Linear
2. Create 6 initial issues for MAAT features
3. Configure MAAT TUI to read from Linear

## Next Steps

After you've set up the LINEAR_API_KEY:
1. Restart Claude Code
2. Return to this chat
3. I'll verify the connection and create the issues
