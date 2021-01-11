package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"strconv"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/common/utils"
	"subd/application/thread"
)

const BySlug = -1

type ThreadHandler struct {
	usecase thread.IUseCaseThread
}

func NewThreadHandler(router *fasthttprouter.Router, usecase thread.IUseCaseThread) {
	f := ThreadHandler{
		usecase: usecase,
	}
	router.POST("/forum/:slug/create", f.ThreadCreate)
	router.GET("/thread/:slug_or_id/details", f.ThreadDetails)
	router.POST("/thread/:slug_or_id/details", f.ThreadUpdate)
	router.POST("/thread/:slug_or_id/vote", f.ThreadVote)
}

func (t ThreadHandler) ThreadCreate(ctx *fasthttp.RequestCtx) {
	var buf models.Thread

	//fmt.Print(time.Now().Format(time.RFC3339))
	if err := buf.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(errors.BadRequestCode)
		ctx.SetBody(errors.BadRequestMsg)
		return
	}
	buf.Forum = utils.GetSlugFromCtx(ctx)
	if buf.Slug == "" {
		buf.Slug = buf.Forum + uuid.New().String()
	}
	threadNew, err := t.usecase.CreateThread(buf)
	if err != nil {
		if err.Code() == errors.ConflictCode {
			threadNew, err = t.usecase.GetBySlugOrId(buf.Slug, BySlug)
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
		threadOld, err = t.usecase.GetBySlugOrId("", id)
	} else {
		threadOld, err = t.usecase.GetBySlugOrId(slugId, BySlug)
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
	threadNew, err := t.usecase.UpdateBySlugOrId(buf, slugId)
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
	threadVoted, err := t.usecase.CreateVote(buf)
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
