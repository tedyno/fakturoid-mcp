# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o fakturoid-mcp .   # build binary
go test ./...                  # run tests
./fakturoid-mcp                # run (communicates via stdio, MCP protocol)
```

## Architecture

MCP server exposing Fakturoid invoicing API as tools for AI assistants. Runs over stdio transport.

**Packages:**
- `main` — entry point, wires config → client → tools → MCP server
- `config` — loads credentials from env vars (`FAKTUROID_CLIENT_ID`, `FAKTUROID_CLIENT_SECRET`, `FAKTUROID_SLUG`) with fallback to `~/.config/fakturoid-mcp/config.json`
- `fakturoid` — HTTP client with OAuth2 client_credentials flow, auto token refresh (5min buffer), thread-safe via sync.Mutex. All API calls go through `client.do(method, endpoint, body, result)`
- `tools` — MCP tool definitions and handlers. Registry pattern: each resource file registers its tools via `register*Tools(server, registry)`

**Data flow:** MCP request → tool handler → fakturoid client → Fakturoid API v3 → JSON response → MCP result

**Key conventions:**
- Financial amounts use `json.Number` (Fakturoid API mixes string/number types for quantities, prices, VAT rates)
- Tool handlers return errors as `mcp.NewToolResultError()`, not Go errors
- Fakturoid API paginates at 40 items/page, pages start at 1
- API base URL: `https://app.fakturoid.cz/api/v3/accounts/{slug}/`
- OAuth token endpoint requires JSON body + Accept header (not form-encoded)

## Git

- Never add `Co-Authored-By` to commit messages.
