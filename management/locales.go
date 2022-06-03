package management

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type LocaleResponse struct {
	Locale Locale `json:"locale"`
}

type LocaleRequest struct {
	Locale LocaleInput `json:"locale"`
}

// Locale represents the global field in contentstack.
type Locale struct {
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name,omitempty"`
	UID            string    `json:"uid,omitempty"`
	Code           string    `json:"code"`
	FallbackLocale string    `json:"fallback_locale"`
}

// LocaleInput is used to create or update a content type
type LocaleInput struct {
	Name           string `json:"name,omitempty"`
	Code           string `json:"code"`
	FallbackLocale string `json:"fallback_locale"`
}

func (si *StackInstance) LocaleCreate(ctx context.Context, input LocaleInput) (*Locale, error) {
	data, err := serializeInput(LocaleRequest{Locale: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.post(
		ctx,
		"/v3/locales/",
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &LocaleResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Locale, nil
}

func (si *StackInstance) LocaleUpdate(ctx context.Context, code string, input LocaleInput) (*Locale, error) {
	data, err := serializeInput(LocaleRequest{Locale: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.put(
		ctx,
		fmt.Sprintf("/v3/locales/%s", code),
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &LocaleResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Locale, nil
}

func (si *StackInstance) LocaleDelete(ctx context.Context, code string) error {
	resp, err := si.client.delete(
		ctx,
		fmt.Sprintf("/v3/locales/%s", code),
		url.Values{},
		si.headers(),
		nil,
	)

	if err != nil {
		return err
	}

	result := &LocaleResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return err
	}

	return nil
}

func (si *StackInstance) LocaleFetch(ctx context.Context, code string) (*Locale, error) {
	resp, err := si.client.get(
		ctx,
		fmt.Sprintf("/v3/locales/%s", code),
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := &LocaleResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Locale, nil
}

func (si *StackInstance) LocaleFetchAll(ctx context.Context) ([]Locale, error) {
	resp, err := si.client.get(
		ctx,
		"/v3/locales",
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := struct {
		Locales []Locale `json:"locales"`
	}{}

	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Locales, nil
}
