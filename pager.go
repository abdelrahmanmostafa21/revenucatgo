package revenuecat

type RVCPager struct {
	NextPage any    `json:"next_page"`
	Object   string `json:"object"`
	URL      string `json:"url"`
}

type RVCPageResp[T any] struct {
	RVCPager
	Items []T `json:"items"`
}
