package mysql

import (
	"douyin-demo/biz/model"
	"douyin-demo/biz/model/douyin_core_model"
)

func CreateFav(fav []*model.Favorites) error {
	return DB.Create(fav).Error
}

func DeleteFav(user_id, video_id int) error {
	return DB.Where("user_id = ? AND video_id = ?", user_id, video_id).Delete(&model.Favorites{}).Error
}

func DeleteComm(comm_id int) error {
	return DB.Where("id = ?", comm_id).Delete(&model.Favorites{}).Error
}

func QueryUserFavorite(user_id, video_id int) (int, error) {
	var fav model.Favorites
	if err := DB.First(&fav, "user_id = ? AND video_id = ?", user_id, video_id).Error; err != nil {
		return -1, err
	}
	return int(fav.Id), nil
}

func AddUserFav(user_id, video_id, num int) error {
	var user model.Users

	if err := DB.Where("id = ?", user_id).First(&user).Error; err != nil {
		return err
	}
	user.FavoriteCount += int64(num)
	DB.Save(&user)

	var video model.Videos
	if err := DB.Where("id = ?", video_id).First(&video).Error; err != nil {
		return err
	}
	video.FavoriteCount += int64(num)
	user_id = int(video.UserId)
	DB.Save(&video)

	if err := DB.Where("id = ?", user_id).First(&user).Error; err != nil {
		return err
	}
	user.TotalFavorited += int64(num)
	DB.Save(&user)

	return nil
}

func QueryFavVideoList(user_id int) ([]*douyin_core_model.Video, error) {
	tpRes := make([]*model.Videos, 0)
	res := make([]*douyin_core_model.Video, 0)
	favs := make([]*model.Favorites, 0)
	if err := DB.Where("user_id = ?", user_id).Find(&favs).Error; err != nil {
		return res, err
	}
	for _, fav := range favs {
		var re model.Videos
		if err := DB.Where("id = ?", fav.VideoId).First(&re).Error; err != nil {
			return res, err
		}
		tpRes = append(tpRes, &re)
	}
	return modelVideos2CoreVideos(tpRes)
}

func CreateComment(com []*model.UserComments) error {
	return DB.Create(com).Error
}

func commentTrans(comm model.UserComments) (douyin_core_model.Comment, error) {
	var newCom douyin_core_model.Comment
	newCom.Id = int64(comm.Id)
	newCom.User, _ = QueryUserByID(int(newCom.Id))
	newCom.Content = comm.Content
	newCom.Create_date = string(comm.Content[0]) + string(comm.Content[1]) + string(comm.Content[2]) + string(comm.Content[3]) + string(comm.Content[4])
	return newCom, nil
}

func QueryCommentByUserId(user_id, video_id int, date string) (douyin_core_model.Comment, error) {
	var comm model.UserComments
	var tp douyin_core_model.Comment
	if err := DB.Where("user_id = ? AND video_id = ? AND create_date = ?", user_id, video_id, date).First(&comm).Error; err != nil {
		return tp, err
	}
	return commentTrans(comm)
}

func QueryVideoComment(video_id int) ([]*douyin_core_model.Comment, error) {
	var res []*douyin_core_model.Comment
	var comms []*model.UserComments
	if err := DB.Where("video_id = ?", video_id).Find(&comms).Error; err != nil {
		return res, err
	}

	for _, comm := range comms {
		var newCom douyin_core_model.Comment

		newCom.Id = int64(comm.Id)
		newCom.User, _ = QueryUserByID(int(newCom.Id))
		newCom.Content = comm.Content
		newCom.Create_date = string(comm.Content[0]) + string(comm.Content[1]) + string(comm.Content[2]) + string(comm.Content[3]) + string(comm.Content[4])

		res = append(res, &newCom)
	}
	return res, nil
}
