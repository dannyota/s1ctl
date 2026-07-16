package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func testServer() *Server {
	tools := []Tool{
		{
			Name:        "echo",
			Description: "Echo the input back",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"text": map[string]any{
						"type":        "string",
						"description": "text to echo",
					},
				},
			},
			Run: func(args map[string]any) (string, error) {
				return args["text"].(string), nil
			},
		},
	}
	resources := []Resource{
		{
			URI:         "guide://test",
			Name:        "test",
			Description: "A test resource",
			MimeType:    "text/plain",
			Read:        func() (string, error) { return "hello from resource", nil },
		},
	}
	return NewServer("test-server", "1.0.0", tools, resources)
}

func roundTrip(t *testing.T, srv *Server, request string) map[string]any {
	t.Helper()
	msgs := roundTripMulti(t, srv, request)
	if len(msgs) == 0 {
		t.Fatal("no response")
	}
	return msgs[0]
}

func roundTripMulti(t *testing.T, srv *Server, request string) []map[string]any {
	t.Helper()
	var out bytes.Buffer
	r := strings.NewReader(request + "\n")
	if err := srv.serve(context.Background(), r, &out); err != nil {
		t.Fatal(err)
	}
	var results []map[string]any
	for line := range strings.SplitSeq(strings.TrimSpace(out.String()), "\n") {
		if line == "" {
			continue
		}
		var msg map[string]any
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			t.Fatalf("unmarshal: %v\nraw: %s", err, line)
		}
		results = append(results, msg)
	}
	return results
}

func TestInitialize(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`)

	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatal("no result in response")
	}
	if v := result["protocolVersion"]; v != protocolVersion {
		t.Errorf("protocolVersion = %v, want %v", v, protocolVersion)
	}
	info, _ := result["serverInfo"].(map[string]any)
	if info["name"] != "test-server" {
		t.Errorf("serverInfo.name = %v, want test-server", info["name"])
	}
	caps, _ := result["capabilities"].(map[string]any)
	if caps["tools"] == nil {
		t.Error("capabilities.tools missing")
	}
	if caps["resources"] == nil {
		t.Error("capabilities.resources missing")
	}
}

func TestToolsList(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)

	result, _ := resp["result"].(map[string]any)
	tools, _ := result["tools"].([]any)
	if len(tools) != 1 {
		t.Fatalf("got %d tools, want 1", len(tools))
	}
	tool, _ := tools[0].(map[string]any)
	if tool["name"] != "echo" {
		t.Errorf("tool name = %v, want echo", tool["name"])
	}
}

func TestToolsCall(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"echo","arguments":{"text":"hello world"}}}`)

	result, _ := resp["result"].(map[string]any)
	content, _ := result["content"].([]any)
	if len(content) != 1 {
		t.Fatalf("got %d content items, want 1", len(content))
	}
	item, _ := content[0].(map[string]any)
	if item["text"] != "hello world" {
		t.Errorf("text = %v, want hello world", item["text"])
	}
}

func TestToolsCallUnknown(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"nonexistent","arguments":{}}}`)

	if resp["error"] == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestResourcesList(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":5,"method":"resources/list","params":{}}`)

	result, _ := resp["result"].(map[string]any)
	resources, _ := result["resources"].([]any)
	if len(resources) != 1 {
		t.Fatalf("got %d resources, want 1", len(resources))
	}
	res, _ := resources[0].(map[string]any)
	if res["uri"] != "guide://test" {
		t.Errorf("uri = %v, want guide://test", res["uri"])
	}
}

func TestResourcesRead(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":6,"method":"resources/read","params":{"uri":"guide://test"}}`)

	result, _ := resp["result"].(map[string]any)
	contents, _ := result["contents"].([]any)
	if len(contents) != 1 {
		t.Fatalf("got %d contents, want 1", len(contents))
	}
	item, _ := contents[0].(map[string]any)
	if item["text"] != "hello from resource" {
		t.Errorf("text = %v, want hello from resource", item["text"])
	}
}

func TestResourcesReadUnknown(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":7,"method":"resources/read","params":{"uri":"guide://nonexistent"}}`)

	if resp["error"] == nil {
		t.Fatal("expected error for unknown resource")
	}
}

func TestPing(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":8,"method":"ping","params":{}}`)

	if resp["error"] != nil {
		t.Errorf("ping returned error: %v", resp["error"])
	}
	if resp["result"] == nil {
		t.Error("ping returned no result")
	}
}

func TestUnknownMethod(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":9,"method":"bogus/method","params":{}}`)

	if resp["error"] == nil {
		t.Fatal("expected error for unknown method")
	}
	errObj, _ := resp["error"].(map[string]any)
	code, _ := errObj["code"].(float64)
	if int(code) != codeMethodNotFound {
		t.Errorf("error code = %v, want %d", code, codeMethodNotFound)
	}
}

