package tools

import "github.com/mark3labs/mcp-go/server"

// RegisterAll registers all network operation tools with the MCP server.
func RegisterAll(s *server.MCPServer) {
	registerDiscover(s)
	registerPortscan(s)
	registerOUI(s)
	registerFingerprint(s)
	registerPing(s)
}
