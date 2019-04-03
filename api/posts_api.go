package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/iamfrommoscow/db_forum/helpers"
	"github.com/iamfrommoscow/db_forum/models"
	"github.com/valyala/fasthttp"
)

func CreatePost(ctx *fasthttp.RequestCtx) {
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
	var newPosts []*models.Post
	if err := json.Unmarshal(ctx.PostBody(), &newPosts); err != nil {
		sendInternalError(ctx, err)
	}
	if len(newPosts) == 0 {
		newPosts = make([]*models.Post, 0)
		if respBody, err := json.Marshal(newPosts); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.Write(respBody)
			ctx.SetStatusCode(fasthttp.StatusCreated)
			return
		}

	}
	created := time.Now().Format("2006-01-02T15:04:05.000Z07:00")
	for _, post := range newPosts {

		post.Thread = thread.ID
		post.Forum = thread.Forum
		if post.Created == "" {
			post.Created = created
		}
		user := helpers.FindByNickname(post.Author)
		if user == nil {
			var errorMessage = models.Error{
				Message: "Can't find user by nickname:" + post.Author,
			}
			if respBody, err := json.Marshal(errorMessage); err != nil {
				sendInternalError(ctx, err)
			} else {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				ctx.Write(respBody)
			}
			return
		}
		if post.Parent != 0 {
			parent := helpers.GetPostById(post.Parent)
			if parent == nil {

				var errorMessage = models.Error{
					Message: "Parent post was created in another thread",
				}
				if respBody, err := json.Marshal(errorMessage); err != nil {
					sendInternalError(ctx, err)
				} else {
					ctx.SetStatusCode(fasthttp.StatusConflict)
					ctx.Write(respBody)
				}
				return
			}
			if parent.Thread != post.Thread {

				var errorMessage = models.Error{
					Message: "Parent post was created in another thread",
				}
				if respBody, err := json.Marshal(errorMessage); err != nil {
					sendInternalError(ctx, err)
				} else {
					ctx.SetStatusCode(fasthttp.StatusConflict)
					ctx.Write(respBody)
				}
				return
			}
		}

	}
	if err := helpers.InsertPosts(newPosts); err == nil {
		if respBody, err := json.Marshal(newPosts); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusCreated)
			ctx.Write(respBody)
		}
	}
}

func GetPostsByThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	limit, sort, since, desc := ctx.QueryArgs().Peek("limit"), ctx.QueryArgs().Peek("sort"), ctx.QueryArgs().Peek("since"), ctx.QueryArgs().Peek("desc")
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

	posts := helpers.GetPostsByThread(thread.ID, limit, sort, since, desc)
	if len(posts) == 0 {
		posts = make([]*models.Post, 0)

	}
	if respBody, err := json.Marshal(posts); err != nil {

		sendInternalError(ctx, err)
	} else {
		ctx.Write(respBody)
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}

func GetPost(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		sendInternalError(ctx, err)
		return
	}
	post := helpers.GetPostById(id)

	if post == nil {
		var errorMessage = models.Error{
			Message: "Can't find post by id:" + string(id),
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}
	var rPost models.ReturnPost
	rPost.Pst = post
	if respBody, err := json.Marshal(rPost); err != nil {

		sendInternalError(ctx, err)
	} else {
		ctx.Write(respBody)
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}

func UpdatePost(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		sendInternalError(ctx, err)
		return
	}
	var newPost *models.Post
	if err := json.Unmarshal(ctx.PostBody(), &newPost); err != nil {
		sendInternalError(ctx, err)
	}
	post := helpers.UpdatePostById(id, newPost.Message)

	if post == nil {
		var errorMessage = models.Error{
			Message: "Can't find post by id:" + string(id),
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}
	var rPost models.ReturnPost
	rPost.Pst = post
	if respBody, err := json.Marshal(rPost); err != nil {

		sendInternalError(ctx, err)
	} else {
		ctx.Write(respBody)
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}
