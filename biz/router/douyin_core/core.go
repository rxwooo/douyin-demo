// Code generated by hertz generator. DO NOT EDIT.

package douyin_core

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	user_service "douyin-demo/biz/handler/user_service"
	"douyin-demo/biz/handler/video_process"
	"douyin-demo/biz/handler/interaction"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {
	r.Static("/static", "./")
	douyin := r.Group("/douyin", rootMw()...)
	{
		douyin.GET("/feed/", video_process.Video_get_30)
		user := douyin.Group("/user")
		{
			user.GET("/", user_service.User_info)
			user.POST("/register/", user_service.User_register)
			user.POST("/login/", user_service.User_login)
		}
		publish := douyin.Group("/publish")
		{
			publish.POST("/action/", video_process.Video_publish)
			publish.GET("/list/", video_process.Publish_list)
		}
		fav := douyin.Group("/favorite")
		{
			fav.POST("/action/", interaction.UserFavorite)
			fav.GET("/list/", interaction.GetFavList)
		}
		douyin.GET("/comment/list/", interaction.GetVideoComment)
		douyin.POST("/comment/action/", interaction.AddComment)
	}
}
