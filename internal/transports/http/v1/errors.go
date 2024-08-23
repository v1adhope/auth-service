package httpv1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/auth-service/internal/models"
)

func setBindError(c *gin.Context, err error) {
	c.Error(err).SetType(gin.ErrorTypeBind)
}

func setAnyError(c *gin.Context, err error) {
	c.Error(err).SetType(gin.ErrorTypeAny)
}

func abortWithErrorMsg(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{
		"errMsg": msg,
	})
}

func errorsHandler(log Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			err, errType := ginErr.Err, ginErr.Type

			switch errType {
			case gin.ErrorTypeBind:
				log.Debug(err, "%s", "StatusUnprocessableEntity")
				c.Status(http.StatusUnprocessableEntity)
				return
			case gin.ErrorTypeAny:
				switch {
				case errors.Is(err, models.ErrNotValidTokens),
					errors.Is(err, models.ErrNotValidGuid):
					log.Debug(err, "%s", "StatusBadRequest")
					abortWithErrorMsg(c, http.StatusBadRequest, err.Error())
					return
				}
			}

			log.Error(err, "%s", "StatusTeapot")
			c.AbortWithStatus(http.StatusTeapot)
			return
		}
	}
}
