package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/iamfrommoscow/db_forum/helpers"
	"github.com/iamfrommoscow/db_forum/models"
	"github.com/valyala/fasthttp"
)

func CreateThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	if slug == "/create" {
		CreateForum(ctx)
		return
	}

	slug = slug[1 : len(slug)-7]

	ctx.SetContentType("application/json")
	var newThread models.Thread
	if err := json.Unmarshal(ctx.PostBody(), &newThread); err != nil {

	}
	forum := helpers.FindBySlug(slug)
	if forum == nil {
		var errorMessage = models.Error{
			Message: "Can't find thread by slug:" + slug,
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}
	user := helpers.FindByNickname(newThread.Author)
	if user == nil {
		var errorMessage = models.Error{
			Message: "Can't find user by nickname:" + newThread.Author,
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}
	if threadID, err := helpers.CreateThread(&newThread); err == nil {
		newThread.ID = threadID
		forum := helpers.FindBySlug(newThread.Forum)
		newThread.Forum = forum.Slug
		if respBody, err := json.Marshal(newThread); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusCreated)
			ctx.Write(respBody)
		}
	}

}

func GetThreadsByForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	limit, desc := ctx.QueryArgs().Peek("limit"), ctx.QueryArgs().Peek("desc")
	slug := ctx.UserValue("slug").(string)
	forum := helpers.FindBySlug(slug)
	if forum == nil {
		var errorMessage = models.Error{
			Message: "Can't find forum by slug:" + slug,
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}

	threads := helpers.GetThreadsByForum(slug, limit, desc)
	if respBody, err := json.Marshal(threads); err != nil {

		sendInternalError(ctx, err)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(respBody)
	}
}

func GetThreadDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slug := ctx.UserValue("slug").(string)
	var thread *models.Thread
	fmt.Println(slug, ":")

	if _, err := strconv.Atoi(slug); err == nil {
		thread = helpers.GetThreadByID(slug)
		fmt.Println("by id")
	} else {
		thread = helpers.GetThreadBySlug(slug)
		fmt.Println("by slug")

	}

	if thread == nil {
		var errorMessage = models.Error{
			Message: "Can't find thread by slug:" + slug,
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}

	if respBody, err := json.Marshal(thread); err != nil {

		sendInternalError(ctx, err)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(respBody)
	}
}
