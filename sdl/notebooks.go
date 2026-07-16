package sdl

import (
	"context"
	"encoding/json"
	"fmt"
)

const notebooksListGQL = `query notebooks {
  purpleConversations {
    accountId
    createdAt
    description
    expiresAt
    id
    isReadOnly
    isAppendable
    isShared
    ownerEmail
    name
    teamToken
    notebookSource
  }
}`

const notebookGetGQL = `query purpleNotebook($id: ID!) {
  purpleConversation(id: $id) {
    id
    accountId
    createdAt
    description
    expiresAt
    isReadOnly
    isAppendable
    isShared
    ownerEmail
    name
    teamToken
    notebookSource
    entitlements { account }
  }
}`

const notebookCreateGQL = `mutation createNotebook($name: String!, $description: String!) {
  createPurpleConversation(name: $name, description: $description) {
    accountId
    createdAt
    description
    expiresAt
    id
    isReadOnly
    isAppendable
    isShared
    ownerEmail
    name
    teamToken
    notebookSource
  }
}`

const notebookUpdateGQL = `mutation updateNotebook($id: ID!, $name: String, $description: String) {
  updatePurpleConversation(id: $id, name: $name, description: $description) {
    id
    name
    description
  }
}`

const notebookDeleteGQL = `mutation deleteNotebook($id: ID!) {
  deletePurpleConversation(id: $id)
}`

// Notebook is a Purple AI notebook (conversation).
type Notebook struct {
	ID             string `json:"id"`
	AccountID      string `json:"accountId"`
	CreatedAt      string `json:"createdAt"`
	Description    string `json:"description"`
	ExpiresAt      string `json:"expiresAt"`
	IsReadOnly     bool   `json:"isReadOnly"`
	IsAppendable   bool   `json:"isAppendable"`
	IsShared       bool   `json:"isShared"`
	OwnerEmail     string `json:"ownerEmail"`
	Name           string `json:"name"`
	TeamToken      string `json:"teamToken"`
	NotebookSource string `json:"notebookSource"`

	Raw json.RawMessage `json:"-"`
}

func (n *Notebook) UnmarshalJSON(b []byte) error {
	type alias Notebook
	if err := json.Unmarshal(b, (*alias)(n)); err != nil {
		return err
	}
	n.Raw = append(n.Raw[:0:0], b...)
	return nil
}

// NotebookDetail is the full representation of a Purple AI notebook
// including entitlements.
type NotebookDetail struct {
	ID             string `json:"id"`
	AccountID      string `json:"accountId"`
	CreatedAt      string `json:"createdAt"`
	Description    string `json:"description"`
	ExpiresAt      string `json:"expiresAt"`
	IsReadOnly     bool   `json:"isReadOnly"`
	IsAppendable   bool   `json:"isAppendable"`
	IsShared       bool   `json:"isShared"`
	OwnerEmail     string `json:"ownerEmail"`
	Name           string `json:"name"`
	TeamToken      string `json:"teamToken"`
	NotebookSource string `json:"notebookSource"`

	Entitlements json.RawMessage `json:"entitlements"`

	Raw json.RawMessage `json:"-"`
}

func (n *NotebookDetail) UnmarshalJSON(b []byte) error {
	type alias NotebookDetail
	if err := json.Unmarshal(b, (*alias)(n)); err != nil {
		return err
	}
	n.Raw = append(n.Raw[:0:0], b...)
	return nil
}

// NotebookUpdateInput holds parameters for updating a notebook.
type NotebookUpdateInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// NotebooksList returns all Purple AI notebooks.
func (c *Client) NotebooksList(ctx context.Context) ([]Notebook, error) {
	var data struct {
		PurpleConversations []Notebook `json:"purpleConversations"`
	}
	if err := c.graphql(ctx, notebooksListGQL, nil, &data); err != nil {
		return nil, err
	}
	return data.PurpleConversations, nil
}

// NotebookGet returns a single notebook by ID with full detail.
func (c *Client) NotebookGet(ctx context.Context, id string) (*NotebookDetail, error) {
	if id == "" {
		return nil, fmt.Errorf("sdl: notebook id is required")
	}
	vars := map[string]any{"id": id}
	var data struct {
		PurpleConversation NotebookDetail `json:"purpleConversation"`
	}
	if err := c.graphql(ctx, notebookGetGQL, vars, &data); err != nil {
		return nil, err
	}
	return &data.PurpleConversation, nil
}

// NotebookCreate creates a new Purple AI notebook.
func (c *Client) NotebookCreate(ctx context.Context, name, description string) (*Notebook, error) {
	if name == "" {
		return nil, fmt.Errorf("sdl: notebook name is required")
	}
	vars := map[string]any{
		"name":        name,
		"description": description,
	}
	var data struct {
		CreatePurpleConversation Notebook `json:"createPurpleConversation"`
	}
	if err := c.graphql(ctx, notebookCreateGQL, vars, &data); err != nil {
		return nil, err
	}
	return &data.CreatePurpleConversation, nil
}

// NotebookUpdate updates a notebook's name and/or description.
func (c *Client) NotebookUpdate(ctx context.Context, id string, input *NotebookUpdateInput) error {
	if id == "" {
		return fmt.Errorf("sdl: notebook id is required")
	}
	if input == nil {
		return fmt.Errorf("sdl: notebook update input is required")
	}
	vars := map[string]any{"id": id}
	if input.Name != nil {
		vars["name"] = *input.Name
	}
	if input.Description != nil {
		vars["description"] = *input.Description
	}
	return c.graphql(ctx, notebookUpdateGQL, vars, nil)
}

// NotebookDelete deletes a notebook by ID.
func (c *Client) NotebookDelete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("sdl: notebook id is required")
	}
	vars := map[string]any{"id": id}
	return c.graphql(ctx, notebookDeleteGQL, vars, nil)
}
