package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lybw/netkit/pkg/fingerprint"
	"github.com/lybw/netkit/pkg/portscan"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerFingerprint(s *server.MCPServer) {
	tool := mcp.Tool{
		Name:        "fingerprint_host",
		Description: "Detect the operating system and device type of a remote host by analyzing TCP/IP stack characteristics (TTL) and open port patterns.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"host": map[string]any{
					"type":        "string",
					"description": "Target host IP or hostname",
				},
				"scan_ports": map[string]any{
					"type":        "boolean",
					"description": "Also scan common ports to determine device type (default: true)",
				},
			},
			Required: []string{"host"},
		},
	}

	s.AddTool(tool, handleFingerprint)
}

func handleFingerprint(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	host, err := request.RequireString("host")
	if err != nil {
		return errResult("missing required parameter: host"), nil
	}

	scanPorts := request.GetBool("scan_ports", true)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Fingerprint results for %s:\n\n", host))

	// OS detection via TTL
	guess, err := fingerprint.DetectByTTL(host, 3*time.Second)
	if err != nil {
		sb.WriteString(fmt.Sprintf("OS Detection: failed (%v)\n", err))
	} else {
		sb.WriteString(fmt.Sprintf("OS:         %s\n", guess.OS))
		sb.WriteString(fmt.Sprintf("Confidence: %d%%\n", guess.Confidence))
		sb.WriteString(fmt.Sprintf("TTL:        %d\n", guess.TTL))
		sb.WriteString(fmt.Sprintf("Method:     %s\n", guess.Method))
	}

	// Device type detection via port scan
	if scanPorts {
		sb.WriteString("\n--- Port-based device classification ---\n\n")

		opts := portscan.Options{
			Timeout: 2 * time.Second,
			Workers: 50,
			Ports:   portscan.CommonPorts(),
		}

		ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		results := portscan.Scan(ctx, host, opts)

		var openPorts []int
		for _, r := range results {
			openPorts = append(openPorts, r.Port)
			sb.WriteString(fmt.Sprintf("  %-6d %s\n", r.Port, r.Service))
		}

		deviceType := fingerprint.DeviceType(openPorts)
		sb.WriteString(fmt.Sprintf("\nDevice Type: %s\n", deviceType))
		sb.WriteString(fmt.Sprintf("Open Ports:  %d\n", len(openPorts)))
	}

	return textResult(sb.String()), nil
}
