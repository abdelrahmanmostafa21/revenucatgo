package revenuecat

import (
	"fmt"
	"net/http"
)

type Offering struct {
	ID          string    `json:"id,omitempty"`
	Description string    `json:"description,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	Identifier  string    `json:"identifier"`
	Packages    []Package `json:"packages,omitempty"`
}

// OverrideOffering overrides the current Offering for a specific user.
// https://docs.revenuecat.com/reference#override-offering
func (c *Client) OverrideOffering(userID string, offeringUUID string) (Subscriber, error) {
	var resp struct {
		Subscriber Subscriber `json:"subscriber"`
	}
	err := c.call("POST", "subscribers/"+userID+"/offerings/"+offeringUUID+"/override", nil, "", &resp)
	return resp.Subscriber, err
}

// DeleteOfferingOverride reset the offering overrides back to the current offering for a specific user.
// https://docs.revenuecat.com/reference#delete-offering-override
func (c *Client) DeleteOfferingOverride(userID string) (Subscriber, error) {
	var resp struct {
		Subscriber Subscriber `json:"subscriber"`
	}
	err := c.call("DELETE", "subscribers/"+userID+"/offerings/override", nil, "", &resp)
	return resp.Subscriber, err
}

// GetOfferings gets the offerings and current offering ID for your app for a specific user
// https://docs.revenuecat.com/reference#get-offerings
func (c *Client) GetOfferings(userID, platform string) ([]Offering, string, error) {
	resp := struct {
		CurrentOfferingId string     `json:"current_offering_id"`
		Offerings         []Offering `json:"offerings"`
	}{}

	err := c.do(http.MethodGet, fmt.Sprintf("subscribers/%s/offerings", userID), nil, platform, &resp, true)
	return resp.Offerings, resp.CurrentOfferingId, err
}

func (c *Client) GetAllOfferings(appID string) ([]Offering, error) {
	var resp []Offering
	err := c.call("GET", fmt.Sprintf("developers/me/apps/%s/new_offerings", appID), nil, "", &resp)
	return resp, err
}

func (c *Client) CreateOffering(appID string, offering *Offering) (Offering, error) {
	var resp Offering
	err := c.call("POST", fmt.Sprintf("developers/me/apps/%s/new_offerings", appID), offering, "", &resp)
	return resp, err
}
