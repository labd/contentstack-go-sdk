package management

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type WebHookResponse struct {
	WebHook WebHook `json:"webhook"`
}

type WebHookRequest struct {
	WebHook WebHookInput `json:"webhook"`
}

// WebHook represents the content type in contentstack.
type WebHook struct {
	UID             string               `json:"uid,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	CreatedBy       string               `json:"created_by"`
	UpdatedBy       string               `json:"updated_by"`
	Name            string               `json:"name"`
	OrganizationUID string               `json:"org_uid,omitempty"`
	Channels        []string             `json:"channels"`
	Branches        []string             `json:"branches"`
	Destinations    []WebhookDestination `json:"destinations"`
	RetryPolicy     string               `json:"retry_policy"`
	Disabled        bool                 `json:"disabled"`
	ConcisePayload  bool                 `json:"concise_payload"`
}

type WebhookDestination struct {
	TargetURL         string          `json:"target_url"`
	HttpBasicAuth     string          `json:"http_basic_auth"`
	HttpBasicPassword string          `json:"http_basic_password"`
	CustomHeader      []WebhookHeader `json:"custom_header"`
}

type WebhookHeader struct {
	Name  string `json:"header_name"`
	Value string `json:"value"`
}

// WebHookInput is used to create or update a content type
type WebHookInput struct {
	Name           string               `json:"name"`
	Branches       []string             `json:"branches"`
	Channels       []string             `json:"channels"`
	Destinations   []WebhookDestination `json:"destinations"`
	RetryPolicy    string               `json:"retry_policy"`
	Disabled       bool                 `json:"disabled"`
	ConcisePayload bool                 `json:"concise_payload"`
}

func (si *StackInstance) WebHookCreate(ctx context.Context, input WebHookInput) (*WebHook, error) {
	data, err := serializeInput(WebHookRequest{WebHook: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.post(
		ctx,
		"/v3/webhooks/",
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &WebHookResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.WebHook, nil
}

func (si *StackInstance) WebHookUpdate(ctx context.Context, uid string, input WebHookInput) (*WebHook, error) {
	data, err := serializeInput(WebHookRequest{WebHook: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.put(
		ctx,
		fmt.Sprintf("/v3/webhooks/%s", uid),
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &WebHookResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.WebHook, nil
}

func (si *StackInstance) WebHookDelete(ctx context.Context, uid string) error {
	resp, err := si.client.delete(
		ctx,
		fmt.Sprintf("/v3/webhooks/%s", uid),
		url.Values{},
		si.headers(),
		nil,
	)

	if err != nil {
		return err
	}

	result := &WebHookResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return err
	}

	return nil
}

func (si *StackInstance) WebHookFetch(ctx context.Context, uid string) (*WebHook, error) {
	resp, err := si.client.get(
		ctx,
		fmt.Sprintf("/v3/webhooks/%s", uid),
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := &WebHookResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.WebHook, nil
}

func (si *StackInstance) WebHookFetchAll(ctx context.Context) ([]WebHook, error) {
	resp, err := si.client.get(
		ctx,
		"/v3/webhooks",
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := struct {
		WebHooks []WebHook `json:"webhooks"`
	}{}

	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.WebHooks, nil
}
