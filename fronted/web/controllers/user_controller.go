package controllers

import (
	"awesomeProject/datamodels"
	"awesomeProject/services"
	"awesomeProject/tool"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"strconv"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
	Session *sessions.Session
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (c *UserController) PostRegister() mvc.Response {
	user := new(datamodels.User)
	if err := c.Ctx.ReadForm(user); err != nil {
		c.Ctx.Application().Logger().Errorf("bind form to User failed: %v", err)
		return mvc.Response{Path: "/user/error"}
	}
	if _, err := c.Service.AddUser(user); err != nil {
		return mvc.Response{Path: "/user/error"}
	}
	return mvc.Response{Path: "/user/login"}
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response {
	//1.获取用户提交的表单信息
	user := new(datamodels.User)
	if err := c.Ctx.ReadForm(user); err != nil {
		c.Ctx.Application().Logger().Errorf("bind form to User failed: %v", err)
	}
	//2、验证账号密码正确
	if _, isOk := c.Service.IsPwdSuccess(user.UserName, user.HashPassword); isOk {
		return mvc.Response{Path: "/user/login"}
	}

	//3、写入用户ID到cookie中
	tool.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	c.Session.Set("userID", strconv.FormatInt(user.ID, 10))

	return mvc.Response{
		Path: "/product/",
	}

}
