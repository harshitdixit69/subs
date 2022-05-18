package middleware

import (
	"github.com/fliptable-io/subscription-service/src/queries"
	"github.com/fliptable-io/subscription-service/src/repos"
	"github.com/fliptable-io/subscription-service/src/utils/errors"
	loggy "github.com/fliptable-io/subscription-service/src/utils/logging"
	"github.com/gin-gonic/gin"
)

func SetBasicHeaders(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.Header().Set("Authorization", "")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}

func PanicRecovery(c *gin.Context, err interface{}) {
	loggy.Error(err)
	errors.UnknownError.Respond(c)
}

func VerifyNotificationAccess(apiKey string, Pg *repos.PostgresRepo) bool {
	type auth struct {
		Auth bool `db:"auth"`
	}

	type key struct {
		APIKey string `db:"apiKey"`
	}

	if apiKey == "" {
		return false
	}

	res := new(auth)
	params := new(key)

	params.APIKey = apiKey

	err := Pg.Query(res, queries.MatchingKey, &params)
	if err != nil {
		loggy.Error(err)
	}

	return res.Auth
}
