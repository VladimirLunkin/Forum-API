package service

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
	"tech-db-forum/internal/app/models"
	"tech-db-forum/internal/pkg/delivery"
)

type Handlers struct {
	ServiceRepo models.ServiceRep
	Logger      *zap.SugaredLogger
}

func (h *Handlers) ServiceStatus(ctx *fasthttp.RequestCtx) {
	status, err := h.ServiceRepo.GetStatus()
	if err != nil {
		delivery.Send(ctx, http.StatusOK, nil)
		return
	}

	delivery.Send(ctx, http.StatusOK, status)
}

func (h *Handlers) ServiceClear(ctx *fasthttp.RequestCtx) {
	if err := h.ServiceRepo.Clear(); err != nil {
		delivery.SendError(ctx, http.StatusInternalServerError, "")
	}

	delivery.SendOK(ctx)
}
