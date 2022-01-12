package forum

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
	"tech-db-forum/internal/app/models"
	"tech-db-forum/internal/pkg/delivery"
)

type Handlers struct {
	ForumRepo models.ForumRep
	Logger    *zap.SugaredLogger
}

func (h *Handlers) Create(ctx *fasthttp.RequestCtx) {
	var forum models.Forum
	err := json.Unmarshal(ctx.PostBody(), &forum)
	if err != nil {
		delivery.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	forum, err = h.ForumRepo.Create(forum)
	if err != nil {
		delivery.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	delivery.Send(ctx, http.StatusCreated, forum)
}
