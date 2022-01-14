package forum

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
	"tech-db-forum/internal/app/models"
	"tech-db-forum/internal/pkg/delivery"
)

type Handlers struct {
	ForumRepo models.ForumRep
	UserRepo  models.UserRep
	Logger    *zap.SugaredLogger
}

func (h *Handlers) Create(ctx *fasthttp.RequestCtx) {
	var newForum models.Forum
	err := json.Unmarshal(ctx.PostBody(), &newForum)
	if err != nil {
		delivery.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.UserRepo.GetUserByNickname(newForum.User)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}
	newForum.User = user.Nickname

	forum, err := h.ForumRepo.CreateForum(newForum)
	if err == models.SlugAlreadyExistsError {
		forum, _ = h.ForumRepo.GetForumBySlug(newForum.Slug)
		delivery.Send(ctx, http.StatusConflict, forum)
		return
	}
	if err != nil {
		delivery.SendError(ctx, http.StatusConflict, err.Error())
		return
	}

	delivery.Send(ctx, http.StatusCreated, forum)
}

func (h *Handlers) GetSlug(ctx *fasthttp.RequestCtx) {
	slug := fmt.Sprintf("%s", ctx.UserValue("slug"))

	forum, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}

	delivery.Send(ctx, http.StatusOK, forum)
	return
}

func (h *Handlers) CreateThread(ctx *fasthttp.RequestCtx) {
	var newThread models.Thread
	err := json.Unmarshal(ctx.PostBody(), &newThread)
	if err != nil {
		delivery.SendError(ctx, http.StatusBadRequest, "")
		return
	}

	user, err := h.UserRepo.GetUserByNickname(newThread.Author)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}
	newThread.Author = user.Nickname

	_, err = h.ForumRepo.GetForumBySlug(fmt.Sprintf("%s", ctx.UserValue("slug")))
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}
	//newThread.Slug = forum.Slug

	thread, err := h.ForumRepo.CreateThread(newThread)
	if err != nil {
		delivery.SendError(ctx, http.StatusConflict, "")
		return
	}

	delivery.Send(ctx, http.StatusCreated, thread)
}

func (h *Handlers) GetThreads(ctx *fasthttp.RequestCtx) {
	slug := fmt.Sprintf("%s", ctx.UserValue("slug"))

	limit := fmt.Sprintf("%s", ctx.FormValue("limit"))
	if limit == "" {
		limit = "100"
	}

	since := fmt.Sprintf("%s", ctx.FormValue("since"))

	desc := fmt.Sprintf("%s", ctx.FormValue("desc"))

	fmt.Println(slug, limit, since, desc)
	fmt.Println()

	threads, err := h.ForumRepo.GetThreads(slug, limit, since, desc)
	if err != nil || len(threads) == 0 {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}

	delivery.Send(ctx, http.StatusCreated, threads)
}
