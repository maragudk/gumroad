// Package gumroad provides an API client for the Gumroad API.
// See https://app.gumroad.com/api
package gumroad

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client for Gumroad.
type Client struct {
	accessToken string
	endpoint    string
	httpClient  *http.Client
}

// NewClientOptions for NewClientWithOptions.
type NewClientOptions struct {
	AccessToken string
	Endpoint    string
	HTTPClient  *http.Client
}

// NewClient with default options.
func NewClient() *Client {
	return NewClientWithOptions(NewClientOptions{})
}

func NewClientWithOptions(opts NewClientOptions) *Client {
	if opts.HTTPClient == nil {
		opts.HTTPClient = &http.Client{Timeout: 3 * time.Second}
	}
	if opts.Endpoint == "" {
		opts.Endpoint = "https://api.gumroad.com"
	}
	opts.Endpoint = strings.TrimSuffix(opts.Endpoint, "/")
	return &Client{
		accessToken: opts.AccessToken,
		endpoint:    opts.Endpoint,
		httpClient:  opts.HTTPClient,
	}
}

type BaseResponse struct {
	Success bool `json:"success"`
}

type GetProductsResponse struct {
	BaseResponse
}

func (c *Client) GetProducts(ctx context.Context) (*GetProductsResponse, error) {
	var r GetProductsResponse
	if err := c.get(ctx, "/products", nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

type GetResourceSubscriptionsResponse struct {
	BaseResponse
	ResourceSubscriptions []struct {
		ID           string `json:"id"`
		ResourceName string `json:"resource_name"`
		PostURL      string `json:"post_url"`
	} `json:"resource_subscriptions"`
}

type ResourceSubscription string

var ResourceSubscriptions = []ResourceSubscription{
	ResourceSubscriptionSale,
	ResourceSubscriptionRefund,
	ResourceSubscriptionDispute,
	ResourceSubscriptionDisputeWon,
	ResourceSubscriptionCancellation,
	ResourceSubscriptionSubscriptionUpdated,
	ResourceSubscriptionSubscriptionEnded,
	ResourceSubscriptionSubscriptionRestarted,
}

const (
	ResourceSubscriptionSale                  = ResourceSubscription("sale")
	ResourceSubscriptionRefund                = ResourceSubscription("refund")
	ResourceSubscriptionDispute               = ResourceSubscription("dispute")
	ResourceSubscriptionDisputeWon            = ResourceSubscription("dispute_won")
	ResourceSubscriptionCancellation          = ResourceSubscription("cancellation")
	ResourceSubscriptionSubscriptionUpdated   = ResourceSubscription("subscription_updated")
	ResourceSubscriptionSubscriptionEnded     = ResourceSubscription("subscription_ended")
	ResourceSubscriptionSubscriptionRestarted = ResourceSubscription("subscription_restarted")
)

func (c *Client) GetResourceSubscriptions(ctx context.Context, name ResourceSubscription) (*GetResourceSubscriptionsResponse, error) {
	var r GetResourceSubscriptionsResponse
	var found bool
	for _, v := range ResourceSubscriptions {
		found = found || v == name
	}
	if !found {
		return nil, fmt.Errorf("name must be one of %v", ResourceSubscriptions)
	}
	if err := c.get(ctx, "/resource_subscriptions", map[string]string{"resource_name": string(name)}, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// PingRequest sent from Gumroad to your API.
type PingRequest struct {
	SaleID                     string            `json:"sale_id"`
	SaleTimestamp              string            `json:"sale_timestamp"`
	OrderNumber                string            `json:"order_number"`
	SellerID                   string            `json:"seller_id"`
	ProductID                  string            `json:"product_id"`
	ProductPermalink           string            `json:"product_permalink"`
	ShortProductID             string            `json:"short_product_id"`
	Email                      string            `json:"email"`
	URLParams                  map[string]string `json:"url_params"`
	FullName                   string            `json:"full_name"`
	PurchaserID                string            `json:"purchaser_id"`
	SubscriptionID             string            `json:"subscription_id"`
	IPCountry                  string            `json:"ip_country"`
	Price                      int               `json:"price"`
	Recurrence                 string            `json:"recurrence"`
	Variants                   map[string]string `json:"variants"`
	OfferCode                  string            `json:"offer_code"`
	Test                       bool              `json:"test"`
	CustomFields               map[string]string `json:"custom_fields"`
	ShippingInformation        map[string]string `json:"shipping_information"`
	IsRecurringCharge          bool              `json:"is_recurring_charge"`
	IsPreorderAuthorization    bool              `json:"is_preorder_authorization"`
	LicenseKey                 string            `json:"license_key"`
	Quantity                   int               `json:"quantity"`
	ShippingRate               int               `json:"shipping_rate"`
	Affiliate                  string            `json:"affiliate"`
	AffiliateCreditAmountCents int               `json:"affiliate_credit_amount_cents"`
	IsGiftReceiverPurchase     bool              `json:"is_gift_receiver_purchase"`
	GifterEmail                string            `json:"gifter_email"`
	GiftPrice                  int               `json:"gift_price"`
	Refunded                   bool              `json:"refunded"`
	DiscoverFeeCharge          bool              `json:"discover_fee_charge"`
	CanContact                 bool              `json:"can_contact"`
	Referrer                   string            `json:"referrer"`
	GumroadFee                 int               `json:"gumroad_fee"`
	Card                       string            `json:"card"`
}

func (c *Client) get(ctx context.Context, path string, args map[string]string, r interface{}) error {
	values := url.Values{}
	for k, v := range args {
		values.Set(k, v)
	}
	values.Set("access_token", c.accessToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/v2"+path,
		strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return fmt.Errorf("error constructing request: %w", err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error requesting: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode > 299 {
		return fmt.Errorf("error requesting, got status code %v", res.StatusCode)
	}
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return fmt.Errorf("error decoding response body: %w", err)
	}
	return nil
}
