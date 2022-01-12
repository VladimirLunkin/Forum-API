package user

import (
	"github.com/fasthttp/router"
)

func SetUserRouting(r *router.Router, h *Handlers) {
	r.POST("/api/user/{nickname}/create", h.CreateUser)
	r.GET("/api/user/{nickname}/profile", h.GetUser)
	//r.POST("/api/user/{nickname}/profile", h.Create)
}
