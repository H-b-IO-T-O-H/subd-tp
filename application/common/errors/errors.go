package models

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"strings"
)

var InvalidNickname = []byte("invalid user nickname")
var BadRequestMsg = []byte("required fields are not filled in")
var ServerErrorMsg = []byte("something went wrong")
var NotFoundMsg = []byte("can't find this user")
var UserAlreadyExists = []byte("user already exists")
var BadRequestCode = 400
var NotFoundCode = 404
var ConflictCode = 409
var ServerErrorCode = 500

type Err interface {
	Msg() string
	Code() int
	SetErrToCtx(ctx *fasthttp.RequestCtx)
}

type RespErr struct {
	Message    []byte
	StatusCode int
}

func (r RespErr) Msg() string {
	return string(r.Message)
}

func (r RespErr) Code() int {
	return r.StatusCode
}

func (r RespErr) SetErrToCtx(ctx *fasthttp.RequestCtx) {
	msg := []byte(fmt.Sprintf("{\"message\": \"%s\"}", r.Message))
	ctx.SetStatusCode(r.StatusCode)
	ctx.SetBody(msg)
}

func UserNotFound(errMsg string) bool {
	return strings.Contains(errMsg, "user")
}

func UserExists(errMsg string) bool {
	return strings.Contains(errMsg, "duplicate")
}

func EmptyResult(errMsg string) bool {
	return strings.Contains(errMsg, "no rows")
}
