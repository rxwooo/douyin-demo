package video_process

import (
	"context"
	"douyin-demo/biz/dal/mysql"
	"douyin-demo/biz/jwt"
	"douyin-demo/biz/model"
	"douyin-demo/biz/model/douyin_core_model"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Video_get_30(ctx context.Context, c *app.RequestContext) {
	var req douyin_core_model.DouyinFeedRequest
	var resp douyin_core_model.DouyinFeedResponse
	resp.StatusCode = 200
	resp.StatusMsg = ""

	token := c.Query("token")

	err := c.BindAndValidate(&req)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	videos, err := mysql.GetLast30Video()
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.VideoList = videos
	resp.StatusCode = 0
	resp.StatusMsg = "video feed get seccess"

	if token != "" {
		claim, err := jwt.ParseToken(token)
		if err != nil {
			resp.StatusMsg = err.Error()
			c.JSON(200, resp)
			return
		}

		user_id := claim.User_id

		for i := 0; i < len(resp.VideoList); i++ {
			_, err := mysql.QueryUserFavorite(user_id, int(resp.VideoList[i].Id))
			if err != nil {
				resp.StatusMsg = err.Error()
				c.JSON(200, resp)
				return
			}
			resp.VideoList[i].IsFavorite = true
		}
	}

	c.JSON(consts.StatusOK, resp)
}

func Publish_list(ctx context.Context, c *app.RequestContext) {
	var req douyin_core_model.DouyinPublishListRequest
	var resp douyin_core_model.DouyinPublishListResponse
	resp.StatusCode = 200
	resp.StatusMsg = ""

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
	user, err := mysql.QueryUserByID(user_id)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}
	videos, err := mysql.QueryVideoByUserId(user_id)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}

	for _, video := range videos {
		video.Author = user
	}

	resp.VideoList = videos
	resp.StatusCode = 0
	resp.StatusMsg = "video list get seccess"
	c.JSON(consts.StatusOK, resp)
}

func Video_publish(ctx context.Context, c *app.RequestContext) {
	//var req douyin_core_model.DouyinPublishActionRequest
	var resp douyin_core_model.DouyinPublishActionResponse
	resp.StatusCode = 200
	resp.StatusMsg = ""

	// err := c.BindAndValidate(&req)
	// if err != nil {
	// 	resp.StatusMsg = err.Error()
	// 	c.JSON(200, resp)
	// 	return
	// }

	token := c.PostForm("token")
	title := c.PostForm("title")
	form, err := c.MultipartForm()
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(200, resp)
		return
	}
	files := form.File["data"]

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
	for _, file := range files {
		suffix := filepath.Ext(file.Filename)
		videoName := GetVideoName(user_id)
		videoPath := "./static/" + videoName + suffix
		imagePath := "./static/" + videoName + "." + "jpg"
		//err = c.SaveUploadedFile(file, videoPath)
		savePath := filepath.Join("./static", videoName+suffix)
		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			resp.StatusMsg = err.Error()
			fmt.Println(err.Error())
			c.JSON(200, resp)
			return
		}

		err := SaveVideoName(videoPath, imagePath)
		if err != nil {
			resp.StatusMsg = err.Error()
			c.JSON(200, resp)
			return
		}

		if err = mysql.CreateVideo([]*model.Videos{
			{
				UserId:   int64(user_id),
				PlayUrl:  "http://" + "192.168.1.140" + ":" + "8888" + "/static/" + videoName + suffix,
				CoverUrl: "http://" + "192.168.1.140" + ":" + "8888" + "/static/" + videoName + "." + "jpg",
				Title:    title,
			},
		}); err != nil {
			resp.StatusMsg = err.Error()
			c.JSON(consts.StatusOK, resp)
			return
		}
		mysql.AddUserWork(user_id)
	}

	resp.StatusCode = 0
	resp.StatusMsg = "video upload seccess"
	c.JSON(consts.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "",
	})
}

func GetVideoName(user_id int) string {

	t := time.Unix(int64(time.Now().Unix()), 0)
	timeStr := t.Format("2006_01_02_15_04_05")
	videoName := strconv.Itoa(user_id) + "_" + timeStr
	return videoName
}

func SaveVideoName(src, dst string) error {
	err := ffmpeg.Input(src, ffmpeg.KwArgs{"ss": 1}).
		Output(dst, ffmpeg.KwArgs{"t": 1}).OverWriteOutput().Run()
	if err != nil && err.Error() != "exit status 0xffffffea" {
		return err
	}
	return nil
}
