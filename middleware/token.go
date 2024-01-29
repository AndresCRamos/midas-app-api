package middleware

import (
	"context"
	"log"

	"firebase.google.com/go/v4/auth"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	"github.com/gin-gonic/gin"
)

func VerifyToken(auth *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(error_utils.EmptyToken{}.GetAPIError())
			c.Abort()
			return
		}

		idToken := authHeader[len("Bearer "):]
		token, err := auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Printf("Error verifying token %s: %v", idToken, err)
			c.JSON(error_utils.InvalidToken{}.GetAPIError())
			c.Abort()
			return
		}

		// Add user information to the context for later use in handlers.
		c.Set("user", token.UID)

		// Continue processing the request.
		c.Next()
	}

}
