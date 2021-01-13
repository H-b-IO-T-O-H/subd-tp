package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"subd/application/common/errors"
	"subd/application/service"
)

type ServiceHandler struct {
	//usecase service.IUseCaseService
	repos service.IRepositoryService
}

func NewUserHandler(router *fasthttprouter.Router, repos service.IRepositoryService) {
	s := ServiceHandler{
		repos: repos,
	}
	router.POST("/api/service/clear", s.Clear)
	router.GET("/api/service/status", s.Status)
}

func (h ServiceHandler) Clear(ctx *fasthttp.RequestCtx) {
	err := h.repos.Clear()
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	ctx.SetBodyString("Database cleanup was successful.")
}

func (h ServiceHandler) Status(ctx *fasthttp.RequestCtx) {
	status, err := h.repos.GetStatus()
	if err != nil {
		err.SetErrToCtx(ctx)
		return
	}
	resp, errMarshal := status.MarshalJSON()
	if errMarshal != nil {
		ctx.SetStatusCode(errors.ServerErrorCode)
		ctx.SetBody(errors.ServerErrorMsg)
		return
	}
	ctx.SetBody(resp)
}
