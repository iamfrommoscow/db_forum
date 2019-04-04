package api

import (
	"encoding/json"

	"github.com/iamfrommoscow/db_forum/helpers"
	"github.com/valyala/fasthttp"
)

func Status(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	service := helpers.GetCount()
	if respBody, err := json.Marshal(service); err != nil {
		sendInternalError(ctx, err)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(respBody)
	}
}

func Clear(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	err := helpers.DropDatabase()
	if err != nil {
		sendInternalError(ctx, err)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}
