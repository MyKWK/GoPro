package datamodels

type User struct {
	ID           int64  `json:"id" form:"ID" sql:"ID" gorm:"column:ID;primary_key"`
	NickName     string `json:"nickName" form:"nickName" sql:"nickName" gorm:"column:nickName"`
	UserName     string `json:"userName" form:"userName" sql:"userName" gorm:"column:userName"`
	HashPassword string `json:"-" form:"passWord" sql:"passWord" gorm:"column:passWord"`
}
