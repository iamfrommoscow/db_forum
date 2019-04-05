package api

import (
	"encoding/json"
	"strconv"

	"db_forum/helpers"

	"db_forum/models"

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

	newThread.Forum = slug
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

	// fmt.Println("")
	// fmt.Println(newThread.Author)
	// fmt.Println(newThread.Forum)
	// fmt.Println(newThread.Message)
	// fmt.Println(newThread.Title)
	// fmt.Println(newThread.Slug)
	// fmt.Println("URI", string(ctx.RequestURI()))
	// fmt.Println("")
	var foundThread *models.Thread
	if newThread.Slug != "" {
		foundThread = helpers.GetThreadBySlug(newThread.Slug)
	}
	if foundThread == nil {

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
	} else {

		if respBody, err := json.Marshal(foundThread); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusConflict)
			ctx.Write(respBody)
		}
		return
	}

}

func GetThreadsByForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	limit, desc, since := ctx.QueryArgs().Peek("limit"), ctx.QueryArgs().Peek("desc"), ctx.QueryArgs().Peek("since")
	slug := ctx.UserValue("slug").(string)
	forum := helpers.FindBySlug(slug)
	// fmt.Println(since)
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
	threads := helpers.GetThreadsByForum(slug, limit, desc, since)
	if len(threads) == 0 {
		threads = make([]*models.Thread, 0)

	}
	if respBody, err := json.Marshal(threads); err != nil {

		sendInternalError(ctx, err)
	} else {
		ctx.Write(respBody)
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}

func GetThreadDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slug := ctx.UserValue("slug").(string)
	var thread *models.Thread

	if _, err := strconv.Atoi(slug); err == nil {
		thread = helpers.GetThreadByID(slug)
	} else {
		thread = helpers.GetThreadBySlug(slug)

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

func VoteForThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slug := ctx.UserValue("slug").(string)
	var thread *models.Thread

	if _, err := strconv.Atoi(slug); err == nil {
		thread = helpers.GetThreadByID(slug)
	} else {
		thread = helpers.GetThreadBySlug(slug)

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
	var newVote models.Vote
	if err := json.Unmarshal(ctx.PostBody(), &newVote); err != nil {

	}
	thread = helpers.VoteForThread(thread.ID, &newVote)
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
	return
}

func UpdateThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slug := ctx.UserValue("slug").(string)
	var thread *models.Thread

	if _, err := strconv.Atoi(slug); err == nil {
		thread = helpers.GetThreadByID(slug)
	} else {
		thread = helpers.GetThreadBySlug(slug)

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
	var newThread models.Thread
	if err := json.Unmarshal(ctx.PostBody(), &newThread); err != nil {

	}
	if newThread.Message == "" {
		newThread.Message = thread.Message
	}
	if newThread.Title == "" {
		newThread.Title = thread.Title
	}

	if _, err := strconv.Atoi(slug); err == nil {
		thread = helpers.UpdateThreadByID(slug, newThread.Message, newThread.Title)
	} else {
		thread = helpers.UpdateThreadBySlug(slug, newThread.Message, newThread.Title)

	}

	if respBody, err := json.Marshal(thread); err != nil {

		sendInternalError(ctx, err)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(respBody)
	}
}
