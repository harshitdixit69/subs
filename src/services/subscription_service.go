package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fliptable-io/subscription-service/src/utils"

	"github.com/fliptable-io/subscription-service/src/models"
	"github.com/fliptable-io/subscription-service/src/queries"
	"github.com/fliptable-io/subscription-service/src/repos"
	loggy "github.com/fliptable-io/subscription-service/src/utils/logging"
	"github.com/stripe/stripe-go/v72/webhook"
)

type SubscriptionService struct {
	Pg *repos.PostgresRepo
}

func (s SubscriptionService) PrintSubInfo(err string, inf models.SubscriptionInfo) {
	errString := ""

	errString += err + "\n" + "BusinessId: " + inf.BusinessId + "\n" + "PriceId: " + inf.PriceId +
		"\n" + "CustomerId: " + inf.CustomerId + "\n" + "StripeCustomerId: " + inf.StripeCustomerId +
		"\n" + "SubId: " + inf.SubscriptionId + "\n" + "StripeSubId: " + inf.StripeSubscriptionId +
		"\n" + "CheckoutId: " + inf.CheckoutSessionId + "\n" + "ProdId: " + inf.ProductId + "\n"

	loggy.Error(errString)
}

func (s SubscriptionService) PrintInvInfo(err string, inf models.InvoiceInfo) {
	errString := ""

	errString += err + "\n" + "BusinessId: " + inf.BusinessId + "\n" + "StripeInvId: " + inf.StripeInvoiceId +
		"\n" + "CustomerId: " + inf.CustomerId + "\n" + "StripeCustomerId: " + inf.StripeCustomerId +
		"\n" + "SubId: " + inf.SubscriptionId + "\n" + "StripeSubId: " + inf.StripeSubscriptionId +
		"\n" + "Total: " + strconv.Itoa(inf.Total) + "\n" + "Paid: " + strconv.FormatBool(inf.Paid) +
		"\n" + "EndDate: " + inf.InvoiceLines.Data[0].Period.EndDate.String() + "\n"

	loggy.Error(errString)
}

func (s SubscriptionService) CreateSubscription(inf models.SubscriptionInfo) {
	fmt.Println("inf===========", inf)
	err := s.Pg.Query(nil, queries.AddSubscription, &inf)
	if err != nil {
		s.PrintSubInfo("[SUB] ADD SUB ERROR", inf)
	}

	err = s.Pg.Query(&inf, queries.GetSubscriptionId, &inf)
	if err != nil {
		s.PrintSubInfo("[SUB] SUB ID ERROR", inf)
	}

	s.CreateCustomer(inf)
}

func (s SubscriptionService) CreateCustomer(inf models.SubscriptionInfo) {

	err := s.Pg.Query(&inf, queries.GetBusinessAndPriceId, &inf)
	if err != nil {
		s.PrintSubInfo("[CUS] BUSINESS ID ERROR: ", inf)
	}

	err = s.Pg.Query(&inf, queries.GetProductId, &inf)
	if err != nil {
		s.PrintSubInfo("[CUS] PRODUCT ID ERROR", inf)
	}

	err = s.Pg.Query(nil, queries.AddCustomer, &inf)
	if err != nil {
		s.PrintSubInfo("[CUS] CUSTOMER ADD ERROR", inf)
	}

	err = s.Pg.Query(nil, queries.RemoveCustomerInTransition, &inf)
	if err != nil {
		s.PrintSubInfo("[CUS] REMOVE CUSTOMER TRANSIT ERROR", inf)
	}
}

func (s SubscriptionService) CreateInvoice(invoice models.InvoiceInfo, paymentFailed bool) {

	err := s.Pg.Query(&invoice, queries.GetCustomerId, &invoice)
	if err != nil {
		s.PrintInvInfo("[INV] CUSTOMER ID RETRIEVE ERROR", invoice)
	}

	err = s.Pg.Query(&invoice, queries.GetSubscriptionId, &invoice)
	if err != nil {
		s.PrintInvInfo("[INV] SUBSCRIPTION ID RETRIEVE ERROR", invoice)
	}

	err = s.Pg.Query(nil, queries.AddInvoice, &invoice)
	if err != nil {
		s.PrintInvInfo("[INV] INV ADD ERROR", invoice)
	}

	s.UpdateSubscriptionEndData(invoice)

	//if paymentFailed {
	//	s.UpdateSubscriptionEndData(invoice)
	//} else {
	//	s.DisableSubscription(invoice)
	//}

	utils.StripeInvoiceEmail(invoice)
}

