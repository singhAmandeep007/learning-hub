package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"learning-hub/constants"
	"learning-hub/models"
	"learning-hub/utils"
)

// ProductValidationMiddleware validates product parameter and adds it to context
func ProductValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		product := c.Param(constants.ProductParamKey)

		if !utils.IsValidProduct(product) {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   constants.InvalidParam,
				Message: "Invalid product parameter",
			})
			return
		}

		// Add product to context for use in handlers
		c.Set(constants.ProductContextKey, product)
		c.Next()
	}
}

// GetProductFromContext extracts product from gin context
func GetProductFromContext(c *gin.Context) string {
	product := c.MustGet(constants.ProductContextKey)
	return product.(string)
}
