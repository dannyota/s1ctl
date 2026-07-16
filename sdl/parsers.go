package sdl

import (
	"context"
	"encoding/json"
	"fmt"
)

const parsersListGQL = `query getConfigurationFiles {
  configFiles {
    udoId
    name
    readOnly
    version
  }
}`

const parserGetGQL = `query configFile($udoId: ID!) {
  configFile(udoId: $udoId) {
    udoId
    name
    content
    createdDate
    modifiedDate
    readOnly
    version
  }
}`

const parserCreateGQL = `mutation addConfigurationFile($name: String, $udoId: ID, $content: String, $expectedVersion: Long) {
  addConfigFile(name: $name, udoId: $udoId, content: $content, expectedVersion: $expectedVersion) {
    udoId
    name
    content
    createdDate
    modifiedDate
    readOnly
    version
  }
}`

const parserDeleteGQL = `mutation deleteConfigurationFile($udoId: ID!, $expectedVersion: Long) {
  deleteConfigFile(udoId: $udoId, expectedVersion: $expectedVersion) {
    udoId
    name
  }
}`

// Parser is an SDL Data Lake configuration file (parser).
type Parser struct {
	UdoID    string `json:"udoId"`
	Name     string `json:"name"`
	ReadOnly bool   `json:"readOnly"`
	Version  int    `json:"version"`

	Raw json.RawMessage `json:"-"`
}

func (p *Parser) UnmarshalJSON(b []byte) error {
	type alias Parser
	if err := json.Unmarshal(b, (*alias)(p)); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0:0], b...)
	return nil
}

// ParserDetail is the full representation of an SDL parser including content
// and timestamps.
type ParserDetail struct {
	UdoID        string `json:"udoId"`
	Name         string `json:"name"`
	Content      string `json:"content"`
	CreatedDate  string `json:"createdDate"`
	ModifiedDate string `json:"modifiedDate"`
	ReadOnly     bool   `json:"readOnly"`
	Version      int    `json:"version"`

	Raw json.RawMessage `json:"-"`
}

func (p *ParserDetail) UnmarshalJSON(b []byte) error {
	type alias ParserDetail
	if err := json.Unmarshal(b, (*alias)(p)); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0:0], b...)
	return nil
}

// ParserCreateInput holds parameters for creating or updating a parser.
type ParserCreateInput struct {
	Name            *string `json:"name,omitempty"`
	UdoID           *string `json:"udoId,omitempty"`
	Content         *string `json:"content,omitempty"`
	ExpectedVersion *int    `json:"expectedVersion,omitempty"`
}

// ParsersList returns all configuration file parsers from the SDL console.
func (c *Client) ParsersList(ctx context.Context) ([]Parser, error) {
	var data struct {
		ConfigFiles []Parser `json:"configFiles"`
	}
	if err := c.graphql(ctx, parsersListGQL, nil, &data); err != nil {
		return nil, err
	}
	return data.ConfigFiles, nil
}

// ParserGet returns a single parser by UDO ID with full detail.
func (c *Client) ParserGet(ctx context.Context, udoID string) (*ParserDetail, error) {
	if udoID == "" {
		return nil, fmt.Errorf("sdl: parser udoId is required")
	}
	vars := map[string]any{"udoId": udoID}
	var data struct {
		ConfigFile ParserDetail `json:"configFile"`
	}
	if err := c.graphql(ctx, parserGetGQL, vars, &data); err != nil {
		return nil, err
	}
	return &data.ConfigFile, nil
}

// ParserCreate creates or updates a configuration file parser.
func (c *Client) ParserCreate(ctx context.Context, input *ParserCreateInput) (*ParserDetail, error) {
	if input == nil {
		return nil, fmt.Errorf("sdl: parser input is required")
	}
	vars := map[string]any{}
	if input.Name != nil {
		vars["name"] = *input.Name
	}
	if input.UdoID != nil {
		vars["udoId"] = *input.UdoID
	}
	if input.Content != nil {
		vars["content"] = *input.Content
	}
	if input.ExpectedVersion != nil {
		vars["expectedVersion"] = *input.ExpectedVersion
	}
	var data struct {
		AddConfigFile ParserDetail `json:"addConfigFile"`
	}
	if err := c.graphql(ctx, parserCreateGQL, vars, &data); err != nil {
		return nil, err
	}
	return &data.AddConfigFile, nil
}

// ParserDelete deletes a configuration file parser by UDO ID.
func (c *Client) ParserDelete(ctx context.Context, udoID string, expectedVersion *int) error {
	if udoID == "" {
		return fmt.Errorf("sdl: parser udoId is required")
	}
	vars := map[string]any{"udoId": udoID}
	if expectedVersion != nil {
		vars["expectedVersion"] = *expectedVersion
	}
	return c.graphql(ctx, parserDeleteGQL, vars, nil)
}
