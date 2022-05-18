package utils

import (
	"bytes"
	"encoding/json"
	"github.com/fliptable-io/subscription-service/src/models"
	"log"
	"net/http"
	"os"
	"strconv"
)

func StripeInvoiceEmail(inf models.InvoiceInfo) {
	notificationKey := os.Getenv("NOTIFICATION_API_KEY")

	postBody, _ := json.Marshal(map[string]string{
		"customerId": inf.CustomerId,
		"total":      strconv.Itoa(inf.Total),
		"action":     "stripe-invoice",
		"apiKey":     notificationKey,
	})

	responseBody := bytes.NewBuffer(postBody)

	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("https://notifications.fliptable.io/api/v1/connections/", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()
}
