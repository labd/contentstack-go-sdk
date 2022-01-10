package management

import (
	"context"
	"net/url"
)

type StackSettings struct {
}

func (si *StackInstance) Settings(ctx context.Context) (*StackSettings, error) {

	resp, err := si.client.get(
		ctx,
		"/v3/settings",
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := struct {
		StackSettings *StackSettings `json:"stack_settings"`
	}{}

	err = si.client.processResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	return result.StackSettings, nil

}
