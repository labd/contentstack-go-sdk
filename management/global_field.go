package management

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type GlobalFieldResponse struct {
	GlobalField GlobalField `json:"global_field"`
}

type GlobalFieldRequest struct {
	GlobalField GlobalFieldInput `json:"global_field"`
}

// GlobalField represents the global field in contentstack.
type GlobalField struct {
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	Title             string          `json:"title,omitempty"`
	UID               string          `json:"uid,omitempty"`
	Schema            json.RawMessage `json:"schema"`
	MaintainRevisions bool            `json:"maintain_revisions"`
	Description       string          `json:"description"`
}

// GlobalFieldInput is used to create or update a content type
type GlobalFieldInput struct {
	Title             *string         `json:"title,omitempty"`
	UID               *string         `json:"uid,omitempty"`
	Description       *string         `json:"description,omitempty"`
	MaintainRevisions bool            `json:"maintain_revisions"`
	Schema            json.RawMessage `json:"schema,omitempty"`
}

func (si *StackInstance) GlobalFieldCreate(ctx context.Context, input GlobalFieldInput) (*GlobalField, error) {
	data, err := serializeInput(GlobalFieldRequest{GlobalField: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.post(
		ctx,
		"/v3/global_fields/",
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &GlobalFieldResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.GlobalField, nil
}

func (si *StackInstance) GlobalFieldUpdate(ctx context.Context, uid string, input GlobalFieldInput) (*GlobalField, error) {
	data, err := serializeInput(GlobalFieldRequest{GlobalField: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.put(
		ctx,
		fmt.Sprintf("/v3/global_fields/%s", uid),
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &GlobalFieldResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.GlobalField, nil
}

func (si *StackInstance) GlobalFieldDelete(ctx context.Context, uid string) error {
	resp, err := si.client.delete(
		ctx,
		fmt.Sprintf("/v3/global_fields/%s", uid),
		url.Values{},
		si.headers(),
		nil,
	)

	if err != nil {
		return err
	}

	result := &GlobalFieldResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return err
	}

	return nil
}

func (si *StackInstance) GlobalFieldFetch(ctx context.Context, uid string) (*GlobalField, error) {
	resp, err := si.client.get(
		ctx,
		fmt.Sprintf("/v3/global_fields/%s", uid),
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := &GlobalFieldResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.GlobalField, nil
}

func (si *StackInstance) GlobalFieldFetchAll(ctx context.Context) ([]GlobalField, error) {
	resp, err := si.client.get(
		ctx,
		"/v3/global_fields",
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := struct {
		GlobalFields []GlobalField `json:"global_fields"`
	}{}

	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.GlobalFields, nil
}
