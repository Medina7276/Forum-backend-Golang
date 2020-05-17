package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"git.01.alem.school/qjawko/forum/apihandlers"
	"git.01.alem.school/qjawko/forum/dao"
	"git.01.alem.school/qjawko/forum/middleware"
	"git.01.alem.school/qjawko/forum/operations"
	"git.01.alem.school/qjawko/forum/service"
	"git.01.alem.school/qjawko/forum/singleton"
)

func main() {
	db := singleton.GetDBInstance()

	//DaoS
	userDao := dao.NewUserStore(db)
	subforumDao := dao.NewSubforumStore(db)
	postDao := dao.NewPostStore(db)
	commentDao := dao.NewCommentStore(db)
	subforumRoleDao := dao.NewSubforumRoleStore(db)
	likeDao := dao.NewLikeStore(db)

	//SERVICES
	userService := service.NewUserService(userDao)
	subforumService := service.NewSubforumService(subforumDao, nil)
	postService := service.NewPostService(postDao)
	commentService := service.NewCommentService(commentDao)
	subforumRoleService := service.NewSubforumRoleService(subforumRoleDao, subforumService)
	subforumService.SubforumRoleService = subforumRoleService
	likeService := service.NewLikeService(likeDao)
	//HANDLERS
	userHandler := apihandlers.NewUserHandler("/api/user/", userService)
	http.Handle("/api/user/", middleware.ContentTypeJson(http.HandlerFunc(userHandler.Route)))

	subforumHandler := apihandlers.NewSubforumHandler("/api/subforum/", subforumService, postService, subforumRoleService)
	http.Handle("/api/subforum/", middleware.ContentTypeJson(http.HandlerFunc(subforumHandler.Route)))

	postHandler := apihandlers.NewPostHandler("/api/post/", postService, commentService)
	http.Handle("/api/post/", middleware.ContentTypeJson(http.HandlerFunc(postHandler.Route)))

	subforumRoleHandler := apihandlers.NewSubforumRoleHandler("/api/subforumrole", subforumRoleService)
	http.Handle("/api/subforumrole/", middleware.ContentTypeJson(http.HandlerFunc(subforumRoleHandler.Route)))

	likeOps := operations.NewLikeOperations(likeService)
	userOps := operations.NewUserOperations(userService)
	http.Handle("/rate", middleware.ContentTypeJson(middleware.Authorization(http.HandlerFunc(likeOps.Rate))))
	http.Handle("/login", middleware.ContentTypeJson(http.HandlerFunc(userOps.Login)))
	http.Handle("/register", middleware.ContentTypeJson(http.HandlerFunc(userOps.Register)))
	http.Handle("/me", middleware.ContentTypeJson(middleware.Authorization(http.HandlerFunc(userOps.Me))))

	fmt.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
