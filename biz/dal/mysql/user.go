package mysql

import (
	"douyin-demo/biz/model"
	"douyin-demo/biz/model/douyin_core_model"
)

func CreateUsers(users []*model.Users) error {
	return DB.Create(users).Error
}

func CreateUser(users []*model.Users) error {
	return DB.Create(users).Error
}

func DeleteUser(userId int64) error {
	return DB.Where("id = ?", userId).Delete(&model.Users{}).Error
}

func UpdateUser(user *model.Users) error {
	return DB.Updates(user).Error
}

func FindUserByName(userName, email string) ([]*model.Users, error) {
	res := make([]*model.Users, 0)
	if err := DB.Where("name = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func CheckUser(username, password string) (*model.Users, error) {
	var res *model.Users
	if err := DB.Where("name = ?", username).Where("passwd_hash = ?", password).First(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func QueryUserId(username string) (int64, error) {
	var res model.Users
	if len(username) != 0 {
		if err := DB.Where("name = ?", username).First(&res).Error; err != nil {
			return -1, err
		}
	}
	return int64(res.Id), nil
}

func modelUsers2CoreUsers(user model.Users) (douyin_core_model.User, error) {
	var newUsr douyin_core_model.User
	newUsr.Id = int64(user.Id)
	newUsr.FavoriteCount = user.FavoriteCount
	newUsr.Name = user.Name
	newUsr.TotalFavorited = user.TotalFavorited
	newUsr.WorkCount = user.WorkCount
	return newUsr, nil
}

func QueryUserByID(user_id int) (douyin_core_model.User, error) {
	var res model.Users
	var tp douyin_core_model.User
	if user_id > 0 {
		if err := DB.Where("id = ?", user_id).First(&res).Error; err != nil {
			return tp, err
		}
	}
	return modelUsers2CoreUsers(res)
}

func AddUserWork(user_id int) error {
	var user model.Users
	if err := DB.Where("id = ?", user_id).First(&user).Error; err != nil {
		return err
	}
	user.WorkCount++
	DB.Save(&user)
	return nil
}
