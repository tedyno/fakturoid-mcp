package tools

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tedyno/fakturoid-mcp/fakturoid"
)

func toJSON(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("json error: %v", err)
	}
	return string(data)
}

func intParam(req mcp.CallToolRequest, name string, defaultVal int) int {
	return req.GetInt(name, defaultVal)
}

func parseInvoiceLines(raw any) ([]fakturoid.InvoiceLine, error) {
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal lines: %w", err)
	}
	var lines []fakturoid.InvoiceLine
	if err := json.Unmarshal(data, &lines); err != nil {
		return nil, fmt.Errorf("unmarshal lines: %w", err)
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("at least one line is required")
	}
	return lines, nil
}
