package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerAccountTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("fakturoid_account_info",
			mcp.WithDescription("Get account information (company name, address, plan, currency, etc.)"),
		),
		accountInfoHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_events",
			mcp.WithDescription("Get recent account events (invoice created, paid, etc.)"),
		),
		eventsHandler(r),
	)
}

func accountInfoHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		account, err := r.client.GetAccount()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get account: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(account)), nil
	}
}

func eventsHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		events, err := r.client.GetEvents(1)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get events: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(events)), nil
	}
}
