package delivery

import (
	"github.com/valyala/fasthttp"
)

type forumHandler struct {
}

func (f forumHandler) ForumCreate(ctx *fasthttp.RequestCtx) {
	//var forum models.Forum

	//body := ctx.PostBody()

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Ok")
}

func (f forumHandler) ForumCreateBranch(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Ok")
}

func (f forumHandler) ForumGetDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Ok")
}

func (f forumHandler) ForumGetBranches(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Ok")
}

func (f forumHandler) ForumGetUsers(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Ok")
}

func (f forumHandler) ForumGetBranchDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Ok")
}
