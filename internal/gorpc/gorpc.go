package gorpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Endpoint string
	Client   *http.Client
}

type NodeResp struct {
}

func (c *Client) Node(ctx context.Context) (*NodeResp, error) {
	nodeUrl := fmt.Sprintf("%s/", c.Endpoint)
	var nodeStatus NodeResp
	if err := c.get(ctx, nodeUrl, &nodeStatus); err != nil {
		return nil, fmt.Errorf("gorpc: Node %w", err)
	}

	return &nodeStatus, nil
}

func (c *Client) validStatus(status int) bool {
	return status >= 200 && status < 300
}

func (c *Client) get(ctx context.Context, endpoint string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.validStatus(resp.StatusCode) {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	errorBuf, _ := ioutil.ReadAll(resp.Body) // ignore returned error
	return fmt.Errorf("HTTP %v: %s", resp.Status, errorBuf)
}
