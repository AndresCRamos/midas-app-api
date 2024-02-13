package middleware

import "github.com/gin-gonic/gin"

func TestMiddleware(id string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user", id)
	}
}
