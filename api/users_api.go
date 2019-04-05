package api

import (
	"encoding/json"

	"db_forum/helpers"

	"db_forum/models"

	"github.com/valyala/fasthttp"
)

func sendInternalError(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	// log.Fatal(err)
}

func CreateUser(ctx *fasthttp.RequestCtx) {
	var newUser models.User
	if err := json.Unmarshal(ctx.PostBody(), &newUser); err != nil {

	}

	newUser.Nickname = ctx.UserValue("nickname").(string)

	ctx.SetContentType("application/json")

	if err := helpers.CreateUser(&newUser); err == nil {
		if respBody, err := json.Marshal(newUser); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusCreated)
			ctx.Write(respBody)
		}
	} else {

		users := helpers.FindByNicknameOrEmail(newUser.Nickname, newUser.Email)
		if respBody, err := json.Marshal(users); err != nil {

			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusConflict)
			ctx.Write(respBody)
		}
	}
}

func GetProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	nickname := ctx.UserValue("nickname").(string)
	user := helpers.FindByNickname(nickname)
	if user == nil {
		var errorMessage = models.Error{
			Message: "Can't find user by nickname:" + nickname,
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}
	if respBody, err := json.Marshal(user); err != nil {
		sendInternalError(ctx, err)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(respBody)
	}
}

func UpdateProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	var user models.User
	if err := json.Unmarshal(ctx.PostBody(), &user); err != nil {

	}
	user.Nickname = ctx.UserValue("nickname").(string)
	findedUser := helpers.FindByNickname(user.Nickname)

	if findedUser == nil {
		var errorMessage = models.Error{
			Message: "Can't find user by nickname:" + user.Nickname,
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(respBody)
		}
		return
	}
	if err := helpers.UpdateProfile(&user); err != nil {
		sameEmailUser := helpers.FindByEmail(user.Email)
		var errorMessage = models.Error{
			Message: "This email is already registered by user:" + sameEmailUser.Nickname,
		}
		if respBody, err := json.Marshal(errorMessage); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusConflict)
			ctx.Write(respBody)
		}
		return
	} else {

		if respBody, err := json.Marshal(user); err != nil {
			sendInternalError(ctx, err)
		} else {
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.Write(respBody)
		}
	}

}
