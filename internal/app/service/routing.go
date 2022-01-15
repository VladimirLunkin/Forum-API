package service

import "github.com/fasthttp/router"

func SetServiceRouting(r *router.Router, h *Handlers) {
	r.GET("/api/service/status", h.ServiceStatus)
	r.POST("/api/service/clear", h.ServiceClear)
}
