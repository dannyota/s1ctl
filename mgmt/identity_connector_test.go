package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestIdentityConnectors(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/getCloudlinkConfigurations") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"cloudlinkId": 10, "mgmtId": 100,
					"status": "ACTIVE", "computerName": "DC1",
					"agentType": "full", "osName": "Windows Server 2022",
					"version": "23.4.1", "guid": "00000000-0000-0000-0000-000000000001",
					"isUnifiedAgent": true, "ipAddress": "10.0.0.1",
					"domainName": "corp.example.com", "lastSeen": "2024-01-01T00:00:00Z",
					"scopePath": "Global / Site1",
				},
			},
		})
	})
	c := testClient(t, handler)
	connectors, err := c.IdentityConnectors(context.Background(), &IdentityParams{SiteIDs: "100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(connectors) != 1 {
		t.Fatalf("expected 1 connector, got %d", len(connectors))
	}
	cn := connectors[0]
	if cn.CloudlinkID != 10 {
		t.Fatalf("unexpected cloudlinkId: %d", cn.CloudlinkID)
	}
	if cn.Status != ConnectorStatusActive {
		t.Fatalf("unexpected status: %s", cn.Status)
	}
	if cn.ComputerName != "DC1" {
		t.Fatalf("unexpected computerName: %s", cn.ComputerName)
	}
	if !cn.IsUnifiedAgent {
		t.Fatal("expected isUnifiedAgent=true")
	}
	if cn.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestIdentityConnector(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/getCloudlinkConfiguration") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"mgmtId": 100, "guid": "00000000-0000-0000-0000-000000000001",
				"status": "ACTIVE", "osName": "Windows Server 2022",
				"version": "23.4.1", "computerName": "DC1",
				"domainName": "corp.example.com", "ipAddress": "10.0.0.1",
				"isUnifiedAgent": true,
			},
		})
	})
	c := testClient(t, handler)
	cn, err := c.IdentityConnector(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cn.ComputerName != "DC1" {
		t.Fatalf("unexpected computerName: %s", cn.ComputerName)
	}
	if cn.Status != ConnectorStatusActive {
		t.Fatalf("unexpected status: %s", cn.Status)
	}
	if cn.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestIdentityConnectorReplace(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/replaceAdConnector") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("agentUuid") != "test-uuid" {
			t.Fatalf("expected agentUuid=test-uuid, got %s", r.URL.Query().Get("agentUuid"))
		}
		json.NewEncoder(w).Encode(map[string]any{"data": nil})
	})
	c := testClient(t, handler)
	err := c.IdentityConnectorReplace(context.Background(), &IdentityParams{SiteIDs: "100"}, "test-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIdentityWindowsAgents(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/getWindowsUnifiedAgents") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": 1, "mgmtId": 100, "uuid": "agent-001",
					"osName": "Windows 11", "ipAddress": "10.0.0.2",
					"agentVersion": "23.4.1", "agentType": "full",
					"domainName": "corp.example.com", "status": "active",
					"hostName": "WS01", "scopePath": "Global / Site1",
				},
			},
		})
	})
	c := testClient(t, handler)
	agents, err := c.IdentityWindowsAgents(context.Background(), &WindowsAgentParams{
		SiteIDs:     "100",
		FilterInput: "WS",
		RequestID:   "req-001",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}
	a := agents[0]
	if a.UUID != "agent-001" {
		t.Fatalf("unexpected uuid: %s", a.UUID)
	}
	if a.HostName != "WS01" {
		t.Fatalf("unexpected hostName: %s", a.HostName)
	}
	if a.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}
