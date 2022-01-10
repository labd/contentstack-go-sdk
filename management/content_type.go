package management

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type ContentTypeResponse struct {
	ContentType ContentType `json:"content_type"`
}

type ContentTypeRequest struct {
	ContentType ContentTypeInput `json:"content_type"`
}

// ContentType represents the content type in contentstack.
type ContentType struct {
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	Title             string              `json:"title,omitempty"`
	UID               string              `json:"uid,omitempty"`
	Schema            json.RawMessage     `json:"schema"`
	Options           *ContentTypeOptions `json:"options"`
	MaintainRevisions bool                `json:"maintain_revisions"`
	Description       string              `json:"description"`
}

// ContentTypeInput is used to create or update a content type
type ContentTypeInput struct {
	Title       *string         `json:"title,omitempty"`
	UID         *string         `json:"uid,omitempty"`
	Description *string         `json:"description,omitempty"`
	Schema      json.RawMessage `json:"schema,omitempty"`
}

type ContentTypeOptions struct {
	Title       string   `json:"title"`
	Publishable bool     `json:"bool"`
	IsPage      bool     `json:"is_page"`
	Singleton   bool     `json:"singleton"`
	SubTitle    []string `json:"sub_title"`
	UrlPattern  string   `json:"url_pattern"`
	UrlPrefix   string   `json:"url_prefix"`
}

func (si *StackInstance) ContentTypeCreate(ctx context.Context, input ContentTypeInput) (*ContentType, error) {
	data, err := serializeInput(ContentTypeRequest{ContentType: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.post(
		ctx,
		"/v3/content_types/",
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &ContentTypeResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.ContentType, nil
}

func (si *StackInstance) ContentTypeUpdate(ctx context.Context, uid string, input ContentTypeInput) (*ContentType, error) {
	data, err := serializeInput(ContentTypeRequest{ContentType: input})
	if err != nil {
		return nil, err
	}

	resp, err := si.client.put(
		ctx,
		fmt.Sprintf("/v3/content_types/%s", uid),
		url.Values{},
		si.headers(),
		data,
	)
	if err != nil {
		return nil, err
	}

	result := &ContentTypeResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.ContentType, nil
}

func (si *StackInstance) ContentTypeDelete(ctx context.Context, uid string) error {
	resp, err := si.client.delete(
		ctx,
		fmt.Sprintf("/v3/content_types/%s", uid),
		url.Values{},
		si.headers(),
		nil,
	)

	if err != nil {
		return err
	}

	result := &ContentTypeResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return err
	}

	return nil
}

func (si *StackInstance) ContentTypeFetch(ctx context.Context, uid string) (*ContentType, error) {
	resp, err := si.client.get(
		ctx,
		fmt.Sprintf("/v3/content_types/%s", uid),
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := &ContentTypeResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.ContentType, nil
}

func (si *StackInstance) ContentTypeFetchAll(ctx context.Context) ([]ContentType, error) {
	resp, err := si.client.get(
		ctx,
		"/v3/content_types",
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := struct {
		ContentTypes []ContentType `json:"content_types"`
	}{}

	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.ContentTypes, nil
}
