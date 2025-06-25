package datamodels

type Product struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id" imooc:"ID" sql:"id"`
	ProductName  string `gorm:"column:ProductName;type:varchar(255)" json:"ProductName" imooc:"ProductName" sql:"ProductName"`
	ProductNum   int64  `gorm:"column:ProductNum" json:"ProductNum" imooc:"ProductNum" sql:"ProductNum"`
	ProductImage string `gorm:"column:ProductImage;type:varchar(512)" json:"ProductImage" imooc:"ProductImage" sql:"ProductImage"`
	ProductUrl   string `gorm:"column:ProductUrl;type:varchar(512)" json:"ProductUrl" imooc:"ProductUrl" sql:"ProductUrl"`
}
