package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lybw/netkit/pkg/discovery"
	"github.com/lybw/netkit/pkg/oui"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerDiscover(s *server.MCPServer) {
	tool := mcp.Tool{
		Name:        "discover_devices",
		Description: "Discover devices on a local network using ARP scanning. Returns IP addresses, MAC addresses, and vendor information for all active devices on the specified subnet.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"cidr": map[string]any{
					"type":        "string",
					"description": "Target subnet in CIDR notation (e.g., 192.168.1.0/24)",
				},
			},
			Required: []string{"cidr"},
		},
	}

	s.AddTool(tool, handleDiscover)
}

func handleDiscover(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cidr, err := request.RequireString("cidr")
	if err != nil {
		return errResult("missing required parameter: cidr"), nil
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	devices, err := discovery.ARP(ctx, cidr)
	if err != nil {
		return errResult(fmt.Sprintf("discovery failed: %v", err)), nil
	}

	if len(devices) == 0 {
		return textResult("No devices found on " + cidr), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d device(s) on %s:\n\n", len(devices), cidr))
	sb.WriteString(fmt.Sprintf("%-16s %-18s %s\n", "IP", "MAC", "VENDOR"))
	sb.WriteString(fmt.Sprintf("%-16s %-18s %s\n", "──────────────", "─────────────────", "──────────────"))

	for _, d := range devices {
		vendor := oui.Lookup(d.MAC)
		sb.WriteString(fmt.Sprintf("%-16s %-18s %s\n", d.IP, d.MAC, vendor))
	}

	return textResult(sb.String()), nil
}
