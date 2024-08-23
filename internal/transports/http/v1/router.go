package httpv1

import "github.com/gin-gonic/gin"

type Router struct {
	as  AuthService
	log Logger
}

func New(
	as AuthService,
	log Logger,
) *Router {
	return &Router{
		as:  as,
		log: log,
	}
}

func (r *Router) New() *gin.Engine {
	e := gin.New()

	e.Use(gin.Logger(), gin.Recovery(), corsHandler(), errorsHandler(r.log))

	apiG := e.Group("/v1")
	{
		initAuthRouter(&authRouter{apiG, r.as})
	}

	return e
}
