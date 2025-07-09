package services

import (
	"awesomeProject/datamodels"
	"awesomeProject/repositories"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {

	user, err := u.UserRepository.Select(userName)
	fmt.Println(user)

	if err != nil {
		return
	}
	fmt.Println("开始验证密码是否正确")

	isOk, _ = ValidatePassword(pwd, user.HashPassword)

	if !isOk {
		return &datamodels.User{}, false
	}

	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	var pwdByte []byte
	var errPwd error
	if pwdByte, errPwd = GeneratePassword(user.HashPassword); errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	fmt.Println("正在准备插入用户信息")
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	fmt.Println("开始验证密码是否正确")
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil

}
