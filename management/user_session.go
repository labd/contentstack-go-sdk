package management

import (
	"context"
	"net/http"
	"net/url"
)

func (c *Client) Login(ctx context.Context, s UserCredentials) error {
	data, err := serializeInput(struct {
		User UserCredentials `json:"user"`
	}{
		User: s,
	})
	if err != nil {
		return err
	}

	resp, err := c.post(
		ctx,
		"/v3/user-session",
		url.Values{},
		http.Header{},
		data,
	)

	if err != nil {
		return err
	}

	result := make(map[string]interface{})
	err = c.processResponse(resp, &result)
	if err != nil {
		return err
	}

	c.authToken = result["user"].(map[string]interface{})["authtoken"].(string)
	return nil

}
