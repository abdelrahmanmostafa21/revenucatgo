package revenuecat

import (
	"encoding/json"
	"time"
)

type SubscriberResp struct {
	RequestDate   time.Time  `json:"request_date"`
	RequestDateMs int64      `json:"request_date_ms"`
	Subscriber    Subscriber `json:"subscriber"`
}

// Subscriber holds a subscriber returned by the RevenueCat API.
type Subscriber struct {
	Entitlements               map[string]Entitlement         `json:"entitlements"`
	FirstSeen                  time.Time                      `json:"first_seen"`
	LastSeen                   time.Time                      `json:"last_seen"`
	ManagementURL              string                         `json:"management_url"`
	OriginalAppUserID          string                         `json:"original_app_user_id"`
	OriginalApplicationVersion *string                        `json:"original_application_version"`
	OriginalPurchaseDate       time.Time                      `json:"original_purchase_date"`
	Subscriptions              map[string]Subscription        `json:"subscriptions"`
	NonSubscriptions           map[string][]NonSubscription   `json:"non_subscriptions"`
	SubscriberAttributes       map[string]SubscriberAttribute `json:"subscriber_attributes"`
}

type Entitlement struct {
	ExpiresDate            time.Time `json:"expires_date"`
	GracePeriodExpiresDate time.Time `json:"grace_period_expires_date"`
	PurchaseDate           time.Time `json:"purchase_date"`
	ProductIdentifier      string    `json:"product_identifier"`
}

type Subscription struct {
	ExpiresDate             *time.Time `json:"expires_date"`
	PurchaseDate            time.Time  `json:"purchase_date"`
	OriginalPurchaseDate    time.Time  `json:"original_purchase_date"`
	PeriodType              PeriodType `json:"period_type"`
	Store                   Store      `json:"store"`
	IsSandbox               bool       `json:"is_sandbox"`
	UnsubscribeDetectedAt   *time.Time `json:"unsubscribe_detected_at"`
	BillingIssuesDetectedAt *time.Time `json:"billing_issues_detected_at"`
}

type NonSubscription struct {
	ID                   string    `json:"id"`
	IsSandbox            bool      `json:"is_sandbox"`
	OriginalPurchaseDate time.Time `json:"original_purchase_date"`
	PurchaseDate         time.Time `json:"purchase_date"`
	Store                Store     `json:"store"`
	StoreTransactionID   string    `json:"store_transaction_id"`
}

type SubscriberAttribute struct {
	Value     string `json:"value"`
	UpdatedAt int64  `json:"updated_at_ms"`
}

// PeriodType holds the predefined values for a subscription period.
type PeriodType string

// https://docs.revenuecat.com/reference#the-subscription-object
const (
	NormalPeriodType PeriodType = "normal"
	TrialPeriodType  PeriodType = "trial"
	IntroPeriodType  PeriodType = "intro"
)

// Store holds the predefined values for a store.
type Store string

// https://docs.revenuecat.com/reference#the-subscription-object
const (
	AppStore         Store = "app_store"
	MacAppStore      Store = "mac_app_store"
	PlayStore        Store = "play_store"
	StripeStore      Store = "stripe"
	PromotionalStore Store = "promotional"
)

// IsEntitledTo returns true if the Subscriber has the given entitlement.
func (s Subscriber) IsEntitledTo(entitlement string) bool {
	e, ok := s.Entitlements[entitlement]
	if !ok {
		return false
	}
	return !e.ExpiresDate.Before(time.Now())
}

// GetSubscriber gets the latest subscriber info or creates one if it doesn't exist.
// https://docs.revenuecat.com/reference#subscribers
func (c *Client) GetSubscriber(userID string) (SubscriberResp, error) {
	return c.GetSubscriberWithPlatform(userID, "")
}

// GetSubscriberWithPlatform gets the latest subscriber info or creates one if it doesn't exist, updating the subscriber record's last_seen
// value for the platform provided.
// https://docs.revenuecat.com/reference#subscribers
func (c *Client) GetSubscriberWithPlatform(userID string, platform string) (SubscriberResp, error) {
	var resp SubscriberResp
	err := c.do("GET", "subscribers/"+userID, nil, platform, &resp, 1)
	return resp, err
}

// UpdateSubscriberAttributes updates subscriber attributes for a user.
// https://docs.revenuecat.com/reference#update-subscriber-attributes
func (c *Client) UpdateSubscriberAttributes(userID string, attributes map[string]SubscriberAttribute) error {
	req := struct {
		Attributes map[string]SubscriberAttribute `json:"attributes"`
	}{
		Attributes: attributes,
	}
	return c.call("POST", "subscribers/"+userID+"/attributes", 1, req, "", nil)
}

// DeleteSubscriber permanently deletes a subscriber.
// https://docs.revenuecat.com/reference#subscribersapp_user_id
func (c *Client) DeleteSubscriber(userID string) error {
	return c.call("DELETE", "subscribers/"+userID, 1, nil, "", nil)
}

func (attr *SubscriberAttribute) MarshalJSON() ([]byte, error) {
	var updatedAt int64
	if !fromMilliseconds(attr.UpdatedAt).IsZero() {
		updatedAt = attr.UpdatedAt
	}
	return json.Marshal(&struct {
		Value     string `json:"value"`
		UpdatedAt int64  `json:"updated_at_ms,omitempty"`
	}{
		Value:     attr.Value,
		UpdatedAt: updatedAt,
	})
}

func (attr *SubscriberAttribute) UnmarshalJSON(data []byte) error {
	var jsonAttr struct {
		Value     string `json:"value"`
		UpdatedAt int64  `json:"updated_at_ms,omitempty"`
	}
	if err := json.Unmarshal(data, &jsonAttr); err != nil {
		return err
	}
	attr.Value = jsonAttr.Value
	if jsonAttr.UpdatedAt > 0 {
		attr.UpdatedAt = jsonAttr.UpdatedAt
	}
	return nil
}