func (s SubscriptionService) UpdateSubscriptionEndData(invoice models.InvoiceInfo) error {

	type endTmp struct {
		StripeSubscriptionId string    `db:"stripeSubscriptionId"`
		EndDate              time.Time `db:"end_ts"`
	}
	timeInt64, _ := invoice.InvoiceLines.Data[0].Period.EndDate.Int64()
	unixTime := time.Unix(timeInt64, 0)

	endDate := endTmp{
		EndDate:              unixTime,
		StripeSubscriptionId: invoice.StripeSubscriptionId,
	}

	err := s.Pg.Query(nil, queries.UpdateSubscriptionEndDate, &endDate)
	if err != nil {
		loggy.Error(err)
		return err
	}
	return nil
}

func (s SubscriptionService) DisableSubscription(invoice models.InvoiceInfo) {

	err := s.Pg.Query(&invoice, queries.DisableSubscription, &invoice)
	if err != nil {
		s.PrintInvInfo("[INV] DISABLE SUB ERR", invoice)
	}

}

func (s SubscriptionService) SubscriptionChanged(subscriptionInfo models.SubscriptionInfoUpdate) {
	err := s.Pg.Query(&subscriptionInfo, queries.GetSubscriptionId, &subscriptionInfo)
	if err != nil {
		loggy.Error("No Subscription ID")
	}

	err = s.Pg.Query(nil, queries.SubscriptionChanged, &subscriptionInfo)
	if err != nil {
		loggy.Error("Error updating subscription")
	}
}

func (s SubscriptionService) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		loggy.Error("ioutil.ReadAll: %v", err)
		return
	}

	webhookSecret, exists := os.LookupEnv("STRIPE_WEBHOOK_SECRET")
	if !exists {
		http.Error(w, err.Error(), http.StatusBadRequest)
		loggy.Error(err)
		return
	}

	event, err := webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		loggy.Error("webhook.ConstructEvent: %v", err)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		inf := models.SubscriptionInfo{}
		err = json.Unmarshal(event.Data.Raw, &inf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.CreateSubscription(inf)

	case "invoice.paid":
		inf := models.InvoiceInfo{}
		err = json.Unmarshal(event.Data.Raw, &inf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.CreateInvoice(inf, false)

	case "invoice.payment_failed":
		inf := models.InvoiceInfo{}
		err = json.Unmarshal(event.Data.Raw, &inf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.CreateInvoice(inf, true)

	case "customer.subscription.updated":
		inf := models.SubscriptionInfoUpdate{}
		err = json.Unmarshal(event.Data.Raw, &inf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		s.SubscriptionChanged(inf)

	default:
		// unhandled event type
	}

	w.WriteHeader(http.StatusOK)
}

func (s SubscriptionService) TestHandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		loggy.Error("ioutil.ReadAll: %v", err)
		return
	}
	webhookSecret, exists := os.LookupEnv("STRIPE_WEBHOOK_TEST_SECRET")
	if !exists {
		http.Error(w, err.Error(), http.StatusBadRequest)
		loggy.Error(err)
		return
	}

	event, err := webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		loggy.Error("webhook.ConstructEvent: %v", err)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		inf := models.SubscriptionInfo{}
		err = json.Unmarshal(event.Data.Raw, &inf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		s.CreateSubscription(inf)

	case "invoice.paid":
		inf := models.InvoiceInfo{}
		err = json.Unmarshal(event.Data.Raw, &inf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.CreateInvoice(inf, false)

	case "invoice.payment_failed":
		inf := models.InvoiceInfo{}
		err = json.Unmarshal(event.Data.Raw, &inf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.CreateInvoice(inf, true)

	case "customer.subscription.updated":
		inf := models.SubscriptionInfoUpdate{}
		err = json.Unmarshal(event.Data.Raw, &inf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		s.SubscriptionChanged(inf)

	default:
		// unhandled event type
	}

	w.WriteHeader(http.StatusOK)
}
