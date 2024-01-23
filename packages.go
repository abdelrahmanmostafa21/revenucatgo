package revenuecat

type Package struct {
	ID                        string    `json:"id,omitempty"`
	Identifier                string    `json:"identifier"`
	PlatformProductIdentifier string    `json:"platform_product_identifier,omitempty"`
	DisplayName               string    `json:"display_name,omitempty"`
	OfferingID                string    `json:"offering_id,omitempty"`
	Store                     string    `json:"store,omitempty"`
	Products                  []Product `json:"products,omitempty"`
}

type Product struct {
	AppID           string `json:"app_id"`
	CreatedAt       int64  `json:"created_at"`
	DisplayName     any    `json:"display_name"`
	ID              string `json:"id"`
	Object          string `json:"object"`
	StoreIdentifier string `json:"store_identifier"`
	Subscription    struct {
		Duration            any `json:"duration"`
		GracePeriodDuration any `json:"grace_period_duration"`
		TrialDuration       any `json:"trial_duration"`
	} `json:"subscription"`
	Type string `json:"type"`
}

func (c *Client) ListAllProducts(projectId string) (RVCPageResp[Product], error) {
	var resp RVCPageResp[Product]
	err := c.call("GET", "projects/"+projectId+"/products", 2, nil, "", &resp)
	return resp, err
}
