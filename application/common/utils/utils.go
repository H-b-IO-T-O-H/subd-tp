package utils

import (
	"github.com/valyala/fasthttp"
	"strconv"
	"subd/application/common/models"
)

const (
	Flat       = "flat"
	ParentTree = "parent_tree"
	Tree       = "tree"
)

func GetSlugFromCtx(ctx *fasthttp.RequestCtx) string {
	return ctx.UserValue("slug").(string)
}

func GetSlugOrIdFromCtx(ctx *fasthttp.RequestCtx) string {
	return ctx.UserValue("slug_or_id").(string)
}

func GetIdFromCtx(ctx *fasthttp.RequestCtx) int64 {
	id, _ := strconv.ParseInt(ctx.UserValue("id").(string), 10, 64)
	return id
}

func GetIntFromCtxQuery(ctx *fasthttp.RequestCtx, key string) int {
	value, _ := strconv.Atoi(string(ctx.URI().QueryArgs().Peek(key)))
	return value
}

func GetBoolFromCtxQuery(ctx *fasthttp.RequestCtx, key string) bool {
	value, _ := strconv.ParseBool(string(ctx.URI().QueryArgs().Peek(key)))
	return value
}

func GetStringFromCtxQuery(ctx *fasthttp.RequestCtx, key string) string {
	return string(ctx.URI().QueryArgs().Peek(key))
}

func MakeQuery(ctx *fasthttp.RequestCtx) models.QueryParams {
	var query models.QueryParams
	query.Slug = ctx.UserValue("slug").(string)
	query.Limit, _ = strconv.Atoi(string(ctx.URI().QueryArgs().Peek("limit")))
	query.Since = string(ctx.URI().QueryArgs().Peek("since"))
	query.Desc, _ = strconv.ParseBool(string(ctx.URI().QueryArgs().Peek("desc")))
	return query
}

func MakeQueryPosts(ctx *fasthttp.RequestCtx) models.QueryPostParams {
	var query models.QueryPostParams
	query.SlugId = ctx.UserValue("slug_or_id").(string)
	query.Limit, _ = strconv.Atoi(string(ctx.URI().QueryArgs().Peek("limit")))
	query.Since, _ = strconv.ParseInt(string(ctx.URI().QueryArgs().Peek("since")), 10, 64)
	query.Desc, _ = strconv.ParseBool(string(ctx.URI().QueryArgs().Peek("desc")))
	query.Sort = string(ctx.URI().QueryArgs().Peek("sort"))
	return query
}
