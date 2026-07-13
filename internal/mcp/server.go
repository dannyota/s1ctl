package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

const protocolVersion = "2024-11-05"

type Server struct {
	name      string
	version   string
	tools     []Tool
	toolIndex map[string]Tool
	resources []Resource
	resIndex  map[string]Resource

	// Dynamic tool loading (listChanged).
	w            io.Writer
	root         *cobra.Command
	allToolIndex map[string]Tool
	metaTools    []Tool
	focused      map[string][]Tool
	toolsVersion uint64

	mu sync.Mutex // guards w, focused, tools, toolIndex, toolsVersion
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

func NewDynamicServer(name, version string, root *cobra.Command, resources []Resource) *Server {
	ri := make(map[string]Resource, len(resources))
	for _, r := range resources {
		ri[r.URI] = r
	}
	allTools := ToolsFromCobra(root)
	ati := make(map[string]Tool, len(allTools))
	for _, t := range allTools {
		ati[t.Name] = t
	}
	s := &Server{
		name:         name,
		version:      version,
		resources:    resources,
		resIndex:     ri,
		root:         root,
		allToolIndex: ati,
		focused:      make(map[string][]Tool),
	}
	s.metaTools = s.buildMetaTools()
	s.rebuildToolList()
	return s
}

// rebuildToolList must be called with s.mu held.
func (s *Server) rebuildToolList() {
	var all []Tool
	all = append(all, s.metaTools...)
	for _, gt := range s.focused {
		all = append(all, gt...)
	}
	s.tools = all
	s.toolIndex = make(map[string]Tool, len(all))
	for _, t := range all {
		s.toolIndex[t.Name] = t
	}
	s.toolsVersion++
}

func (s *Server) notifyToolsChanged() {
	if s.w == nil {
		return
	}
	data, _ := json.Marshal(struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
	}{JSONRPC: "2.0", Method: "notifications/tools/list_changed"})
	s.mu.Lock()
	fmt.Fprintf(s.w, "%s\n", data)
	s.mu.Unlock()
}

func (s *Server) Serve(ctx context.Context) error {
	return s.serve(ctx, os.Stdin, os.Stdout)
}

func (s *Server) serve(ctx context.Context, r io.Reader, w io.Writer) error {
	s.w = w
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<20)

	var wg sync.WaitGroup
	defer wg.Wait()

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
			s.writeError(nil, codeParseError, "parse error")
			continue
		}

		if msg.Method == "" {
			continue
		}

		if msg.Method == "tools/call" {
			wg.Add(1)
			go func(m jsonrpcMessage) {
				defer wg.Done()
				s.dispatch(&m)
			}(msg)
			continue
		}

		s.dispatch(&msg)
	}

	return scanner.Err()
}

func (s *Server) dispatch(msg *jsonrpcMessage) {
	switch msg.Method {
	case "initialize":
		tc := &toolCapability{}
		if s.root != nil {
			tc.ListChanged = true
		}
		result := initializeResult{
			ProtocolVersion: protocolVersion,
			Capabilities: capabilities{
				Tools: tc,
			},
			ServerInfo:   serverInfo{Name: s.name, Version: s.version},
			Instructions: serverInstructions,
		}
		if len(s.resources) > 0 {
			result.Capabilities.Resources = &resourceCapability{}
		}
		s.writeResult(msg.ID, result)

	case "notifications/initialized":

	case "tools/list":
		s.mu.Lock()
		defs := make([]toolDef, len(s.tools))
		for i, t := range s.tools {
			defs[i] = toolDef{
				Name:        t.Name,
				Description: t.Description,
				InputSchema: t.InputSchema,
			}
		}
		s.mu.Unlock()
		s.writeResult(msg.ID, toolListResult{Tools: defs})

	case "tools/call":
		s.handleToolCall(msg)

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
		s.writeResult(msg.ID, resourceListResult{Resources: defs})

	case "resources/read":
		s.handleResourceRead(msg)

	case "ping":
		s.writeResult(msg.ID, struct{}{})

	default:
		s.writeError(msg.ID, codeMethodNotFound, fmt.Sprintf("unknown method: %s", msg.Method))
	}
}

func (s *Server) handleToolCall(msg *jsonrpcMessage) {
	var params toolCallParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.writeError(msg.ID, codeInvalidParams, "invalid params")
		return
	}

	s.mu.Lock()
	tool, ok := s.toolIndex[params.Name]
	prevVersion := s.toolsVersion
	s.mu.Unlock()

	if !ok {
		s.writeError(msg.ID, codeInvalidParams, fmt.Sprintf("unknown tool: %s", params.Name))
		return
	}

	output, err := tool.Run(params.Arguments)
	if err != nil {
		errText := err.Error()
		if output != "" {
			errText = output
		}
		s.writeResult(msg.ID, toolCallResult{
			Content: []content{{Type: "text", Text: fmt.Sprintf("error: %s", errText)}},
			IsError: true,
		})
		return
	}

	s.writeResult(msg.ID, toolCallResult{
		Content: []content{{Type: "text", Text: output}},
	})

	s.mu.Lock()
	changed := s.toolsVersion != prevVersion
	s.mu.Unlock()
	if changed {
		s.notifyToolsChanged()
	}
}

func (s *Server) handleResourceRead(msg *jsonrpcMessage) {
	var params resourceReadParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.writeError(msg.ID, codeInvalidParams, "invalid params")
		return
	}

	res, ok := s.resIndex[params.URI]
	if !ok {
		s.writeError(msg.ID, codeInvalidParams, fmt.Sprintf("unknown resource: %s", params.URI))
		return
	}

	text, err := res.Read()
	if err != nil {
		s.writeError(msg.ID, codeInternalError, err.Error())
		return
	}

	s.writeResult(msg.ID, resourceReadResult{
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
	Instructions    string       `json:"instructions,omitempty"`
}

const serverInstructions = `s1ctl — CLI and SDK for SentinelOne Singularity Platform.

Discovery flow:
1. help → list command groups
2. help {group} → list subcommands
3. usage {command} → flags, args, and description for one command
4. focus {group} → load typed tool schemas (enables structured calls)
5. run {command} → execute any command (e.g. "agents list --site-id 123 --limit 5")

Use "run" for quick one-off commands. Use "focus" when you need repeated structured calls within a group.
Prefer "focus" for commands with filter expressions — typed parameters avoid quoting issues.
In "run", use shell-style quoting for values with spaces: --filter 'event.type = "Login"'.

All mutations are dry-run by default — pass --yes to apply.
Always scope to the correct --site-id.
Output is JSON by default when called via MCP.`

type capabilities struct {
	Tools     *toolCapability     `json:"tools,omitempty"`
	Resources *resourceCapability `json:"resources,omitempty"`
}

type toolCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

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

func (s *Server) writeResult(id json.RawMessage, result any) {
	resp := jsonrpcResponse{JSONRPC: "2.0", ID: id, Result: result}
	data, _ := json.Marshal(resp)
	s.mu.Lock()
	fmt.Fprintf(s.w, "%s\n", data)
	s.mu.Unlock()
}

func (s *Server) writeError(id json.RawMessage, code int, message string) {
	resp := jsonrpcResponse{JSONRPC: "2.0", ID: id, Error: &jsonrpcError{Code: code, Message: message}}
	data, _ := json.Marshal(resp)
	s.mu.Lock()
	fmt.Fprintf(s.w, "%s\n", data)
	s.mu.Unlock()
}
