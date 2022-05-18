package services

import (
	"fmt"
	"os"

	"github.com/fliptable-io/subscription-service/src/models"
	"github.com/fliptable-io/subscription-service/src/queries"
	"github.com/fliptable-io/subscription-service/src/repos"
	"github.com/stripe/stripe-go/v72"
	billingSession "github.com/stripe/stripe-go/v72/billingportal/session"
	checkoutSession "github.com/stripe/stripe-go/v72/checkout/session"
)

const baseUrl = "https://business.fliptable.io/locations/"

func CreateCheckoutSession(priceId string, businessId string, testKey bool) (*stripe.CheckoutSession, error) {
	pg := repos.NewPostgresRepo()

	key, exists := os.LookupEnv("STRIPE_LIVE_KEY")
	if exists {
		stripe.Key = key
	} else {
		fmt.Println("Error")
	}
	if testKey == true {
		key, exists = os.LookupEnv("STRIPE_TEST_KEY")
		if exists {
			stripe.Key = key
		} else {
			fmt.Println("Error")
		}
	}

	successUrl := baseUrl + businessId
	cancelUrl := baseUrl + businessId

	params := &stripe.CheckoutSessionParams{
		SuccessURL: &successUrl,
		CancelURL:  &cancelUrl,
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price: stripe.String(priceId),
				// For metered billing, do not pass quantity
				Quantity: stripe.Int64(1),
			},
		},
	}

	sess, _ := checkoutSession.New(params)

	type customer_trans struct {
		BusinessId        string `db:"businessId"`
		PriceId           string `db:"priceId"`
		CheckoutSessionId string `db:"checkoutSessionId"`
	}

	ct := customer_trans{
		BusinessId:        businessId,
		PriceId:           priceId,
		CheckoutSessionId: sess.ID,
	}
	err := pg.Query(nil, queries.AddCustomerInTransit, &ct)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func CreateBillingPortalSession(inf models.BusinessPortalSession) (*stripe.BillingPortalSession, error) {
	pg := repos.NewPostgresRepo()

	key, exists := os.LookupEnv("STRIPE_LIVE_KEY")
	if exists {
		stripe.Key = key
	} else {
		fmt.Println("Error")
	}

	if inf.TestKey == true {
		key, exists = os.LookupEnv("STRIPE_TEST_KEY")
		if exists {
			stripe.Key = key
		} else {
			fmt.Println("Error")
		}
	}

	err := pg.Query(&inf, queries.GetStripeCustomerId, &inf)
	if err != nil {
		return nil, err
	}
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(inf.CustomerId),
		ReturnURL: stripe.String(baseUrl + inf.BusinessId),
	}
	fmt.Println("param", params)
	sess, _ := billingSession.New(params)

	return sess, nil
}

func GetFreeTrialCode(freeTrial models.FreeTrial) (*models.FreeTrial, error) {
	pg := repos.NewPostgresRepo()

	type tempStruct struct {
		Exists bool `db:"exists"`
	}

	exists := tempStruct{}

	err := pg.Query(&exists, queries.CheckFreeTrialCode, &freeTrial)
	if err != nil {
		return nil, err
	}

	if exists.Exists {
		err = pg.Query(&freeTrial, queries.GetFreeTrialCode, &freeTrial)
		if err != nil {
			return nil, err
		}
	}

	return &freeTrial, nil
}

func CreateFreeTrialCheckoutSession(priceId string, businessId string, testKey bool, freeTrialCodeId string) (*stripe.CheckoutSession, error) {
	pg := repos.NewPostgresRepo()

	key, exists := os.LookupEnv("STRIPE_LIVE_KEY")
	if exists {
		stripe.Key = key
	} else {
		fmt.Println("Error")
	}

	if testKey == true {
		key, exists = os.LookupEnv("STRIPE_TEST_KEY")
		if exists {
			stripe.Key = key
		} else {
			fmt.Println("Error")
		}
	}

	successUrl := baseUrl + businessId
	cancelUrl := baseUrl + businessId
	trialPeriodDays := int64(14)

	params := &stripe.CheckoutSessionParams{
		SuccessURL: &successUrl,
		CancelURL:  &cancelUrl,
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price: stripe.String(priceId),
				// For metered billing, do not pass quantity
				Quantity: stripe.Int64(1),
			},
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			TrialPeriodDays: &trialPeriodDays,
		},
	}

	sess, _ := checkoutSession.New(params)

	type customer_trans struct {
		BusinessId        string `db:"businessId"`
		PriceId           string `db:"priceId"`
		CheckoutSessionId string `db:"checkoutSessionId"`
	}

	ct := customer_trans{
		BusinessId:        businessId,
		PriceId:           priceId,
		CheckoutSessionId: sess.ID,
	}

	err := pg.Query(nil, queries.AddCustomerInTransit, &ct)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
