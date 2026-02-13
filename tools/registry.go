package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/fakturoid-mcp/fakturoid"
)

// RegisterAll registers all Fakturoid MCP tools on the given server.
func RegisterAll(s *server.MCPServer, client *fakturoid.Client) {
	r := &registry{client: client}

	registerAccountTools(s, r)
	registerInvoiceTools(s, r)
	registerSubjectTools(s, r)
	registerExpenseTools(s, r)
}

type registry struct {
	client *fakturoid.Client
}
