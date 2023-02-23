package mysql

import (
	"douyin-demo/biz/model"
	"douyin-demo/biz/model/douyin_core_model"
)

func CreateVideo(videos []*model.Videos) error {
	return DB.Create(videos).Error
}

func modelVideos2CoreVideos(videos []*model.Videos) ([]*douyin_core_model.Video, error) {
	res := make([]*douyin_core_model.Video, 0)
	for _, video := range videos {
		var newVid douyin_core_model.Video
		var err error
		newVid.Author, err = QueryUserByID(int(video.UserId))
		if err != nil {
			return res, err
		}
		newVid.Id = int64(video.Id)
		newVid.CommentCount = video.CommentCount
		newVid.CoverUrl = video.CoverUrl
		newVid.FavoriteCount = video.FavoriteCount
		newVid.IsFavorite = false
		newVid.PlayUrl = video.PlayUrl
		newVid.Title = video.Title

		res = append(res, &newVid)
	}
	return res, nil
}

func QueryVideoByUserId(user_id int) ([]*douyin_core_model.Video, error) {
	res := make([]*model.Videos, 0)
	if err := DB.Where("user_id = ?", user_id).Find(&res).Error; err != nil {
		return nil, err
	}
	return modelVideos2CoreVideos(res)
}

func GetLast30Video() ([]*douyin_core_model.Video, error) {
	res := make([]*model.Videos, 0)
	if err := DB.Limit(30).Find(&res).Error; err != nil {
		return nil, err
	}
	return modelVideos2CoreVideos(res)
}

func QueryUserIdByVideoId(video_id int) (int, error) {
	var res model.Videos
	if video_id > 0 {
		if err := DB.Where("id = ?", video_id).First(&res).Error; err != nil {
			return -1, err
		}
	}
	return int(res.UserId), nil
}
