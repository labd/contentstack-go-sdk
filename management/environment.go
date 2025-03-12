package management

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type EnvironmentResponse struct {
	Environment Environment `json:"environment"`
}

type EnvironmentRequest struct {
	Environment EnvironmentInput `json:"environment"`
}

// Environment represents the environment in contentstack
type Environment struct {
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Name      string           `json:"name"`
	UID       string           `json:"uid,omitempty"`
	URLs      []EnvironmentUrl `json:"urls"`
}

type EnvironmentUrl struct {
	Locale string `json:"locale"`
	URL    string `json:"url"`
}

// EnvironmentInput is used to create or update an environment
type EnvironmentInput struct {
	Name string           `json:"name"`
	URLs []EnvironmentUrl `json:"urls"`
}

func (si *StackInstance) EnvironmentCreate(ctx context.Context, input EnvironmentInput) (*Environment, error) {
	data, err := serializeInput(EnvironmentRequest{Environment: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.post(
		ctx,
		"/v3/environments/",
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &EnvironmentResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Environment, nil
}

func (si *StackInstance) EnvironmentUpdate(ctx context.Context, name string, input EnvironmentInput) (*Environment, error) {
	data, err := serializeInput(EnvironmentRequest{Environment: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.put(
		ctx,
		fmt.Sprintf("/v3/environments/%s", name),
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &EnvironmentResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Environment, nil
}

func (si *StackInstance) EnvironmentDelete(ctx context.Context, name string) error {
	resp, err := si.client.delete(
		ctx,
		fmt.Sprintf("/v3/environments/%s", name),
		url.Values{},
		si.headers(),
		nil,
	)
	if err != nil {
		return err
	}

	result := &EnvironmentResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return err
	}

	return nil
}

func (si *StackInstance) EnvironmentFetch(ctx context.Context, name string) (*Environment, error) {
	resp, err := si.client.get(
		ctx,
		fmt.Sprintf("/v3/environments/%s", name),
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := &EnvironmentResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Environment, nil
}

func (si *StackInstance) EnvironmentFetchAll(ctx context.Context, name string) ([]Environment, error) {
	resp, err := si.client.get(
		ctx,
		"/v3/environments",
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := struct {
		Environments []Environment `json:"environments"`
	}{}

	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Environments, nil
}