func TestMultipleMessages(t *testing.T) {
	srv := testServer()
	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
	}, "\n") + "\n"

	var out bytes.Buffer
	if err := srv.serve(context.Background(), strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d response lines, want 2 (initialize + tools/list; notification has no response)", len(lines))
	}
}

// --- Dynamic server tests ---

func testDynamicServer() *Server {
	root := testCobraTree()
	return NewDynamicServer("test-dynamic", "2.0.0", root, nil)
}

func TestListChangedCapability(t *testing.T) {
	srv := testDynamicServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`)

	result, _ := resp["result"].(map[string]any)
	caps, _ := result["capabilities"].(map[string]any)
	tools, _ := caps["tools"].(map[string]any)
	if tools["listChanged"] != true {
		t.Errorf("listChanged = %v, want true", tools["listChanged"])
	}
}

func TestStaticServerNoListChanged(t *testing.T) {
	srv := testServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`)

	result, _ := resp["result"].(map[string]any)
	caps, _ := result["capabilities"].(map[string]any)
	tools, _ := caps["tools"].(map[string]any)
	if tools["listChanged"] == true {
		t.Error("static server should not advertise listChanged")
	}
}

func TestDynamicServerMetaTools(t *testing.T) {
	srv := testDynamicServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`)

	result, _ := resp["result"].(map[string]any)
	tools, _ := result["tools"].([]any)

	names := make(map[string]bool)
	for _, t := range tools {
		tool, _ := t.(map[string]any)
		names[tool["name"].(string)] = true
	}

	for _, want := range []string{"run", "help", "usage", "focus", "unfocus"} {
		if !names[want] {
			t.Errorf("missing meta-tool %q", want)
		}
	}
	if len(tools) != 5 {
		t.Errorf("got %d tools, want 5 meta-tools only", len(tools))
	}
}

func TestFocusUnfocus(t *testing.T) {
	srv := testDynamicServer()

	// Focus on "agents" — should return response + notification.
	msgs := roundTripMulti(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"focus","arguments":{"group":"agents"}}}`)
	if len(msgs) < 2 {
		t.Fatalf("got %d messages, want at least 2 (response + notification)", len(msgs))
	}

	// First message is the response.
	result, _ := msgs[0]["result"].(map[string]any)
	content, _ := result["content"].([]any)
	if len(content) == 0 {
		t.Fatal("no content in focus response")
	}
	item, _ := content[0].(map[string]any)
	text, _ := item["text"].(string)
	if !strings.Contains(text, "agents") {
		t.Errorf("focus response should mention agents: %s", text)
	}

	// Second message is the notification.
	notif := msgs[1]
	if notif["method"] != "notifications/tools/list_changed" {
		t.Errorf("notification method = %v, want notifications/tools/list_changed", notif["method"])
	}

	// Now tools/list should have meta-tools + agents tools.
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)
	toolsResult, _ := resp["result"].(map[string]any)
	tools, _ := toolsResult["tools"].([]any)
	if len(tools) <= 5 {
		t.Errorf("after focus, got %d tools, want more than 5", len(tools))
	}

	hasAgentsList := false
	for _, tt := range tools {
		tool, _ := tt.(map[string]any)
		if tool["name"] == "agents_list" {
			hasAgentsList = true
		}
	}
	if !hasAgentsList {
		t.Error("agents_list should be present after focusing on agents")
	}

	// Unfocus — should shrink back to 5.
	msgs = roundTripMulti(t, srv, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"unfocus","arguments":{"group":"agents"}}}`)
	if len(msgs) < 2 {
		t.Fatalf("unfocus: got %d messages, want at least 2", len(msgs))
	}

	resp = roundTrip(t, srv, `{"jsonrpc":"2.0","id":4,"method":"tools/list","params":{}}`)
	toolsResult, _ = resp["result"].(map[string]any)
	tools, _ = toolsResult["tools"].([]any)
	if len(tools) != 5 {
		t.Errorf("after unfocus, got %d tools, want 5", len(tools))
	}
}

func TestHelpTool(t *testing.T) {
	srv := testDynamicServer()

	// Help with no group — lists all groups.
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"help","arguments":{}}}`)
	result, _ := resp["result"].(map[string]any)
	content, _ := result["content"].([]any)
	item, _ := content[0].(map[string]any)
	text, _ := item["text"].(string)
	if !strings.Contains(text, "agents") {
		t.Errorf("help output should list agents group: %s", text)
	}

	// Help with group — lists subcommands.
	resp = roundTrip(t, srv, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"help","arguments":{"group":"agents"}}}`)
	result, _ = resp["result"].(map[string]any)
	content, _ = result["content"].([]any)
	item, _ = content[0].(map[string]any)
	text, _ = item["text"].(string)
	if !strings.Contains(text, "list") || !strings.Contains(text, "isolate") {
		t.Errorf("help agents should list subcommands: %s", text)
	}
}

func TestFocusUnknownGroup(t *testing.T) {
	srv := testDynamicServer()
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"focus","arguments":{"group":"nonexistent"}}}`)

	result, _ := resp["result"].(map[string]any)
	if result["isError"] != true {
		t.Error("focusing unknown group should return isError")
	}
}

