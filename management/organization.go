package management

type OrganizationClient struct {
	organization_uid string
	client           *Client
}

func (c *Client) Organization(ouid string) *OrganizationClient {
	return &OrganizationClient{
		client:           c,
		organization_uid: ouid,
	}
}
