package routers

import (
	"github.com/fliptable-io/subscription-service/src/services"
	"github.com/gin-gonic/gin"
)

type SubscriptionRouter struct {
	SubscriptionService services.SubscriptionService
}

func (sr SubscriptionRouter) AddSubscriptionRoutes(group *gin.RouterGroup) {
	router := group.Group("/subscription")

	router.POST("/webhook", sr.HandleWebhooks)
	router.POST("/test-webhook", sr.TestHandleWebhooks)
}

func (sr SubscriptionRouter) HandleWebhooks(c *gin.Context) {
	sr.SubscriptionService.HandleWebhook(c.Writer, c.Request)
}

func (sr SubscriptionRouter) TestHandleWebhooks(c *gin.Context) {
	sr.SubscriptionService.TestHandleWebhook(c.Writer, c.Request)
}
