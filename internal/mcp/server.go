package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const protocolVersion = "2024-11-05"

type Server struct {
	name      string
	version   string
	tools     []Tool
	toolIndex map[string]Tool
	resources []Resource
	resIndex  map[string]Resource
}

type Resource struct {
	URI         string
	Name        string
	Description string
	MimeType    string
	Read        func() (string, error)
}

func NewServer(name, version string, tools []Tool, resources []Resource) *Server {
	ti := make(map[string]Tool, len(tools))
	for _, t := range tools {
		ti[t.Name] = t
	}
	ri := make(map[string]Resource, len(resources))
	for _, r := range resources {
		ri[r.URI] = r
	}
	return &Server{
		name: name, version: version,
		tools: tools, toolIndex: ti,
		resources: resources, resIndex: ri,
	}
}

func (s *Server) Serve(ctx context.Context) error {
	return s.serve(ctx, os.Stdin, os.Stdout)
}

func (s *Server) serve(ctx context.Context, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<20)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var msg jsonrpcMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			writeError(w, nil, codeParseError, "parse error")
			continue
		}

		if msg.Method == "" {
			continue
		}

		s.dispatch(w, &msg)
	}

	return scanner.Err()
}

func (s *Server) dispatch(w io.Writer, msg *jsonrpcMessage) {
	switch msg.Method {
	case "initialize":
		result := initializeResult{
			ProtocolVersion: protocolVersion,
			Capabilities: capabilities{
				Tools: &toolCapability{},
			},
			ServerInfo: serverInfo{Name: s.name, Version: s.version},
		}
		if len(s.resources) > 0 {
			result.Capabilities.Resources = &resourceCapability{}
		}
		writeResult(w, msg.ID, result)

	case "notifications/initialized":

	case "tools/list":
		defs := make([]toolDef, len(s.tools))
		for i, t := range s.tools {
			defs[i] = toolDef{
				Name:        t.Name,
				Description: t.Description,
				InputSchema: t.InputSchema,
			}
		}
		writeResult(w, msg.ID, toolListResult{Tools: defs})

	case "tools/call":
		s.handleToolCall(w, msg)

	case "resources/list":
		defs := make([]resourceDef, len(s.resources))
		for i, r := range s.resources {
			defs[i] = resourceDef{
				URI:         r.URI,
				Name:        r.Name,
				Description: r.Description,
				MimeType:    r.MimeType,
			}
		}
		writeResult(w, msg.ID, resourceListResult{Resources: defs})

	case "resources/read":
		s.handleResourceRead(w, msg)

	case "ping":
		writeResult(w, msg.ID, struct{}{})

	default:
		writeError(w, msg.ID, codeMethodNotFound, fmt.Sprintf("unknown method: %s", msg.Method))
	}
}

func (s *Server) handleToolCall(w io.Writer, msg *jsonrpcMessage) {
	var params toolCallParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		writeError(w, msg.ID, codeInvalidParams, "invalid params")
		return
	}

	tool, ok := s.toolIndex[params.Name]
	if !ok {
		writeError(w, msg.ID, codeInvalidParams, fmt.Sprintf("unknown tool: %s", params.Name))
		return
	}

	output, err := tool.Run(params.Arguments)
	if err != nil {
		writeResult(w, msg.ID, toolCallResult{
			Content: []content{{Type: "text", Text: fmt.Sprintf("error: %s", err)}},
			IsError: true,
		})
		return
	}

	writeResult(w, msg.ID, toolCallResult{
		Content: []content{{Type: "text", Text: output}},
	})
}

func (s *Server) handleResourceRead(w io.Writer, msg *jsonrpcMessage) {
	var params resourceReadParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		writeError(w, msg.ID, codeInvalidParams, "invalid params")
		return
	}

	res, ok := s.resIndex[params.URI]
	if !ok {
		writeError(w, msg.ID, codeInvalidParams, fmt.Sprintf("unknown resource: %s", params.URI))
		return
	}

	text, err := res.Read()
	if err != nil {
		writeError(w, msg.ID, codeInternalError, err.Error())
		return
	}

	writeResult(w, msg.ID, resourceReadResult{
		Contents: []resourceContent{{
			URI:      res.URI,
			MimeType: res.MimeType,
			Text:     text,
		}},
	})
}

// JSON-RPC types

type jsonrpcMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonrpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	codeParseError     = -32700
	codeMethodNotFound = -32601
	codeInvalidParams  = -32602
	codeInternalError  = -32603
)

// MCP protocol types

type initializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    capabilities `json:"capabilities"`
	ServerInfo      serverInfo   `json:"serverInfo"`
}

type capabilities struct {
	Tools     *toolCapability     `json:"tools,omitempty"`
	Resources *resourceCapability `json:"resources,omitempty"`
}

type toolCapability struct{}

type resourceCapability struct{}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type toolDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

type toolListResult struct {
	Tools []toolDef `json:"tools"`
}

type toolCallParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type toolCallResult struct {
	Content []content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

type content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type resourceDef struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

type resourceListResult struct {
	Resources []resourceDef `json:"resources"`
}

type resourceReadParams struct {
	URI string `json:"uri"`
}

type resourceReadResult struct {
	Contents []resourceContent `json:"contents"`
}

type resourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text"`
}

func writeResult(w io.Writer, id json.RawMessage, result any) {
	resp := jsonrpcResponse{JSONRPC: "2.0", ID: id, Result: result}
	data, _ := json.Marshal(resp)
	fmt.Fprintf(w, "%s\n", data)
}

func writeError(w io.Writer, id json.RawMessage, code int, message string) {
	resp := jsonrpcResponse{JSONRPC: "2.0", ID: id, Error: &jsonrpcError{Code: code, Message: message}}
	data, _ := json.Marshal(resp)
	fmt.Fprintf(w, "%s\n", data)
}
