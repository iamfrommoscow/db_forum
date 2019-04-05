package api

import (
	"encoding/json"

	"db_forum/helpers"

	"db_forum/models"

	"github.com/valyala/fasthttp"
)

func CreateForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	var newForum models.Forum
	if err := json.Unmarshal(ctx.PostBody(), &newForum); err != nil {

	}
	user := helpers.FindByNickname(newForum.User)
	if user == nil {
		var errorMessage = models.Error{
			Message: "Can't find user by nickname:" + newForum.User,
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}
	if err, name := helpers.CreateForum(&newForum); err == nil {
		newForum.User = name
		if respBody, err := json.Marshal(newForum); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusCreated)
			ctx.Write(respBody)
		}
	} else {
		conflictForum := helpers.FindBySlug(newForum.Slug)
		if respBody, err := json.Marshal(conflictForum); err != nil {

			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusConflict)
			ctx.Write(respBody)
		}
	}
}

func GetForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug").(string)
	forum := helpers.FindBySlug(slug)
	if forum == nil {
		var errorMessage = models.Error{
			Message: "Can't find user by nickname:" + slug,
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}
	forum.Posts = helpers.GetPostsCountByForum(forum.Slug)
	forum.Threads = helpers.GetThreadsCountByForum(forum.Slug)
	if respBody, err := json.Marshal(forum); err != nil {

		sendInternalError(ctx, err)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(respBody)
	}
}

func GetUsersByForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug").(string)
	limit, desc, since := ctx.QueryArgs().Peek("limit"), ctx.QueryArgs().Peek("desc"), ctx.QueryArgs().Peek("since")
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
	users := helpers.GetUsersBySlug(slug, limit, desc, since)
	if len(users) == 0 {
		users = make([]*models.User, 0)

	}
	if respBody, err := json.Marshal(users); err != nil {

		sendInternalError(ctx, err)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(respBody)
	}
}
