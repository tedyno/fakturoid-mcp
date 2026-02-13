# fakturoid-mcp

MCP server for [Fakturoid](https://www.fakturoid.cz/) written in Go. Lightweight alternative to the [TypeScript version](https://github.com/cookielab/fakturoid-mcp) — single ~6 MB binary, ~12 MB RAM idle, no Node.js required.

## Setup

1. Get OAuth credentials in Fakturoid: Settings → User Account → API & Integrations
2. Create config file:

```bash
mkdir -p ~/.config/fakturoid-mcp
cat > ~/.config/fakturoid-mcp/config.json << 'EOF'
{
  "client_id": "your-client-id",
  "client_secret": "your-client-secret",
  "slug": "your-fakturoid-slug"
}
EOF
```

Or use environment variables: `FAKTUROID_CLIENT_ID`, `FAKTUROID_CLIENT_SECRET`, `FAKTUROID_SLUG`.

3. Build:

```bash
go build -o fakturoid-mcp .
```

## Claude Code integration

Add to `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "fakturoid": {
      "command": "/path/to/fakturoid-mcp"
    }
  }
}
```

## Docker

```json
{
  "mcpServers": {
    "fakturoid": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "FAKTUROID_CLIENT_ID=your-client-id",
        "-e", "FAKTUROID_CLIENT_SECRET=your-client-secret",
        "-e", "FAKTUROID_SLUG=your-slug",
        "tedyno/fakturoid-mcp:latest"
      ]
    }
  }
}
```

## Tools

| Tool | Description |
|------|-------------|
| `fakturoid_account_info` | Account details (company, address, plan, currency) |
| `fakturoid_events` | Recent account events |
| `fakturoid_invoice_list` | List invoices (filter by status, subject, date) |
| `fakturoid_invoice_detail` | Invoice detail with line items |
| `fakturoid_invoice_search` | Search by number, subject name, or note |
| `fakturoid_invoice_create` | Create invoice |
| `fakturoid_invoice_delete` | Delete invoice |
| `fakturoid_invoice_send` | Send invoice via email |
| `fakturoid_invoice_payments` | List payments for an invoice |
| `fakturoid_subject_list` | List contacts/clients |
| `fakturoid_subject_detail` | Contact detail |
| `fakturoid_subject_search` | Search contacts |
| `fakturoid_subject_create` | Create contact |
| `fakturoid_subject_update` | Update contact |
| `fakturoid_subject_delete` | Delete contact |
| `fakturoid_expense_list` | List expenses |
| `fakturoid_expense_detail` | Expense detail |
