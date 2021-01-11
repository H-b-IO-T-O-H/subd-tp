package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/common/validation"
	"subd/application/user"
)

type UserHandler struct {
	usecase user.IUseCaseUser
}

func NewUserHandler(router *fasthttprouter.Router, usecase user.IUseCaseUser) {
	u := UserHandler{
		usecase: usecase,
	}
	router.POST("/api/user/:nickname/create", u.UserCreate)
	router.GET("/api/user/:nickname/profile", u.UserGetProfile)
	router.POST("/api/user/:nickname/profile", u.UserUpdateProfile)
}

func (u UserHandler) UserCreate(ctx *fasthttp.RequestCtx) {
	var buf models.User

	nick := ctx.UserValue("nickname").(string)
	if err := validation.NicknameValid(nick); err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	if err := buf.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(errors.BadRequestCode)
		ctx.SetBody(errors.BadRequestMsg)
		return
	}
	buf.Nickname = nick
	userNew, err := u.usecase.CreateUser(buf)
	if err != nil {
		if err.Code() == errors.ConflictCode {
			existingUsers := u.usecase.GetUsersOnConflict(buf.Email, buf.Nickname)
			resp, _ := existingUsers.MarshalJSON()
			ctx.SetStatusCode(errors.ConflictCode)
			ctx.SetBody(resp)
		} else {
			err.SetErrToCtx(ctx)
		}
		return
	}
	resp, errMarshal := userNew.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetStatusCode(201)
	ctx.SetBody(resp)
}

func (u UserHandler) UserGetProfile(ctx *fasthttp.RequestCtx) {
	nick := ctx.UserValue("nickname").(string)
	if err := validation.NicknameValid(nick); err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	userOld, err := u.usecase.GetUser(nick)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := userOld.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}

func (u UserHandler) UserUpdateProfile(ctx *fasthttp.RequestCtx) {
	var buf models.User

	nick := ctx.UserValue("nickname").(string)
	if err := validation.NicknameValid(nick); err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	if err := buf.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(errors.BadRequestCode)
		ctx.SetBody(errors.BadRequestMsg)
		return
	}
	buf.Nickname = nick
	userNew, err := u.usecase.UpdateUser(buf)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := userNew.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}
