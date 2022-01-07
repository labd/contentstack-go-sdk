package management

import (
	"net/http"
)

type StackInstance struct {
	client *Client
	auth   StackAuth
}

type StackAuth struct {
	ApiKey          string
	ManagementToken string
}

func (c *Client) Stack(s *StackAuth) *StackInstance {
	return &StackInstance{
		client: c,
		auth:   *s,
	}
}

func (si *StackInstance) headers() http.Header {
	header := http.Header{}
	header.Add("api_key", si.auth.ApiKey)
	if si.auth.ManagementToken != "" {
		header.Add("authorization", si.auth.ManagementToken)
	} else {
		header.Add("authtoken", si.client.authToken)
	}
	return header
}
