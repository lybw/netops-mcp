package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lybw/netkit/pkg/portscan"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPortscan(s *server.MCPServer) {
	tool := mcp.Tool{
		Name:        "scan_ports",
		Description: "Scan TCP ports on a target host. Returns open ports with service names and response latency.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"host": map[string]any{
					"type":        "string",
					"description": "Target host IP or hostname",
				},
				"ports": map[string]any{
					"type":        "string",
					"description": "Comma-separated ports or range (e.g., '22,80,443' or '1-1024'). Defaults to common ports if omitted.",
				},
				"timeout_ms": map[string]any{
					"type":        "integer",
					"description": "Connection timeout in milliseconds (default: 2000)",
				},
				"workers": map[string]any{
					"type":        "integer",
					"description": "Number of concurrent workers (default: 100)",
				},
			},
			Required: []string{"host"},
		},
	}

	s.AddTool(tool, handlePortscan)
}

func handlePortscan(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	host, err := request.RequireString("host")
	if err != nil {
		return errResult("missing required parameter: host"), nil
	}

	opts := portscan.DefaultOptions()

	if portsStr := request.GetString("ports", ""); portsStr != "" {
		ports, err := parsePorts(portsStr)
		if err != nil {
			return errResult(fmt.Sprintf("invalid ports: %v", err)), nil
		}
		opts.Ports = ports
	}

	if timeout := request.GetInt("timeout_ms", 0); timeout > 0 {
		opts.Timeout = time.Duration(timeout) * time.Millisecond
	}

	if workers := request.GetInt("workers", 0); workers > 0 {
		opts.Workers = workers
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	start := time.Now()
	results := portscan.Scan(ctx, host, opts)
	elapsed := time.Since(start)

	if len(results) == 0 {
		return textResult(fmt.Sprintf("No open ports found on %s (scanned %d ports in %s)",
			host, len(opts.Ports), elapsed.Round(time.Millisecond))), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Open ports on %s (%d found, scanned in %s):\n\n",
		host, len(results), elapsed.Round(time.Millisecond)))
	sb.WriteString(fmt.Sprintf("%-8s %-8s %-12s %s\n", "PORT", "STATE", "SERVICE", "LATENCY"))
	sb.WriteString(fmt.Sprintf("%-8s %-8s %-12s %s\n", "────", "─────", "───────", "───────"))

	for _, r := range results {
		sb.WriteString(fmt.Sprintf("%-8d %-8s %-12s %s\n",
			r.Port, "open", r.Service, r.Latency.Round(time.Millisecond)))
	}

	return textResult(sb.String()), nil
}

func parsePorts(spec string) ([]int, error) {
	var ports []int
	parts := strings.Split(spec, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			var start, end int
			if _, err := fmt.Sscanf(bounds[0], "%d", &start); err != nil {
				return nil, fmt.Errorf("invalid port: %s", bounds[0])
			}
			if _, err := fmt.Sscanf(bounds[1], "%d", &end); err != nil {
				return nil, fmt.Errorf("invalid port: %s", bounds[1])
			}
			ports = append(ports, portscan.PortRange(start, end)...)
		} else {
			var p int
			if _, err := fmt.Sscanf(part, "%d", &p); err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			ports = append(ports, p)
		}
	}
	return ports, nil
}
