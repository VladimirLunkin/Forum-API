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

func (h *Handlers) CreateForum(ctx *fasthttp.RequestCtx) {
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

	forum, err := h.ForumRepo.GetForumBySlug(fmt.Sprintf("%s", ctx.UserValue("slug")))
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}
	newThread.Forum = forum.Slug

	thread, err := h.ForumRepo.CreateThread(newThread)
	if err != nil {
		if err == models.ThreadAlreadyExistsError {
			delivery.Send(ctx, http.StatusConflict, thread)
			return
		}
		delivery.SendError(ctx, http.StatusConflict, "")
		return
	}

	delivery.Send(ctx, http.StatusCreated, thread)
}

func (h *Handlers) GetThreads(ctx *fasthttp.RequestCtx) {
	slug := fmt.Sprintf("%s", ctx.UserValue("slug"))
	_, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}

	limit := fmt.Sprintf("%s", ctx.FormValue("limit"))
	if limit == "" {
		limit = "100"
	}

	since := fmt.Sprintf("%s", ctx.FormValue("since"))

	desc := ""
	if fmt.Sprintf("%s", ctx.FormValue("desc")) == "true" {
		desc = "desc"
	}

	threads, err := h.ForumRepo.GetThreads(slug, limit, since, desc)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}

	delivery.Send(ctx, http.StatusOK, threads)
}

func (h *Handlers) CreatePosts(ctx *fasthttp.RequestCtx) {
	slugOrId := fmt.Sprintf("%s", ctx.UserValue("slug_or_id"))
	thread, err := h.ForumRepo.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		fmt.Println("assssssss ", thread, err)
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}

	var newPosts []models.Post
	err = json.Unmarshal(ctx.PostBody(), &newPosts)
	if err != nil {
		delivery.SendError(ctx, http.StatusBadRequest, "")
		return
	}

	posts, err := h.ForumRepo.CreatePosts(thread, newPosts)
	if err != nil {
		delivery.SendError(ctx, http.StatusConflict, err.Error())
		return
	}

	delivery.Send(ctx, http.StatusCreated, posts)
}

func (h *Handlers) Vote(ctx *fasthttp.RequestCtx) {
	slugOrId := fmt.Sprintf("%s", ctx.UserValue("slug_or_id"))
	thread, err := h.ForumRepo.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}

	var vote models.Vote
	err = json.Unmarshal(ctx.PostBody(), &vote)
	if err != nil {
		delivery.SendError(ctx, http.StatusBadRequest, "")
		return
	}

	thread, err = h.ForumRepo.Vote(thread, vote)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}

	delivery.Send(ctx, http.StatusOK, thread)
}

func (h *Handlers) ThreadDetails(ctx *fasthttp.RequestCtx) {
	slugOrId := fmt.Sprintf("%s", ctx.UserValue("slug_or_id"))
	thread, err := h.ForumRepo.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}

	delivery.Send(ctx, http.StatusOK, thread)
}

func (h *Handlers) GetPosts(ctx *fasthttp.RequestCtx) {
	slugOrId := fmt.Sprintf("%s", ctx.UserValue("slug_or_id"))
	thread, err := h.ForumRepo.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, "")
		return
	}

	limit := fmt.Sprintf("%s", ctx.FormValue("limit"))
	if limit == "" {
		limit = "100"
	}

	since := fmt.Sprintf("%s", ctx.FormValue("since"))

	sort := fmt.Sprintf("%s", ctx.FormValue("sort"))
	if sort == "" {
		sort = "flat"
	}

	desc := ""
	if fmt.Sprintf("%s", ctx.FormValue("desc")) == "true" {
		desc = "desc"
	}

	posts, err := h.ForumRepo.GetPosts(thread, limit, since, sort, desc)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	delivery.Send(ctx, http.StatusOK, posts)
}