func TestToolCallErrorIncludesMessage(t *testing.T) {
	tools := []Tool{
		{
			Name:        "fail",
			Description: "Always fails with a message",
			InputSchema: map[string]any{"type": "object"},
			Run: func(args map[string]any) (string, error) {
				return "", fmt.Errorf("connection refused: dial tcp 10.0.0.1:443")
			},
		},
	}
	srv := NewServer("test", "1.0.0", tools, nil)
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"fail","arguments":{}}}`)

	result, _ := resp["result"].(map[string]any)
	if result["isError"] != true {
		t.Fatal("expected isError")
	}
	content, _ := result["content"].([]any)
	if len(content) == 0 {
		t.Fatal("expected content")
	}
	text, _ := content[0].(map[string]any)["text"].(string)

	var env struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(text), &env); err != nil {
		t.Fatalf("error text is not JSON envelope: %v\n%s", err, text)
	}
	if !strings.Contains(env.Error.Message, "connection refused") {
		t.Errorf("error.message = %q, want it to contain the actual error message", env.Error.Message)
	}
}

func TestReadOnlyModeFiltersMutations(t *testing.T) {
	root := testCobraTree()
	srv := NewDynamicServer("test-ro", "1.0.0", root, nil, WithReadOnly(true))
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`)

	result, _ := resp["result"].(map[string]any)
	tools, _ := result["tools"].([]any)
	for _, tt := range tools {
		tool, _ := tt.(map[string]any)
		name, _ := tool["name"].(string)
		if name == "agents_isolate" {
			t.Error("read-only server should not expose mutation tool agents_isolate")
		}
	}
}

func TestReadOnlyModeFocusFiltersMutations(t *testing.T) {
	root := testCobraTree()
	srv := NewDynamicServer("test-ro", "1.0.0", root, nil, WithReadOnly(true))
	roundTripMulti(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"focus","arguments":{"group":"agents"}}}`)
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)

	result, _ := resp["result"].(map[string]any)
	tools, _ := result["tools"].([]any)
	for _, tt := range tools {
		tool, _ := tt.(map[string]any)
		name, _ := tool["name"].(string)
		if name == "agents_isolate" {
			t.Error("read-only server should not expose mutation tool agents_isolate after focus")
		}
	}
	// agents_list should still be present.
	hasAgentsList := false
	for _, tt := range tools {
		tool, _ := tt.(map[string]any)
		if tool["name"] == "agents_list" {
			hasAgentsList = true
		}
	}
	if !hasAgentsList {
		t.Error("agents_list should be present after focus in read-only mode")
	}
}

func TestReadOnlyInstructions(t *testing.T) {
	root := testCobraTree()
	srv := NewDynamicServer("test-ro", "1.0.0", root, nil, WithReadOnly(true))
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`)

	result, _ := resp["result"].(map[string]any)
	instructions, _ := result["instructions"].(string)
	if !strings.Contains(instructions, "Read-only mode") {
		t.Error("read-only server instructions should mention read-only mode")
	}
}

func TestReadOnlyHelpOmitsMutations(t *testing.T) {
	root := testCobraTree()
	srv := NewDynamicServer("test-ro", "1.0.0", root, nil, WithReadOnly(true))

	// Help for the agents group should not show isolate (a mutation).
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"help","arguments":{"group":"agents"}}}`)
	result, _ := resp["result"].(map[string]any)
	content, _ := result["content"].([]any)
	text, _ := content[0].(map[string]any)["text"].(string)
	if strings.Contains(text, "[mutation]") {
		t.Errorf("read-only help should not contain [mutation], got:\n%s", text)
	}
	if strings.Contains(text, "isolate") {
		t.Errorf("read-only help should not list mutation command 'isolate', got:\n%s", text)
	}
	if !strings.Contains(text, "list") {
		t.Errorf("read-only help should still list read commands, got:\n%s", text)
	}

	// Top-level help should show 0 mutations.
	resp = roundTrip(t, srv, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"help","arguments":{}}}`)
	result, _ = resp["result"].(map[string]any)
	content, _ = result["content"].([]any)
	text, _ = content[0].(map[string]any)["text"].(string)
	if strings.Contains(text, "/1m") {
		t.Errorf("read-only top-level help should not report mutations, got:\n%s", text)
	}
}

