package middleware

import (
	"github.com/gin-gonic/gin"

	"learninghub/constants"
	"learninghub/errors"
	"learninghub/utils"
)

// ProductValidationMiddleware validates product parameter and adds it to context
func ProductValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		product := c.Param(constants.ProductParamKey)

		if !utils.IsValidProduct(product) {
			errors.AbortWithError(c, errors.ErrInvalidProduct, "Invalid product parameter")
			return
		}

		// Add product to context for use in handlers
		c.Set(constants.ProductContextKey, product)
		c.Next()
	}
}

// GetProductFromContext extracts product from gin context
func GetProductFromContext(c *gin.Context) (string, bool) {
	value, exists := c.Get(constants.ProductContextKey)
	if !exists {
		return "", false
	}

	product, ok := value.(string)
	if !ok {
		return "", false
	}

	return product, true
}
