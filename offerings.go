package revenuecat

import "fmt"

type Offering struct {
	CreatedAt   int64  `json:"created_at"`
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	IsCurrent   bool   `json:"is_current"`
	LookupKey   string `json:"lookup_key"`
	Metadata    any    `json:"metadata"`
	Object      string `json:"object"`
	ProjectID   string `json:"project_id"`
}

func (c *Client) ListAllProjectOfferings(projectId string) (RVCPageResp[Offering], error) {
	var resp RVCPageResp[Offering]
	err := c.call("GET", fmt.Sprintf("projects/%s/offerings", projectId), 2, nil, "", &resp)
	return resp, err
}

func (c *Client) GetProjectOffering(projectId, offeringId string) (RVCPageResp[Offering], error) {
	var resp RVCPageResp[Offering]
	err := c.call("GET", fmt.Sprintf("projects/%s/offerings/%s", projectId, offeringId), 2, nil, "", &resp)
	return resp, err
}
