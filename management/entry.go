package management

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type EntryResponse struct {
	Entry json.RawMessage `json:"entry"`
}

type EntryRequest struct {
	Entry json.RawMessage `json:"entry"`
}

// Entry represents the content type in contentstack.
type Entry struct {
	UID       string    `json:"uid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	Locale    string    `json:"locale"`
	Version   int       `json:"_version"`

	Fields map[string]interface{} `json:"-"`
}

func (e *EntryResponse) deserialize() (*Entry, error) {
	return deserializeEntry(e.Entry)
}

// EntryInput is used to create or update a entry
type EntryInput struct {
	ContentTypeUID string `json:"-"`
	Locale         string `json:"-"`
	Fields         map[string]interface{}
}

func (e *EntryInput) serialize() (json.RawMessage, url.Values, error) {
	data, err := json.MarshalIndent(e.Fields, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	params := url.Values{
		"locale": []string{e.Locale},
	}

	return json.RawMessage(data), params, nil
}

// EntryDeleteInput is used to delete an entry
type EntryContextInput struct {
	ContentTypeUID string
	Locale         string
	UID            string
}

func (si *StackInstance) EntryCreate(ctx context.Context, input *EntryInput) (*Entry, error) {
	data, params, err := input.serialize()
	if err != nil {
		return nil, err
	}

	body, err := json.MarshalIndent(&EntryRequest{
		Entry: data,
	}, "", "  ")
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/v3/content_types/%s/entries", input.ContentTypeUID)
	resp, err := si.client.post(
		ctx,
		endpoint,
		params,
		si.headers(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	result := EntryResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.deserialize()
}

func (si *StackInstance) EntryUpdate(ctx context.Context, uid string, input *EntryInput) (*Entry, error) {
	data, params, err := input.serialize()
	if err != nil {
		return nil, err
	}

	body, err := json.MarshalIndent(&EntryRequest{
		Entry: data,
	}, "", "  ")
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/v3/content_types/%s/entries/%s", input.ContentTypeUID, uid)
	resp, err := si.client.put(
		ctx,
		endpoint,
		params,
		si.headers(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	result := EntryResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.deserialize()
}

func (si *StackInstance) EntryDelete(ctx context.Context, input *EntryContextInput) error {
	endpoint := fmt.Sprintf("/v3/content_types/%s/entries/%s", input.ContentTypeUID, input.UID)
	resp, err := si.client.delete(
		ctx,
		endpoint,
		url.Values{},
		si.headers(),
		nil,
	)

	if err != nil {
		return err
	}

	result := &EntryResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return err
	}

	return nil
}

func (si *StackInstance) EntryFetch(ctx context.Context, input *EntryContextInput) (*Entry, error) {
	endpoint := fmt.Sprintf("/v3/content_types/%s/entries/%s", input.ContentTypeUID, input.UID)
	resp, err := si.client.get(
		ctx,
		endpoint,
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	result := &EntryResponse{}
	if err = si.client.processResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.deserialize()
}

func (si *StackInstance) EntryFetchAll(ctx context.Context, contentTypeUID string) ([]Entry, error) {
	endpoint := fmt.Sprintf("/v3/content_types/%s/entries", contentTypeUID)
	resp, err := si.client.get(
		ctx,
		endpoint,
		url.Values{},
		si.headers(),
	)
	if err != nil {
		return nil, err
	}

	response := struct {
		Entries []json.RawMessage `json:"entries"`
	}{}
	if err = si.client.processResponse(resp, &response); err != nil {
		return nil, err
	}

	result := make([]Entry, len(response.Entries))
	for i := range response.Entries {
		entry, err := deserializeEntry(response.Entries[i])
		if err != nil {
			return nil, err
		}
		result[i] = *entry
	}

	return result, nil
}

func deserializeEntry(data json.RawMessage) (*Entry, error) {
	result := &Entry{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}

	// Unmarshall again for the Fields
	err = json.Unmarshal(data, &result.Fields)
	if err != nil {
		return nil, err
	}

	// Delete internal fields
	known_fields := []string{"tags", "locale", "uid", "created_by", "updated_by", "created_at", "updated_at", "ACL", "_version", "_in_progress", "publish_details"}
	for _, field := range known_fields {
		delete(result.Fields, field)
	}
	return result, nil

}
