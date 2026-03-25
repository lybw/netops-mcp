# netops-mcp

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/lybw/netops-mcp)](https://goreportcard.com/report/github.com/lybw/netops-mcp)

A Model Context Protocol (MCP) server for network operations. Gives AI agents the ability to discover devices, scan ports, fingerprint operating systems, and look up MAC vendors on local networks.

Built on [netkit](https://github.com/lybw/netkit) and [mcp-go](https://github.com/mark3labs/mcp-go).

## Tools

| Tool | Description |
|------|-------------|
| `discover_devices` | ARP scan a subnet to find active devices with IP, MAC, and vendor info |
| `scan_ports` | Concurrent TCP port scan with service detection |
| `fingerprint_host` | OS detection (TTL analysis) + device type classification (port patterns) |
| `lookup_mac_vendor` | MAC address to manufacturer lookup (IEEE OUI database) |
| `ping` | ICMP ping to check host reachability and latency |

## Install

```bash
go install github.com/lybw/netops-mcp@latest
```

Or build from source:

```bash
git clone https://github.com/lybw/netops-mcp.git
cd netops-mcp
go build -o netops-mcp .
```

## Configuration

### Claude Desktop

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "netops": {
      "command": "netops-mcp"
    }
  }
}
```

### Claude Code

Add to settings:

```json
{
  "mcpServers": {
    "netops": {
      "command": "netops-mcp"
    }
  }
}
```

## Usage Examples

Once connected, the AI agent can:

```
"Discover all devices on my local network 192.168.1.0/24"
"Scan ports 22, 80, 443 on 10.0.0.1"
"What OS is running on 192.168.1.100?"
"Look up the vendor for MAC address 00:0C:29:AA:BB:CC"
"Ping 8.8.8.8 to check connectivity"
```

## Tool Details

### discover_devices

```
Parameters:
  cidr (required) - Subnet in CIDR notation, e.g., "192.168.1.0/24"
```

Performs an ARP sweep: pings all IPs in the subnet to populate the ARP table, then reads it to extract IP/MAC pairs. Enriches results with vendor names from the OUI database.

### scan_ports

```
Parameters:
  host (required)  - Target IP or hostname
  ports            - Port spec: "22,80,443" or "1-1024" (default: common ports)
  timeout_ms       - Per-port timeout in ms (default: 2000)
  workers          - Concurrent scan workers (default: 100)
```

### fingerprint_host

```
Parameters:
  host (required)  - Target IP or hostname
  scan_ports       - Also scan ports for device type detection (default: true)
```

Combines TTL-based OS detection with port-pattern device classification (Server, Printer, Camera, Network Device, etc.).

### lookup_mac_vendor

```
Parameters:
  mac (required) - MAC address in any format
```

### ping

```
Parameters:
  host (required) - Target IP or hostname
  count           - Number of packets (default: 4, max: 20)
```

## Architecture

```
netops-mcp
├── main.go          # MCP server setup and stdio transport
└── tools/
    ├── register.go  # Tool registration
    ├── discover.go  # LAN device discovery
    ├── portscan.go  # Port scanning
    ├── fingerprint.go # OS + device type detection
    ├── oui.go       # MAC vendor lookup
    ├── ping.go      # ICMP ping
    └── helpers.go   # Shared response helpers
```

Dependencies:
- [netkit](https://github.com/lybw/netkit) — Network operations library (discovery, scanning, fingerprinting, OUI)
- [mcp-go](https://github.com/mark3labs/mcp-go) — MCP protocol implementation for Go

## License

MIT
