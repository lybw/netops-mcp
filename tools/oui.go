package tools

import (
	"context"
	"fmt"

	"github.com/lybw/netkit/pkg/oui"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerOUI(s *server.MCPServer) {
	tool := mcp.Tool{
		Name:        "lookup_mac_vendor",
		Description: "Look up the manufacturer/vendor of a network device by its MAC address using the IEEE OUI database.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"mac": map[string]any{
					"type":        "string",
					"description": "MAC address in any format (e.g., '00:0C:29:AA:BB:CC', '00-0C-29-AA-BB-CC', or '000C29AABBCC')",
				},
			},
			Required: []string{"mac"},
		},
	}

	s.AddTool(tool, handleOUI)
}

func handleOUI(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mac, err := request.RequireString("mac")
	if err != nil {
		return errResult("missing required parameter: mac"), nil
	}

	vendor := oui.Lookup(mac)
	return textResult(fmt.Sprintf("MAC: %s\nVendor: %s", mac, vendor)), nil
}
