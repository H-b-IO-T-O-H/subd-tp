package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strconv"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/common/utils"
	"subd/application/thread"
)

const BySlug = -1

type ThreadHandler struct {
	//usecase thread.IUseCaseThread
	repos thread.IRepositoryThread
}

func NewThreadHandler(router *fasthttprouter.Router, repos thread.IRepositoryThread) {
	f := ThreadHandler{
		repos: repos,
	}
	router.POST("/api/forum/:slug/create", f.ThreadCreate)
	router.GET("/api/thread/:slug_or_id/details", f.ThreadDetails)
	router.POST("/api/thread/:slug_or_id/details", f.ThreadUpdate)
	router.POST("/api/thread/:slug_or_id/vote", f.ThreadVote)
	router.GET("/api/thread/:slug_or_id/posts", f.ThreadPosts)
}

func (t ThreadHandler) ThreadCreate(ctx *fasthttp.RequestCtx) {
	var buf models.Thread

	if err := buf.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(errors.BadRequestCode)
		ctx.SetBody(errors.BadRequestMsg)
		return
	}
	buf.Forum = utils.GetSlugFromCtx(ctx)
	threadNew, err := t.repos.CreateThread(buf)
	if err != nil {
		if err.Code() == errors.ConflictCode {
			threadNew, err = t.repos.GetBySlug(buf.Slug)
			ctx.SetStatusCode(errors.ConflictCode)
		} else {
			err.SetErrToCtx(ctx)
			return
		}
	} else {
		ctx.SetStatusCode(201)
	}

	resp, errMarshal := threadNew.MarshalJSON()
	if errMarshal != nil || err != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}

func (t ThreadHandler) ThreadDetails(ctx *fasthttp.RequestCtx) {
	var (
		threadOld models.Thread
		err       errors.Err
	)
	slugId := utils.GetSlugOrIdFromCtx(ctx)
	if id, _ := strconv.Atoi(slugId); id != 0 {
		threadOld, err = t.repos.GetById(id)
	} else {
		threadOld, err = t.repos.GetBySlug(slugId)
	}
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := threadOld.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}

func (t ThreadHandler) ThreadUpdate(ctx *fasthttp.RequestCtx) {
	var buf models.Thread

	if err := buf.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(errors.BadRequestCode)
		ctx.SetBody(errors.BadRequestMsg)
		return
	}
	slugId := utils.GetSlugOrIdFromCtx(ctx)
	threadNew, err := t.repos.UpdateBySlugOrId(buf, slugId)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := threadNew.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}

func (t ThreadHandler) ThreadVote(ctx *fasthttp.RequestCtx) {
	var buf models.Vote
	if err := buf.UnmarshalJSON(ctx.PostBody()); err != nil || (buf.Voice != -1 && buf.Voice != 1) {
		ctx.SetStatusCode(errors.BadRequestCode)
		ctx.SetBody(errors.BadRequestMsg)
		return
	}
	buf.SlugOrId = utils.GetSlugOrIdFromCtx(ctx)
	threadVoted, err := t.repos.UpsertVote(buf)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := threadVoted.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}

func (t ThreadHandler) ThreadPosts(ctx *fasthttp.RequestCtx) {
	query := utils.MakeQueryPosts(ctx)

	if query.Sort == "" {
		query.Sort = utils.Flat
	}
	posts, err := t.repos.GetThreadPosts(query)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := posts.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}
