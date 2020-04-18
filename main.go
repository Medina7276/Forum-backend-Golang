package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"git.01.alem.school/qjawko/forum/apihandlers"
	"git.01.alem.school/qjawko/forum/middleware"
	"git.01.alem.school/qjawko/forum/operations"
	"git.01.alem.school/qjawko/forum/repo"
	"git.01.alem.school/qjawko/forum/service"
	"git.01.alem.school/qjawko/forum/singleton"
)

func main() {
	db := singleton.GetDBInstance()

	//REPOS
	userRepo := repo.NewUserStore(db)
	subforumRepo := repo.NewSubforumStore(db)
	postRepo := repo.NewPostStore(db)
	//SERVICES
	userService := service.NewUserService(userRepo)
	subforumService := service.NewSubforumService(subforumRepo)
	postService := service.NewPostService(postRepo)
	//HANDLERS
	userHandler := apihandlers.NewUserHandler("/api/user/", userService)
	http.Handle("/api/user/", middleware.ContentTypeJson(http.HandlerFunc(userHandler.Route)))

	subforumHandler := apihandlers.NewSubforumHandler("/api/subforum/", subforumService)
	http.Handle("/api/subforum/", middleware.ContentTypeJson(http.HandlerFunc(subforumHandler.Route)))

	postHandler := apihandlers.NewPostHandler("/api/post/", postService)
	http.Handle("/api/post/", middleware.ContentTypeJson(http.HandlerFunc(postHandler.Route)))

	userOps := operations.NewUserOperations(userService)

	http.Handle("/login", middleware.ContentTypeJson(middleware.Unauthorized(http.HandlerFunc(userOps.Login))))
	http.Handle("/register", middleware.ContentTypeJson(middleware.Unauthorized(http.HandlerFunc(userOps.Register))))
	http.Handle("/me", middleware.ContentTypeJson(middleware.Authorization(http.HandlerFunc(userOps.Me))))

	fmt.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
