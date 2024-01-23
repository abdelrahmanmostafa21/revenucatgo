package revenuecat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client makes authorized calls to the RevenueCat API.
type Client struct {
	apiKey    string
	publicKey string
	cookie    string
	apiURL    string
	http      doer
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// New returns a new *Client for the provided API key.
// For more information on authentication, see https://docs.revenuecat.com/docs/authentication.
func New(apiKey, publicKey, cookie string, client *http.Client) *Client {
	if client == nil {
		client = &http.Client{
			// Set a long timeout here since calls to Apple are probably invloved.
			Timeout: 10 * time.Second,
		}
	}

	return &Client{
		apiKey:    apiKey,
		publicKey: publicKey,
		apiURL:    "https://api.revenuecat.com/v1/",
		http:      client,
		cookie:    cookie,
	}
}

func (c *Client) do(method, path string, reqBody interface{}, platform string, respBody interface{}, public bool) error {
	var reqBodyJSON io.Reader
	if reqBody != nil {
		js, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %v", err)
		}
		reqBodyJSON = bytes.NewBuffer(js)
	}
	req, err := http.NewRequest(method, c.apiURL+path, reqBodyJSON)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	if strings.Contains(path, "developers/me/") {
		req.Header.Add("Cookie", c.cookie)
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
	} else {
		auth := "Bearer " + c.apiKey
		if public {
			auth = "Bearer " + c.publicKey
		}
		req.Header.Add("Authorization", auth)
	}
	req.Header.Add("Content-Type", "application/json")
	if platform != "" {
		req.Header.Add("X-Platform", platform)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var errResp Error
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return err
		}
		return errResp
	}
	if respBody == nil {
		// Expecting an empty body.
		return nil
	}
	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}
	return nil
}

func (c *Client) call(method, path string, reqBody interface{}, platform string, respBody interface{}) error {
	return c.do(method, path, reqBody, platform, respBody, false)
}
