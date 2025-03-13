package management

import (
	"fmt"
	"net/http"
)

type StackInstance struct {
	client *Client
	auth   StackAuth
}

type StackAuth struct {
	ApiKey          string
	ManagementToken string
	Branch          string
}

// Stack creates a new StackInstance which can be used for actions on the
// given stack instance.
func (c *Client) Stack(s *StackAuth) (*StackInstance, error) {

	if c.authToken == "" && s.ManagementToken == "" {
		return nil, fmt.Errorf("the management token is required when no auth token is used")
	}

	instance := &StackInstance{
		client: c,
		auth:   *s,
	}

	return instance, nil
}

func (si *StackInstance) headers() http.Header {
	header := http.Header{}
	header.Add("api_key", si.auth.ApiKey)
	if si.auth.ManagementToken != "" {
		header.Add("authorization", si.auth.ManagementToken)
	} else {
		header.Add("authtoken", si.client.authToken)
	}
	if si.auth.Branch != "" {
		header.Add("branch", si.auth.Branch)
	}
	return header
}
