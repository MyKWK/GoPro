package main

import (
	"awesomeProject/common"
	"awesomeProject/fronted/middleware"
	"awesomeProject/fronted/web/controllers"
	"awesomeProject/repositories"
	"awesomeProject/services"
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

func main() {
	//1.创建iris 实例
	app := iris.New()

	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	template := iris.
		HTML("./fronted/web/views", ".html").
		Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	//4.设置模板，让远程路径public和本地public对应起来
	app.HandleDir("/public", "./fronted/web/public")
	// 同上
	app.HandleDir("/html", "./fronted/web/htmlProductShow")
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {

	}
	sess := sessions.New(sessions.Config{
		Cookie:  "AdminCookie",
		Expires: 600 * time.Minute,
	})
	app.UseRouter(sess.Handler())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册user
	user := repositories.NewUserRepository("user", db) // 封装了db操作
	userService := services.NewService(user)           // 封装了数据操作层
	userParty := mvc.New(app.Party("/user"))           // 设定分组 user开头
	userParty.Register(userService, ctx, sess)         // 注册
	userParty.Handle(new(controllers.UserController))  // 将一个控制器实例（UserController）注册到这个 MVC 应用上。

	//注册product控制器
	orderManager := repositories.NewOrderManagerRepository(db)
	orderService := services.NewOrderService(orderManager)

	product := repositories.NewProductManager(db)
	productService := services.NewProductService(product)
	productParty := app.Party("/product")
	productParty.Use(middleware.AuthConProduct) // 挂载中间件

	pro := mvc.New(productParty) //?
	pro.Register(productService, orderService)
	pro.Handle(new(controllers.ProductController))

	err = app.Run(
		iris.Addr("0.0.0.0:9998"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
	if err != nil {
		return
	}

}
