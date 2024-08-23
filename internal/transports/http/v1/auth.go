package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/auth-service/internal/models"
)

type authRouter struct {
	apiG *gin.RouterGroup
	as   AuthService
}

func initAuthRouter(r *authRouter) {
	tokensG := r.apiG.Group("/tokens")
	{
		tokensG.POST("/:userId", r.tokenPair)
		tokensG.POST("/refresh", r.refreshTokenPair)
	}
}

type tokenPairPathParam struct {
	UserId string `uri:"userId"`
}

func (r *authRouter) tokenPair(c *gin.Context) {
	pathParams := tokenPairPathParam{}

	if err := c.ShouldBindUri(&pathParams); err != nil {
		setBindError(c, err)
		return
	}

	tp, err := r.as.GenerateTokenPair(c.Request.Context(), pathParams.UserId, c.ClientIP())
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusCreated, tp)
}

type refreshTokenPairReq struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

func (r *authRouter) refreshTokenPair(c *gin.Context) {
	req := refreshTokenPairReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	tp := models.TokenPair{
		Access:  req.Access,
		Refresh: req.Refresh,
	}

	newTp, err := r.as.RefreshTokenPair(c.Request.Context(), tp, c.ClientIP())
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusCreated, newTp)
}
