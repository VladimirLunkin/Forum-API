package forum

import (
	"github.com/fasthttp/router"
)

func SetForumRouting(r *router.Router, h *Handlers) {
	r.POST("/api/forum/create", h.Create)
	r.GET("/api/forum/{slug}/details", h.GetSlug)
	r.POST("/api/forum/{slug}/create", h.CreateThread)
	r.GET("/api/forum/{slug}/threads", h.GetThreads)
}
