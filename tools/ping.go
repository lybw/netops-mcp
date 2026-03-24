package tools

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPing(s *server.MCPServer) {
	tool := mcp.Tool{
		Name:        "ping",
		Description: "Ping a host to check if it is reachable and measure round-trip latency.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"host": map[string]any{
					"type":        "string",
					"description": "Target host IP or hostname",
				},
				"count": map[string]any{
					"type":        "integer",
					"description": "Number of ping packets to send (default: 4)",
				},
			},
			Required: []string{"host"},
		},
	}

	s.AddTool(tool, handlePing)
}

func handlePing(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	host, err := request.RequireString("host")
	if err != nil {
		return errResult("missing required parameter: host"), nil
	}

	count := request.GetInt("count", 4)
	if count < 1 {
		count = 1
	}
	if count > 20 {
		count = 20
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	countStr := fmt.Sprintf("%d", count)

	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "ping", "-n", countStr, host)
	default:
		cmd = exec.CommandContext(ctx, "ping", "-c", countStr, host)
	}

	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	if err != nil {
		if output != "" {
			return textResult(fmt.Sprintf("Ping %s failed:\n\n%s", host, output)), nil
		}
		return errResult(fmt.Sprintf("ping failed: %v", err)), nil
	}

	return textResult(fmt.Sprintf("Ping %s:\n\n%s", host, output)), nil
}
