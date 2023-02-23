package interaction

import (
	"context"
	"douyin-demo/biz/dal/mysql"
	"douyin-demo/biz/jwt"
	"douyin-demo/biz/model"
	"strconv"
	"time"

	"douyin-demo/biz/model/douyin_core_model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// type douyin_favorite_action_request struct {
// 	token       string
// 	video_id    int64 // 视频id
// 	action_type int32 // 1-点赞，2-取消点赞
// }

type douyin_favorite_action_response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type douyin_favorite_list_response struct {
	StatusCode string                     `json:"status_code"`
	StatusMsg  string                     `json:"status_msg"`
	Video_list []*douyin_core_model.Video `json:"video_list"`
}

type douyin_comment_action_response struct {
	StatusCode int32                     `json:"status_code"`
	StatusMsg  string                    `json:"status_msg"`
	Comment    douyin_core_model.Comment `json:"comment"`
}

type douyin_comment_list_response struct {
	StatusCode   int32                        `json:"status_code"`
	StatusMsg    string                       `json:"status_msg"`
	Comment_list []*douyin_core_model.Comment `json:"comment_list"`
}

func UserFavorite(ctx context.Context, c *app.RequestContext) {
	var resp douyin_favorite_action_response
	var err error
	resp.StatusCode = 200
	resp.StatusMsg = ""

	token := c.Query("token")
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	action_type, _ := strconv.Atoi(c.Query("action_type"))

	resp.StatusCode, err = jwt.JwtAuthMiddleware(token)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	claim, err := jwt.ParseToken(token)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	user_id := claim.User_id

	fav, err := mysql.QueryUserFavorite(user_id, video_id)
	if err != nil {
		resp.StatusMsg = err.Error()
		if resp.StatusMsg == "record not found" {
			fav = 0
		} else {
			c.JSON(200, resp)
			return
		}
	}

	if action_type == 1 {
		if fav == 0 {
			resp.StatusCode = 0
			//add fav
			if err = mysql.CreateFav([]*model.Favorites{
				{
					UserId:  int64(user_id),
					VideoId: int64(video_id),
				},
			}); err != nil {
				resp.StatusMsg = err.Error()
				c.JSON(consts.StatusOK, resp)
				return
			}
			mysql.AddUserFav(user_id, video_id, 1)
		}
	} else if action_type == 2 {
		if fav != 0 {
			resp.StatusCode = 0
			//delete fav
			if err = mysql.DeleteFav(user_id, video_id); err != nil {
				resp.StatusMsg = err.Error()
				c.JSON(consts.StatusOK, resp)
				return
			}
			mysql.AddUserFav(user_id, video_id, -1)
		}
	} else {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(consts.StatusOK, resp)
}

func GetFavList(ctx context.Context, c *app.RequestContext) {
	var resp douyin_favorite_list_response
	var err error
	resp.StatusCode = "200"
	resp.StatusMsg = ""

	user_id, _ := strconv.Atoi(c.Query("user_id"))
	token := c.Query("token")
	_, err = jwt.JwtAuthMiddleware(token)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Video_list, err = mysql.QueryFavVideoList(user_id)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}
	for i := 0; i < len(resp.Video_list); i++ {
		_, err := mysql.QueryUserFavorite(user_id, int(resp.Video_list[i].Id))
		if err != nil {
			resp.StatusMsg = err.Error()
			c.JSON(200, resp)
			return
		}
		resp.Video_list[i].IsFavorite = true
	}

	resp.StatusCode = "0"
	resp.StatusMsg = "get fav list success"

	c.JSON(consts.StatusOK, resp)
}

func AddComment(ctx context.Context, c *app.RequestContext) {
	var resp douyin_comment_action_response
	resp.StatusCode = 200
	resp.StatusMsg = ""

	var err error

	token := c.Query("token")
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	action_type, _ := strconv.Atoi(c.Query("action_type"))
	comment_text := c.Query("comment_text")

	var user_id int
	if token != "" {
		claim, err := jwt.ParseToken(token)
		if err != nil {
			resp.StatusMsg = err.Error()
			c.JSON(200, resp)
			return
		}
		user_id = claim.User_id
	}

	if action_type != 1 && action_type != 2 {
		resp.StatusMsg = "error request"
		c.JSON(200, resp)
		return
	}

	if action_type == 1 {
		addtime := time.Unix(int64(time.Now().Unix()), 0).Format("01_02_15_04_05")
		if err = mysql.CreateComment([]*model.UserComments{
			{
				UserId:     int64(user_id),
				VideoId:    int64(video_id),
				Content:    comment_text,
				CreateDate: addtime,
			},
		}); err != nil {
			resp.StatusMsg = err.Error()
			c.JSON(consts.StatusOK, resp)
			return
		}
		resp.StatusMsg = "add comment success"
		resp.Comment, err = mysql.QueryCommentByUserId(user_id, user_id, addtime)
		if err != nil {
			resp.StatusMsg = err.Error()
			c.JSON(consts.StatusOK, resp)
			return
		}
		c.JSON(consts.StatusOK, resp)
	}

	comment_id, _ := strconv.Atoi(c.Query("comment_id"))
	err = mysql.DeleteComm(comment_id)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(consts.StatusOK, resp)
		return
	}
	c.JSON(consts.StatusOK, resp)

}

func GetVideoComment(ctx context.Context, c *app.RequestContext) {
	var resp douyin_comment_list_response
	resp.StatusCode = 200
	resp.StatusMsg = ""

	token := c.Query("token")
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	_, err := jwt.ParseToken(token)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Comment_list, err = mysql.QueryVideoComment(video_id)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	c.JSON(200, resp)
}
