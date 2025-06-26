package main

import (
	"awesomeProject/backend/web/controllers"
	"awesomeProject/common"
	"awesomeProject/repositories"
	"awesomeProject/services"
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	// 一、View
	// 声明：html文件都在对应目录下
	template := iris.HTML("./backend/web/views", ".html").
		Layout("shared/layout.html"). // 声明全局布局,前端常用的“母版页”
		Reload(true)                  // 文件有变动会自动重新加载模板
	app.RegisterView(template) //注册该模板引擎到 iris 应用

	// 静态资源：浏览器能直接访问静态资源
	// 路由映射：Web向本地映射： http://localhost:8080/assets/ <-> ./backend/web/assets/
	app.HandleDir("/assets", "./backend/web/assets")

	// 处理异常： 404, 500
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "<UNK>"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	// DataBase
	db, err := common.NewMysqlConn()
	if err != nil {
		panic(err)
	}

	// 上下文
	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	// 注册控制器:product
	productRepo := repositories.NewProductManager(db)         // 数据层操作
	productService := services.NewProductService(productRepo) // 服务层操作
	productParty := app.Party("/product")                     // 创建一个路由分组，product开头的HTTP路径，会归这个分组
	product := mvc.New(productParty)                          // MVC实例，与对应分组绑定
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	// 注册控制器:order
	orderRepo := repositories.NewOrderManagerRepository(db)
	orderService := services.NewOrderService(orderRepo)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	// 三、启动
	err = app.Run(
		iris.Addr("localhost:9999"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
	if err != nil {
		return
	}
}
