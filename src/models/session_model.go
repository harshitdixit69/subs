package models

//type CheckoutSession struct {
//	Body CheckoutSessionData `json:"data"`
//}

type BusinessPortalSession struct {
	CustomerId string `db:"stripeCustomerId" json:"stripeCustomerId"`
	BusinessId string `db:"businessId" json:"businessId"`
	TestKey    bool   `json:"testKey"`
}

type CheckoutSession struct {
	PriceId         string `db:"priceId" json:"priceId"`
	BusinessId      string `db:"businessId" json:"businessId"`
	FreeTrialCodeId string `db:"freeTrialCodeId" json:"freeTrialCodeId"`
	TestKey         bool   `json:"testKey"`
}

type FreeTrial struct {
	FreeTrialCodeId string `db:"freeTrialCodeId"`
	FreeTrialCode   string `db:"freeTrialCode" json:"freeTrialCode"`
	BusinessId      string `db:"businessId" json:"businessId"`
}
