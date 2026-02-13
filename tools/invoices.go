package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/fakturoid-mcp/fakturoid"
)

func registerInvoiceTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("fakturoid_invoice_list",
			mcp.WithDescription("List invoices (paginated, 40 per page)"),
			mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
			mcp.WithString("status", mcp.Description("Filter by status: open, sent, overdue, paid, cancelled")),
			mcp.WithNumber("subject_id", mcp.Description("Filter by subject (contact) ID")),
			mcp.WithString("since", mcp.Description("Filter invoices updated since date (ISO 8601)")),
		),
		invoiceListHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_invoice_detail",
			mcp.WithDescription("Get full detail of a specific invoice including lines"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Invoice ID")),
		),
		invoiceDetailHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_invoice_search",
			mcp.WithDescription("Search invoices by number, subject name, or note"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
			mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		),
		invoiceSearchHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_invoice_create",
			mcp.WithDescription("Create a new invoice"),
			mcp.WithNumber("subject_id", mcp.Required(), mcp.Description("Subject (contact) ID")),
			mcp.WithArray("lines", mcp.Required(), mcp.Description("Invoice lines (array of {name, quantity, unit_price, vat_rate, unit_name})")),
			mcp.WithString("currency", mcp.Description("Currency code (default: account currency)")),
			mcp.WithString("note", mcp.Description("Invoice note")),
			mcp.WithString("due_on", mcp.Description("Due date (YYYY-MM-DD)")),
			mcp.WithString("issued_on", mcp.Description("Issue date (YYYY-MM-DD)")),
		),
		invoiceCreateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_invoice_delete",
			mcp.WithDescription("Delete an invoice"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Invoice ID")),
		),
		invoiceDeleteHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_invoice_send",
			mcp.WithDescription("Send an invoice via email. Uses #link# in message to insert invoice link."),
			mcp.WithNumber("invoice_id", mcp.Required(), mcp.Description("Invoice ID")),
			mcp.WithString("email", mcp.Required(), mcp.Description("Recipient email address")),
			mcp.WithString("email_copy", mcp.Description("CC email address")),
			mcp.WithString("subject", mcp.Description("Email subject (Fakturoid uses default if empty)")),
			mcp.WithString("message", mcp.Description("Email body (use #link# for invoice link, Fakturoid uses default if empty)")),
		),
		invoiceSendHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_invoice_payments",
			mcp.WithDescription("List payments for an invoice"),
			mcp.WithNumber("invoice_id", mcp.Required(), mcp.Description("Invoice ID")),
		),
		invoicePaymentsHandler(r),
	)
}

func invoiceListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := intParam(req, "page", 1)
		params := url.Values{}
		if status := req.GetString("status", ""); status != "" {
			params.Set("status", status)
		}
		if subjectID := intParam(req, "subject_id", 0); subjectID != 0 {
			params.Set("subject_id", fmt.Sprintf("%d", subjectID))
		}
		if since := req.GetString("since", ""); since != "" {
			params.Set("since", since)
		}

		invoices, err := r.client.GetInvoices(page, params)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list invoices: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(invoices)), nil
	}
}

func invoiceDetailHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := intParam(req, "id", 0)
		if id == 0 {
			return mcp.NewToolResultError("id is required"), nil
		}

		invoice, err := r.client.GetInvoice(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get invoice: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(invoice)), nil
	}
}

func invoiceSearchHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := req.GetString("query", "")
		if query == "" {
			return mcp.NewToolResultError("query is required"), nil
		}
		page := intParam(req, "page", 1)

		invoices, err := r.client.SearchInvoices(query, page)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to search invoices: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(invoices)), nil
	}
}

func invoiceCreateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		subjectID := intParam(req, "subject_id", 0)
		if subjectID == 0 {
			return mcp.NewToolResultError("subject_id is required"), nil
		}

		args := req.GetArguments()
		linesRaw, ok := args["lines"]
		if !ok {
			return mcp.NewToolResultError("lines is required"), nil
		}

		lines, err := parseInvoiceLines(linesRaw)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid lines: %v", err)), nil
		}

		createReq := fakturoid.CreateInvoiceRequest{
			SubjectID: subjectID,
			Lines:     lines,
			Currency:  req.GetString("currency", ""),
			Note:      req.GetString("note", ""),
			DueOn:     req.GetString("due_on", ""),
			IssuedOn:  req.GetString("issued_on", ""),
		}

		invoice, err := r.client.CreateInvoice(createReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create invoice: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(invoice)), nil
	}
}

func invoiceDeleteHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := intParam(req, "id", 0)
		if id == 0 {
			return mcp.NewToolResultError("id is required"), nil
		}

		err := r.client.DeleteInvoice(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete invoice: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Invoice %d deleted", id)), nil
	}
}

func invoiceSendHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		invoiceID := intParam(req, "invoice_id", 0)
		if invoiceID == 0 {
			return mcp.NewToolResultError("invoice_id is required"), nil
		}
		email := req.GetString("email", "")
		if email == "" {
			return mcp.NewToolResultError("email is required"), nil
		}

		sendReq := fakturoid.SendInvoiceRequest{
			Email:     email,
			EmailCopy: req.GetString("email_copy", ""),
			Subject:   req.GetString("subject", ""),
			Message:   req.GetString("message", ""),
		}

		err := r.client.SendInvoice(invoiceID, sendReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to send invoice: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Invoice %d sent to %s", invoiceID, email)), nil
	}
}

func invoicePaymentsHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		invoiceID := intParam(req, "invoice_id", 0)
		if invoiceID == 0 {
			return mcp.NewToolResultError("invoice_id is required"), nil
		}

		payments, err := r.client.GetInvoicePayments(invoiceID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get payments: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(payments)), nil
	}
}
