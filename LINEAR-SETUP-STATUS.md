# Linear MCP Setup Status

## Current Status: Configuration Added, Awaiting API Key

### Completed Steps

1. ✅ Linear MCP server added to Claude Code configuration
   - Location: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Server: `@linear/mcp-server` via npx
   - Environment variable: `LINEAR_API_KEY`

2. ✅ Documentation created
   - Setup guide: `docs/LINEAR-MCP-SETUP.md`
   - Issue templates: `docs/INITIAL-ISSUES.md`

### Remaining Steps

1. ⏳ **You need to**: Get Linear API key
   - Go to: https://linear.app/settings/api
   - Create personal API key
   - Name it: "MAAT Terminal MCP"

2. ⏳ **You need to**: Add key to environment
   ```bash
   # Add to ~/.zshrc
   echo 'export LINEAR_API_KEY="lin_api_..."' >> ~/.zshrc
   source ~/.zshrc
   ```

3. ⏳ **You need to**: Restart Claude Code
   - Quit completely (Cmd+Q)
   - Relaunch
   - MCP server will initialize

4. ⏳ **After restart**: I'll create the issues
   - Verify Linear MCP connection
   - Create or identify MAAT team
   - Create 6 initial issues from templates

## Issues to Create

Once Linear MCP is active, we'll create:

1. **[High Priority]** Add GitHub API Integration
2. **[High Priority]** Add Linear API Integration
3. **[Medium Priority]** Improve Search/Filter Capabilities
4. **[Medium Priority]** Add Node Actions Menu
5. **[Low Priority]** Fix Graph View Scrolling
6. **[Low Priority]** Add README and Usage Guide

All issue details are in: `docs/INITIAL-ISSUES.md`

## Security Note

Following Commandment #1 of CLAUDE.md security rules:
- ❌ DO NOT paste your Linear API key in chat
- ✅ Add it directly to ~/.zshrc manually
- ✅ Use placeholder in examples: `your_api_key_here`

## Next Steps

1. Follow instructions in `docs/LINEAR-MCP-SETUP.md`
2. Restart Claude Code after setting LINEAR_API_KEY
3. Return to this chat
4. I'll verify connection and create issues
5. Configure MAAT to connect to Linear data source

## References

- Setup guide: `docs/LINEAR-MCP-SETUP.md`
- Issue templates: `docs/INITIAL-ISSUES.md`
- MCP config: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Linear API settings: https://linear.app/settings/api
