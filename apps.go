package revenuecat

type App struct {
	AppStore struct {
		BundleID string `json:"bundle_id"`
	} `json:"app_store,omitempty"`
	CreatedAt int64  `json:"created_at"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Object    string `json:"object"`
	ProjectID string `json:"project_id"`
	Type      string `json:"type"`
	PlayStore struct {
		PackageName string `json:"package_name"`
	} `json:"play_store,omitempty"`
}

func (c *Client) ListAllProjectApps(projectId string) (RVCPageResp[Project], error) {
	var resp RVCPageResp[Project]
	err := c.call("GET", "projects/"+projectId+"/apps", 2, nil, "", &resp)
	return resp, err
}

func (c *Client) GetProjectApp(projectId, appId string) (RVCPageResp[Project], error) {
	var resp RVCPageResp[Project]
	err := c.call("GET", "projects/"+projectId+"/apps/"+appId, 2, nil, "", &resp)
	return resp, err
}
