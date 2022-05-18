package routers

import (
	"encoding/json"
	"fmt"

	"github.com/fliptable-io/subscription-service/src/models"
	"github.com/fliptable-io/subscription-service/src/services"
	loggy "github.com/fliptable-io/subscription-service/src/utils/logging"
	"github.com/gin-gonic/gin"
)

func AddSessionRoutes(group *gin.RouterGroup) {
	router := group.Group("/session")

	router.POST("/checkout", CreateCheckoutSession)
	router.POST("/portal", CreateBillingPortalSession)
	router.POST("/freeTrialCheck", GetFreeTrialCode)
	router.POST("/freeTrialCheckout", CreateFreeTrialSession)
}

func CreateCheckoutSession(c *gin.Context) {
	requestBody := models.CheckoutSession{TestKey: true}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return
	}
	retVal, err := services.CreateCheckoutSession(requestBody.PriceId, requestBody.BusinessId, requestBody.TestKey)
	if err != nil {
		loggy.Error(nil)
		c.JSON(400, nil)
		return
	}
	c.JSON(200, retVal)
}

func CreateBillingPortalSession(c *gin.Context) {
	requestBody := models.BusinessPortalSession{TestKey: true}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println(err)
	}
	retVal, err := services.CreateBillingPortalSession(requestBody)
	if err != nil {
		c.Status(400)
		c.Done()
		return
	}
	fmt.Println("=retVal==========", retVal)
	// price_1JzMRzDfMZZgzqA4EZQ59qWD
	//price_1K4rqZDfMZZgzqA4gDnGFxdn

	c.JSON(200, retVal)
}

func GetFreeTrialCode(c *gin.Context) {

	requestBody := models.FreeTrial{}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println(err)
	}

	freeTrial, err := services.GetFreeTrialCode(requestBody)
	if err != nil {
		c.Status(400)
		c.Done()
		return
	}

	c.JSON(200, freeTrial)
}
func CreateFreeTrialSession(c *gin.Context) {

	requestBody := models.CheckoutSession{TestKey: false}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println(err)
	}

	retVal, err := services.CreateFreeTrialCheckoutSession(requestBody.PriceId, requestBody.BusinessId, requestBody.TestKey, requestBody.FreeTrialCodeId)
	if err != nil {
		c.Status(400)
		c.Done()
		return
	}

	c.JSON(200, retVal)
}

//862de255-38d7-4751-ae07-faeb15326ab2

// price_1KX5PzDfMZZgzqA4TSXuowxe

// [default]
// aws_access_key_id = AKIAZEX5NKSU2EFMDE6J
// aws_secret_access_key = 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
