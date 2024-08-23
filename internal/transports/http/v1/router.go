package httpv1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/v1adhope/auth-service/internal/services"
)

type Router struct {
	as  *services.Services
	log Logger
}

func New(
	as *services.Services,
	log Logger,
) *Router {
	return &Router{
		as:  as,
		log: log,
	}
}

func (r *Router) Handler(opts ...Option) *gin.Engine {
	cfg := config(opts...)

	gin.SetMode(cfg.Mode)

	e := gin.New()

	e.Use(
		gin.Logger(),
		gin.Recovery(),
		cors.New(cfg.Cors),
		errorsHandler(r.log),
	)

	apiG := e.Group("/v1")
	{
		initAuthRouter(&authRouter{apiG, r.as})
	}

	return e
}
