package forum

import (
	"github.com/fasthttp/router"
)

func SetForumRouting(r *router.Router, h *Handlers) {
	r.POST("/api/forum/create", h.Create)
}
