package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/fakturoid-mcp/fakturoid"
)

func registerSubjectTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("fakturoid_subject_list",
			mcp.WithDescription("List subjects (contacts/clients), paginated"),
			mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		),
		subjectListHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_subject_detail",
			mcp.WithDescription("Get full detail of a specific subject (contact)"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Subject ID")),
		),
		subjectDetailHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_subject_search",
			mcp.WithDescription("Search subjects by name, email, registration number, etc."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
			mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		),
		subjectSearchHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_subject_create",
			mcp.WithDescription("Create a new subject (contact/client)"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Company or person name")),
			mcp.WithString("street", mcp.Description("Street address")),
			mcp.WithString("city", mcp.Description("City")),
			mcp.WithString("zip", mcp.Description("ZIP/postal code")),
			mcp.WithString("country", mcp.Description("Country code (e.g. CZ, SK)")),
			mcp.WithString("registration_no", mcp.Description("Registration number (IČO)")),
			mcp.WithString("vat_no", mcp.Description("VAT number (DIČ)")),
			mcp.WithString("email", mcp.Description("Email")),
			mcp.WithString("phone", mcp.Description("Phone")),
		),
		subjectCreateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_subject_update",
			mcp.WithDescription("Update an existing subject"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Subject ID")),
			mcp.WithString("name", mcp.Description("Company or person name")),
			mcp.WithString("street", mcp.Description("Street address")),
			mcp.WithString("city", mcp.Description("City")),
			mcp.WithString("zip", mcp.Description("ZIP/postal code")),
			mcp.WithString("country", mcp.Description("Country code")),
			mcp.WithString("registration_no", mcp.Description("Registration number (IČO)")),
			mcp.WithString("vat_no", mcp.Description("VAT number (DIČ)")),
			mcp.WithString("email", mcp.Description("Email")),
			mcp.WithString("phone", mcp.Description("Phone")),
		),
		subjectUpdateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("fakturoid_subject_delete",
			mcp.WithDescription("Delete a subject"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Subject ID")),
		),
		subjectDeleteHandler(r),
	)
}

func subjectListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := intParam(req, "page", 1)

		subjects, err := r.client.GetSubjects(page)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list subjects: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(subjects)), nil
	}
}

func subjectDetailHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := intParam(req, "id", 0)
		if id == 0 {
			return mcp.NewToolResultError("id is required"), nil
		}

		subject, err := r.client.GetSubject(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get subject: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(subject)), nil
	}
}

func subjectSearchHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := req.GetString("query", "")
		if query == "" {
			return mcp.NewToolResultError("query is required"), nil
		}
		page := intParam(req, "page", 1)

		subjects, err := r.client.SearchSubjects(query, page)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to search subjects: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(subjects)), nil
	}
}

func subjectCreateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := req.GetString("name", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		createReq := fakturoid.CreateSubjectRequest{
			Name:           name,
			Street:         req.GetString("street", ""),
			City:           req.GetString("city", ""),
			Zip:            req.GetString("zip", ""),
			Country:        req.GetString("country", ""),
			RegistrationNo: req.GetString("registration_no", ""),
			VATNo:          req.GetString("vat_no", ""),
			Email:          req.GetString("email", ""),
			Phone:          req.GetString("phone", ""),
		}

		subject, err := r.client.CreateSubject(createReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create subject: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(subject)), nil
	}
}

func subjectUpdateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := intParam(req, "id", 0)
		if id == 0 {
			return mcp.NewToolResultError("id is required"), nil
		}

		updateReq := fakturoid.UpdateSubjectRequest{
			Name:           req.GetString("name", ""),
			Street:         req.GetString("street", ""),
			City:           req.GetString("city", ""),
			Zip:            req.GetString("zip", ""),
			Country:        req.GetString("country", ""),
			RegistrationNo: req.GetString("registration_no", ""),
			VATNo:          req.GetString("vat_no", ""),
			Email:          req.GetString("email", ""),
			Phone:          req.GetString("phone", ""),
		}

		subject, err := r.client.UpdateSubject(id, updateReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update subject: %v", err)), nil
		}
		return mcp.NewToolResultText(toJSON(subject)), nil
	}
}

func subjectDeleteHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := intParam(req, "id", 0)
		if id == 0 {
			return mcp.NewToolResultError("id is required"), nil
		}

		err := r.client.DeleteSubject(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete subject: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Subject %d deleted", id)), nil
	}
}
