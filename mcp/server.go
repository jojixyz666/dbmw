package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"dbmw/core/connection"
	"dbmw/core/explorer"
	"dbmw/core/query"
)

// MCPRequest represents an incoming JSON-RPC 2.0 request.
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents a JSON-RPC 2.0 response.
type MCPResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *MCPError `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Server handles the MCP protocol loop.
type Server struct {
	connSvc     *connection.Service
	schemaTools *SchemaTools
	queryTools  *QueryTools
}

func NewServer(connSvc *connection.Service, expSvc *explorer.Service, qSvc *query.Service) *Server {
	return &Server{
		connSvc:     connSvc,
		schemaTools: NewSchemaTools(expSvc),
		queryTools:  NewQueryTools(qSvc),
	}
}

// ServeStdio listens on stdin/stdout for MCP JSON-RPC protocol requests.
func (s *Server) ServeStdio(ctx context.Context, reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req MCPRequest
		if err := json.Unmarshal(line, &req); err != nil {
			resp := MCPResponse{
				JSONRPC: "2.0",
				ID:      nil,
				Error:   &MCPError{Code: -32700, Message: "Parse error"},
			}
			out, _ := json.Marshal(resp)
			fmt.Fprintln(writer, string(out))
			continue
		}

		resp := s.handleRequest(ctx, req)
		out, _ := json.Marshal(resp)
		fmt.Fprintln(writer, string(out))
	}
	return scanner.Err()
}

func (s *Server) handleRequest(ctx context.Context, req MCPRequest) MCPResponse {
	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		resp.Result = map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    "dbmw-mcp",
				"version": "0.0.1",
			},
		}

	case "tools/list":
		resp.Result = map[string]any{
			"tools": []map[string]any{
				{
					"name":        "read_schemas",
					"description": "List schemas in a chosen database (read-only)",
					"inputSchema": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"connectionId": map[string]any{"type": "string"},
							"database":     map[string]any{"type": "string"},
						},
						"required": []string{"connectionId"},
					},
				},
				{
					"name":        "read_tables",
					"description": "List tables and views in a schema (read-only)",
					"inputSchema": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"connectionId": map[string]any{"type": "string"},
							"schema":       map[string]any{"type": "string"},
						},
						"required": []string{"connectionId"},
					},
				},
				{
					"name":        "read_columns",
					"description": "Read column metadata and types for a table (read-only)",
					"inputSchema": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"connectionId": map[string]any{"type": "string"},
							"schema":       map[string]any{"type": "string"},
							"table":        map[string]any{"type": "string"},
						},
						"required": []string{"connectionId", "table"},
					},
				},
				{
					"name":        "execute_read_only_query",
					"description": "Execute a SELECT query against a connected database (prohibits all insert/update/delete/ddl writes)",
					"inputSchema": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"connectionId": map[string]any{"type": "string"},
							"query":        map[string]any{"type": "string"},
						},
						"required": []string{"connectionId", "query"},
					},
				},
			},
		}

	case "tools/call":
		var params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			resp.Error = &MCPError{Code: -32602, Message: "Invalid params"}
			return resp
		}

		res, err := s.callTool(ctx, params.Name, params.Arguments)
		if err != nil {
			resp.Result = map[string]any{
				"isError": true,
				"content": []map[string]any{
					{"type": "text", "text": err.Error()},
				},
			}
		} else {
			marshaled, _ := json.Marshal(res)
			resp.Result = map[string]any{
				"content": []map[string]any{
					{"type": "text", "text": string(marshaled)},
				},
			}
		}

	default:
		resp.Error = &MCPError{Code: -32601, Message: fmt.Sprintf("Method not found: %s", req.Method)}
	}

	return resp
}

func (s *Server) callTool(ctx context.Context, name string, rawArgs json.RawMessage) (any, error) {
	var args struct {
		ConnectionID string `json:"connectionId"`
		Database     string `json:"database"`
		Schema       string `json:"schema"`
		Table        string `json:"table"`
		Query        string `json:"query"`
	}
	if err := json.Unmarshal(rawArgs, &args); err != nil {
		return nil, fmt.Errorf("argument format error: %w", err)
	}

	if args.ConnectionID == "" {
		args.ConnectionID = s.connSvc.GetActiveID()
	}
	if args.ConnectionID == "" {
		return nil, fmt.Errorf("no active database connection selected")
	}

	conn, cfg, err := s.connSvc.GetConnector(ctx, args.ConnectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to load connector: %w", err)
	}

	switch name {
	case "read_schemas":
		return s.schemaTools.ReadSchemas(ctx, conn, args.Database)
	case "read_tables":
		return s.schemaTools.ReadTables(ctx, conn, args.Schema)
	case "read_columns":
		return s.schemaTools.ReadColumns(ctx, conn, args.Schema, args.Table)
	case "execute_read_only_query":
		return s.queryTools.ExecuteReadOnlyQuery(ctx, conn, args.ConnectionID, cfg.Name, args.Query)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func RunStdio(connSvc *connection.Service, expSvc *explorer.Service, qSvc *query.Service) error {
	srv := NewServer(connSvc, expSvc, qSvc)
	return srv.ServeStdio(context.Background(), os.Stdin, os.Stdout)
}
