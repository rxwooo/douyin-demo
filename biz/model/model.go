package model

// User information table
type Users struct {
	Id              uint64 `gorm:"column:id" db:"id" json:"id" form:"id"`
	Name            string `gorm:"column:name" db:"name" json:"name" form:"name"`
	PasswdHash      string `gorm:"column:passwd_hash" db:"passwd_hash" json:"passwd_hash" form:"passwd_hash"`
	Signature       string `gorm:"column:signature" db:"signature" json:"signature" form:"signature"`
	FollowCount     uint64 `gorm:"column:follow_count" db:"follow_count" json:"follow_count" form:"follow_count"`
	FollowerCount   uint64 `gorm:"column:follower_count" db:"follower_count" json:"follower_count" form:"follower_count"`
	Avatar          string `gorm:"column:avatar" db:"avatar" json:"avatar" form:"avatar"`
	BackgroundImage string `gorm:"column:background_image" db:"background_image" json:"background_image" form:"background_image"`
	TotalFavorited  int64  `gorm:"column:total_favorited" db:"total_favorited" json:"total_favorited" form:"total_favorited"`
	WorkCount       int64  `gorm:"column:work_count" db:"work_count" json:"work_count" form:"work_count"`
	FavoriteCount   int64  `gorm:"column:favorite_count" db:"favorite_count" json:"favorite_count" form:"favorite_count"`
}

type Videos struct {
	Id            uint64 `gorm:"column:id" db:"id" json:"id" form:"id"`                                                 //auto increment id
	UserId        int64  `gorm:"column:user_id" db:"user_id" json:"user_id" form:"user_id"`                             //User id
	PlayUrl       string `gorm:"column:play_url" db:"play_url" json:"play_url" form:"play_url"`                         //video src
	CoverUrl      string `gorm:"column:cover_url" db:"cover_url" json:"cover_url" form:"cover_url"`                     //cover path
	FavoriteCount int64  `gorm:"column:favorite_count" db:"favorite_count" json:"favorite_count" form:"favorite_count"` //video favorite count
	CommentCount  int64  `gorm:"column:comment_count" db:"comment_count" json:"comment_count" form:"comment_count"`     //video comment count
	Title         string `gorm:"column:title" db:"title" json:"title" form:"title"`                                     //video title
}

type Favorites struct {
	Id      uint64 `gorm:"column:id" db:"id" json:"id" form:"id"`
	UserId  int64  `gorm:"column:user_id" db:"user_id" json:"user_id" form:"user_id"`     //user id
	VideoId int64  `gorm:"column:video_id" db:"video_id" json:"video_id" form:"video_id"` //followed user ID
}

// User comments table
type UserComments struct {
	Id uint64 `gorm:"column:id" db:"id" json:"id" form:"id"`
	UserId int64 `gorm:"column:user_id" db:"user_id" json:"user_id" form:"user_id"` //user id
	VideoId int64 `gorm:"column:video_id" db:"video_id" json:"video_id" form:"video_id"` //video ID
	Content string `gorm:"column:content" db:"content" json:"content" form:"content"`
	CreateDate string `gorm:"column:create_date" db:"create_date" json:"create_date" form:"create_date"`
}