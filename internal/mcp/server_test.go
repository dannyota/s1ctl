package mcp

import (
	"bytes"
	"context"
	"encoding/json"
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
	var out bytes.Buffer
	r := strings.NewReader(request + "\n")
	if err := srv.serve(context.Background(), r, &out); err != nil {
		t.Fatal(err)
	}
	var resp map[string]any
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v\nraw: %s", err, out.String())
	}
	return resp
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
