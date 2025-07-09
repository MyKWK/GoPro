package repositories

import (
	"awesomeProject/common"
	"awesomeProject/datamodels"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
}

func NewUserRepository(table string, db *gorm.DB) IUserRepository {
	return &UserManagerRepository{db}
}

type UserManagerRepository struct {
	mysqlConn *gorm.DB
}

func (u *UserManagerRepository) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = mysql
	}
	return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return nil, errors.New("条件不能为空！")
	}
	if err = u.Conn(); err != nil {
		return nil, err
	}
	user = &datamodels.User{}
	result := u.mysqlConn.Where("userName = ?", userName).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在！")
		}
		return nil, result.Error
	}
	return user, nil
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		fmt.Println("数据库没连上")
		return
	}
	fmt.Println("正在新建用户信息")
	if err = u.mysqlConn.Create(user).Error; err != nil {
		fmt.Println("插入失败了：")
		fmt.Println(err)
		return
	} else {
		userId = user.ID
	}
	fmt.Println("新用户注册成功")
	return
}

func (u *UserManagerRepository) SelectByID(userId int64) (user *datamodels.User, err error) {
	// ensure connection
	if err = u.Conn(); err != nil {
		return nil, err
	}

	user = &datamodels.User{}
	// use GORM to fetch by primary key
	result := u.mysqlConn.First(user, userId)
	// 获取row，传入结构体中
	// 根据user来指定表结构
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在！")
		}
		return nil, result.Error
	}
	return user, nil
}
