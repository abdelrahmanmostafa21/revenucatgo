package revenuecat

type Project struct {
	CreatedAt int64  `json:"created_at"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Object    string `json:"object"`
}

func (c *Client) ListAllProjects() (RVCPageResp[Project], error) {
	var resp RVCPageResp[Project]
	err := c.call("GET", "projects", 2, nil, "", &resp)
	return resp, err
}
