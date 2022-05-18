package models

import (
	"encoding/json"
)

type SubscriptionInfo struct {
	CustomerId           string `db:"customerId"           json:"customerId"`
	PriceId              string `db:"priceId"              json:"priceId"`
	ProductId            string `db:"productId"            json:"productId"`
	StripeCustomerId     string `db:"stripeCustomerId"     json:"customer"`
	CheckoutSessionId    string `db:"checkoutSessionId"    json:"id"   `
	SubscriptionId       string `db:"subscriptionId"       json:"subscriptionId"      `
	StripeSubscriptionId string `db:"stripeSubscriptionId" json:"subscription"     `
	BusinessId           string `db:"businessId"           json:"businessId"                  `
}

type InvoiceInfo struct {
	CustomerId           string               `db:"customerId"           json:"customerId"`
	StripeCustomerId     string               `db:"stripeCustomerId"     json:"customer"`
	StripeInvoiceId      string               `db:"stripeInvoiceId"      json:"id"     `
	Paid                 bool                 `db:"paid"                 json:"paid"                `
	Total                int                  `db:"total"                json:"total"    `
	SubscriptionId       string               `db:"subscriptionId"        json:"subscriptionId"      `
	StripeSubscriptionId string               `db:"stripeSubscriptionId" json:"subscription"     `
	BusinessId           string               `db:"id"                   json:"businessId"  `
	InvoiceLines         InvoiceLineDataArray `json:"lines"`
}

type SubscriptionInfoUpdate struct {
	Plan                 PlanObj `json:"plan"`
	StripeSubscriptionId string  `db:"stripeSubscriptionId" json:"subscription"`
	SubscriptionId       string  `db:"id"`
}

type PlanObj struct {
	StripePriceId string `db:"priceId" json:"id"`
}

type InvoiceLineDataArray struct {
	Data []InvoiceLineData `json:"data"`
}

type InvoiceLineData struct {
	Period Period `json:"period"`
}

type Period struct {
	EndDate json.Number `db:"end_ts" json:"end"`
}
