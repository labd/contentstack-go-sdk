package management

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type Stack struct {
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UID             string    `json:"uid"`
	OrganizationUID string    `json:"org_uid"`
	ApiKey          string    `json:"api_key"`
	Name            string    `json:"name"`
	MasterLocale    string    `json:"master_locale"`
}

type StacksInput struct {
	OrganizationUid string
	IncludeCount    bool
	Limit           int
	Skip            int
	Asc             string
	Desc            string
}

func (c *Client) Stacks(ctx context.Context, input StacksInput) ([]Stack, error) {
	header := http.Header{}
	header.Add("authtoken", c.authToken)

	resp, err := c.get(
		ctx,
		"/v3/stacks",
		url.Values{},
		header,
	)
	if err != nil {
		return nil, err
	}

	result := struct {
		Stacks []Stack `json:"stacks"`
	}{}

	err = c.processResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	return result.Stacks, nil
}
