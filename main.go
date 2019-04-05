package main

import (
	"fmt"
	"log"

	"db_forum/api"

	"db_forum/database"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {

	router := fasthttprouter.New()

	router.POST("/api/user/:nickname/create", api.CreateUser)
	router.GET("/api/user/:nickname/profile", api.GetProfile)
	router.POST("/api/user/:nickname/profile", api.UpdateProfile)

	router.GET("/api/forum/:slug/details", api.GetForum)
	router.GET("/api/forum/:slug/users", api.GetUsersByForum)
	router.GET("/api/forum/:slug/threads", api.GetThreadsByForum)
	router.POST("/api/forum/*slug", api.CreateThread)

	router.GET("/api/thread/:slug/details", api.GetThreadDetails)
	router.POST("/api/thread/:slug/details", api.UpdateThread)
	router.POST("/api/thread/:slug/create", api.CreatePost)
	router.GET("/api/thread/:slug/posts", api.GetPostsByThread)
	router.POST("/api/thread/:slug/vote", api.VoteForThread)

	router.GET("/api/post/:id/details", api.GetPost)
	router.POST("/api/post/:id/details", api.UpdatePost)

	router.GET("/api/service/status", api.Status)
	router.POST("/api/service/clear", api.Clear)

	database.Connect()
	fmt.Println("Connected")

	if err := fasthttp.ListenAndServe(":5000", router.Handler); err != nil {
		log.Fatalf("error in ListenAndServe: %s", err)
	}

}
