package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strings"
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/common/utils"
	"subd/application/post"
)

type PostHandler struct {
	//usecase post.IUseCasePost
	repos post.IRepositoryPost
}

func NewPostHandler(router *fasthttprouter.Router, repos post.IRepositoryPost) {
	f := PostHandler{
		repos: repos,
	}
	router.POST("/api/thread/:slug_or_id/create", f.PostCreate)
	router.GET("/api/post/:id/details", f.PostDetails)
	router.POST("/api/post/:id/details", f.PostUpdate)
}

func (p PostHandler) PostCreate(ctx *fasthttp.RequestCtx) {
	buf := models.PostsList{}
	if err := buf.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(errors.BadRequestCode)
		ctx.SetBody(errors.BadRequestMsg)
		return
	}
	slugId := utils.GetSlugOrIdFromCtx(ctx)
	posts, err := p.repos.CreatePost(buf, slugId)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	var resp[] byte
	var errM error
	if posts == nil || len(*posts) == 0 {
		resp, _ = models.PostsList{}.MarshalJSON()
	} else {
		resp, errM = posts.MarshalJSON()
		if errM != nil {
			ctx.SetStatusCode(errors.ServerErrorCode)
			ctx.SetBody(errors.ServerErrorMsg)
			return
		}
	}
	ctx.SetStatusCode(201)
	ctx.SetBody(resp)
}

func (p PostHandler) PostDetails(ctx *fasthttp.RequestCtx) {
	getParams := models.PostGetParams{PostId: utils.GetIdFromCtx(ctx)}
	params := strings.Split(utils.GetStringFromCtxQuery(ctx, "related"), ",")
	for _, p := range params {
		if p == "user" {
			getParams.HaveUser = true
		} else if p == "forum" {
			getParams.HaveForum = true
		} else if p == "thread" {
			getParams.HaveThread = true
		}
	}
	postOld, err := p.repos.GetFull(getParams)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := postOld.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}

func (p PostHandler) PostUpdate(ctx *fasthttp.RequestCtx) {
	var buf models.PostUpdate

	buf.ID = utils.GetIdFromCtx(ctx)
	if err := buf.UnmarshalJSON(ctx.PostBody()); err != nil || buf.ID == 0 {
		ctx.SetStatusCode(errors.BadRequestCode)
		ctx.SetBody(errors.BadRequestMsg)
		return
	}
	postUpdate, err := p.repos.UpdatePost(buf)
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := postUpdate.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}
