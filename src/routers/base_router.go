package routers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/fliptable-io/subscription-service/src/models"

	"github.com/fliptable-io/subscription-service/src/data"
	"github.com/fliptable-io/subscription-service/src/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type BaseRouter struct {
	*gin.Engine
}

func NewBaseRouter(subRouter SubscriptionRouter) *BaseRouter {
	router := &BaseRouter{Engine: gin.Default()}

	// Disable Console Color
	//gin.DisableConsoleColor()
	//gin.SetMode(gin.ReleaseMode)

	router.Use(gin.CustomRecovery(middleware.PanicRecovery))
	router.Use(CORS())

	v1 := router.Group("api/v1")

	AddSessionRoutes(v1)
	subRouter.AddSubscriptionRoutes(v1)

	v1.GET("/ping", pingTest)

	router.GET("/", route404)

	return router
}

func route404(c *gin.Context) {
	ind := rand.Intn(len(data.Facts))
	randomFact := data.Facts[ind]
	c.String(404, "404 Page Not Found, but here's a random fact: "+randomFact)
}

func pingTest(c *gin.Context) {
	c.String(http.StatusOK, "pong <=========> "+time.Now().String())
}

func ModifiedCORSConfig() gin.HandlerFunc {
	config := ModifiedConfig()
	config.AllowAllOrigins = true
	return cors.New(cors.Config(config))
}

func ModifiedConfig() models.Config {
	return models.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}
