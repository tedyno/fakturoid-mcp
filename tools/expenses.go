package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerExpenseTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("fakturoid_expense_list",
			mcp.WithDescription("List expenses (paginated, 40 per page)"),
			mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
			mcp.WithString("status", mcp.Description("Filter by status: open, overdue, paid")),
			mcp.WithNumber("subject_id", mcp.Description("Filter by subject (contact) ID")),
		),
		expenseListHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_expense_detail",
			mcp.WithDescription("Get full detail of a specific expense"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Expense ID")),
		),
		expenseDetailHandler(r),
	)
}

func expenseListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := intParam(req, "page", 1)

		expenses, err := r.client.GetExpenses(page, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list expenses: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(expenses)), nil
	}
}

func expenseDetailHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := intParam(req, "id", 0)
		if id == 0 {
			return mcp.NewToolResultError("id is required"), nil
		}

		expense, err := r.client.GetExpense(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get expense: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(expense)), nil
	}
}
