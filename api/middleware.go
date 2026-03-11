package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Long4Changes/MySimpleBank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	// assume our system only allow one type of authorization
	authorizationTypeBearer = "bearer"
	// we can get the payload from context with this key easily
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	// this anonymous function is actually the middleware 
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// the format: Bearer + "<space>" + access token
		// Bearer prefix is to let the server know the type of authorization
		// strings.Fields() can split authorization header by space
		fields := strings.Fields(authorizationHeader)
		// we expect fields have two parts
		if len(fields) < 2 {
			err := errors.New("invalid authorizaiton header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// the first element of fields is the type
		// lowercase helps comparison
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// the access token is the second element of fileds slice
		accessToken := fields[1]
		// verify & get the payload
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// store the payload to the context
		ctx.Set(authorizationPayloadKey, payload)
		// forward request to the next handler
		ctx.Next()
	}
}
