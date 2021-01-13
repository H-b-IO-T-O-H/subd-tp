package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/common/utils"
	"subd/application/forum"
)

type ForumHandler struct {
	//usecase forum.IUseCaseForum
	repos forum.IRepositoryForum
}

func NewForumHandler(router *fasthttprouter.Router, repos forum.IRepositoryForum) {
	f := ForumHandler{
		repos: repos,
	}
	router.POST("/api/forum/:slug", f.ForumCreate)
	router.GET("/api/forum/:slug/details", f.ForumGetDetails)
	router.GET("/api/forum/:slug/users", f.ForumGetUsers)
	router.GET("/api/forum/:slug/threads", f.ForumGetThreads)
}

func (f ForumHandler) ForumCreate(ctx *fasthttp.RequestCtx) {
	var buf models.Forum

	body := ctx.PostBody()
	if err := buf.UnmarshalJSON(body); err != nil {
		ctx.SetStatusCode(errors.BadRequestCode)
		ctx.SetBody(errors.BadRequestMsg)
		return
	}

	forumNew, err := f.repos.CreateForum(buf)
	if err != nil {
		if err.Code() == errors.ConflictCode {
			forumNew, err = f.repos.GetBySlug(buf.Slug)
			ctx.SetStatusCode(errors.ConflictCode)
		} else {
			err.SetErrToCtx(ctx)
			return
		}
	} else {
		ctx.SetStatusCode(201)
	}

	resp, errMarshal := forumNew.MarshalJSON()
	if errMarshal != nil || err != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}

func (f ForumHandler) ForumThreadCreate(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Ok")
}

func (f ForumHandler) ForumGetDetails(ctx *fasthttp.RequestCtx) {
	slug := utils.GetSlugFromCtx(ctx)
	//TODO: slug-pattern ???

	forumOld, err := f.repos.GetBySlug(slug)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := forumOld.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}

func (f ForumHandler) ForumGetThreads(ctx *fasthttp.RequestCtx) {
	query := utils.MakeQuery(ctx)
	threadsList, err := f.repos.GetThreads(query)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := threadsList.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}

func (f ForumHandler) ForumGetUsers(ctx *fasthttp.RequestCtx) {
	query := utils.MakeQuery(ctx)
	usersList, err := f.repos.GetUsers(query)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := usersList.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}
