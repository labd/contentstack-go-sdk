package management

import (
	"context"
	"net/url"
)

type ContentType struct {
}

func (si *StackInstance) ContentTypeCreate() {

}

func (si *StackInstance) ContentTypeUpdate() {

}

func (si *StackInstance) ContentTypeDelete() {

}

func (si *StackInstance) ContentTypeFetch() {

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

	err = si.client.processResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	return result.ContentTypes, nil

}