func TestReadOnlyRunToolDescription(t *testing.T) {
	root := testCobraTree()
	srv := NewDynamicServer("test-ro", "1.0.0", root, nil, WithReadOnly(true))
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`)

	result, _ := resp["result"].(map[string]any)
	tools, _ := result["tools"].([]any)
	for _, tt := range tools {
		tool, _ := tt.(map[string]any)
		if tool["name"] == "run" {
			desc, _ := tool["description"].(string)
			if !strings.Contains(desc, "read-only mode") {
				t.Errorf("run tool description in read-only mode = %q, want read-only mention", desc)
			}
			return
		}
	}
	t.Error("run tool not found")
}

func TestToolAnnotations(t *testing.T) {
	root := testCobraTree()
	srv := NewDynamicServer("test-ann", "1.0.0", root, nil)
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`)

	result, _ := resp["result"].(map[string]any)
	tools, _ := result["tools"].([]any)
	for _, tt := range tools {
		tool, _ := tt.(map[string]any)
		name, _ := tool["name"].(string)
		ann, _ := tool["annotations"].(map[string]any)
		switch name {
		case "help", "usage", "focus", "unfocus":
			if ann == nil {
				t.Errorf("%s: expected annotations", name)
				continue
			}
			if ann["readOnlyHint"] != true {
				t.Errorf("%s: readOnlyHint = %v, want true", name, ann["readOnlyHint"])
			}
			if ann["destructiveHint"] != false {
				t.Errorf("%s: destructiveHint = %v, want false", name, ann["destructiveHint"])
			}
		case "run":
			if ann != nil {
				t.Errorf("run: should have no annotations in normal mode, got %v", ann)
			}
		}
	}
}

func TestToolAnnotationsOnFocusedTools(t *testing.T) {
	root := testCobraTree()
	srv := NewDynamicServer("test-ann", "1.0.0", root, nil)
	roundTripMulti(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"focus","arguments":{"group":"agents"}}}`)
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)

	result, _ := resp["result"].(map[string]any)
	tools, _ := result["tools"].([]any)
	for _, tt := range tools {
		tool, _ := tt.(map[string]any)
		name, _ := tool["name"].(string)
		ann, _ := tool["annotations"].(map[string]any)
		switch name {
		case "agents_list":
			if ann == nil {
				t.Fatal("agents_list: expected annotations")
			}
			if ann["readOnlyHint"] != true {
				t.Errorf("agents_list: readOnlyHint = %v, want true", ann["readOnlyHint"])
			}
			if ann["destructiveHint"] != false {
				t.Errorf("agents_list: destructiveHint = %v, want false", ann["destructiveHint"])
			}
		case "agents_isolate":
			if ann == nil {
				t.Fatal("agents_isolate: expected annotations")
			}
			if ann["readOnlyHint"] != false {
				t.Errorf("agents_isolate: readOnlyHint = %v, want false", ann["readOnlyHint"])
			}
			if ann["destructiveHint"] != true {
				t.Errorf("agents_isolate: destructiveHint = %v, want true", ann["destructiveHint"])
			}
		}
	}
}

func TestToolCallErrorKeepsOutputAndCause(t *testing.T) {
	tools := []Tool{
		{
			Name:        "partial",
			Description: "Fails but has output",
			InputSchema: map[string]any{"type": "object"},
			Run: func(args map[string]any) (string, error) {
				return `{"detail":"Bad Request - could not parse query"}`, fmt.Errorf("exit status 1")
			},
		},
	}
	srv := NewServer("test", "1.0.0", tools, nil)
	resp := roundTrip(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"partial","arguments":{}}}`)

	result, _ := resp["result"].(map[string]any)
	if result["isError"] != true {
		t.Fatal("expected isError")
	}
	content, _ := result["content"].([]any)
	text, _ := content[0].(map[string]any)["text"].(string)

	var env struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(text), &env); err != nil {
		t.Fatalf("error text is not JSON envelope: %v\n%s", err, text)
	}
	if !strings.Contains(env.Error.Message, "could not parse query") {
		t.Errorf("error.message = %q, want partial output included", env.Error.Message)
	}
	if !strings.Contains(env.Error.Message, "exit status 1") {
		t.Errorf("error.message = %q, want the error cause included", env.Error.Message)
	}
}
