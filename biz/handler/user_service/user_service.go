package user_service

import (
	"context"
	"douyin-demo/biz/dal/mysql"
	"douyin-demo/biz/jwt"
	"douyin-demo/biz/model"
	"douyin-demo/biz/model/douyin_core_model"
	"fmt"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func User_register(ctx context.Context, c *app.RequestContext) {
	//fmt.Println(c.Request.Body())
	var req douyin_core_model.DouyinUserLoginRequest
	var resp douyin_core_model.DouyinUserRegisterResponse
	resp.StatusCode = 400
	resp.UserId = -1
	resp.Token = ""
	err := c.BindAndValidate(&req)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	username := c.Query("username")
	password := c.Query("password")

	if err = mysql.CreateUser([]*model.Users{
		{
			Name:       username,
			PasswdHash: jwt.MD5(password),
		},
	}); err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "register success"
	resp.UserId, _ = mysql.QueryUserId(username)
	resp.Token, _ = jwt.GetToken(int(resp.UserId))

	c.JSON(consts.StatusOK, resp)
}

func User_login(ctx context.Context, c *app.RequestContext) {
	var req douyin_core_model.DouyinUserLoginRequest
	var resp douyin_core_model.DouyinUserLoginResponse
	resp.StatusCode = 400
	resp.UserId = -1
	resp.Token = ""

	err := c.BindAndValidate(&req)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}
	username := c.Query("username")
	password := c.Query("password")

	u, err := mysql.CheckUser(username, jwt.MD5(password))
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "login success"
	resp.UserId = int64(u.Id)
	resp.Token, _ = jwt.GetToken(int(resp.UserId))

	c.JSON(consts.StatusOK, resp)
}

func User_info(ctx context.Context, c *app.RequestContext) {
	var req douyin_core_model.DouyinUserRequest
	var resp douyin_core_model.DouyinUserResponse
	resp.StatusCode = 400

	err := c.BindAndValidate(&req)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	user_id, _ := strconv.Atoi(c.Query("user_id"))
	token := c.Query("token")
	fmt.Println(token)
	resp.StatusCode, err = jwt.JwtAuthMiddleware(token)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.User, err = mysql.QueryUserByID(user_id)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "user info get success"

	c.JSON(consts.StatusOK, resp)
}
