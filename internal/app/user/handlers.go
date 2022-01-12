package user

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
	UserRepo models.UserRep
	Logger   *zap.SugaredLogger
}

func (h *Handlers) CreateUser(ctx *fasthttp.RequestCtx) {
	var user models.User
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		delivery.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user.Nickname = fmt.Sprintf("%s", ctx.UserValue("nickname"))

	users, err := h.UserRepo.CreateUser(user)
	if err != nil {
		if err == models.UserExistsError {
			delivery.Send(ctx, http.StatusConflict, users)
			return
		}
		delivery.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	delivery.Send(ctx, http.StatusCreated, user)
}

func (h *Handlers) GetUser(ctx *fasthttp.RequestCtx) {
	nickname := fmt.Sprintf("%s", ctx.UserValue("nickname"))

	user, err := h.UserRepo.GetUserByNickname(nickname)
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	delivery.Send(ctx, http.StatusOK, user)
}

func (h *Handlers) UpdateUser(ctx *fasthttp.RequestCtx) {
	newUserData, err := h.UserRepo.GetUserByNickname(fmt.Sprintf("%s", ctx.UserValue("nickname")))
	if err != nil {
		delivery.SendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	err = json.Unmarshal(ctx.PostBody(), &newUserData)
	if err != nil {
		delivery.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.UserRepo.UpdateUser(newUserData)
	if err != nil {
		if err == models.NewUserDataError {
			delivery.SendError(ctx, http.StatusConflict, err.Error())
			return
		}
		delivery.SendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	delivery.Send(ctx, http.StatusOK, user)
}
