package main

import (
	"context"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"imooc-product/backend/web/controllers"
	"imooc-product/common"
	"imooc-product/repositories"
	"imooc-product/services"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	tmplate := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板目标
	app.HandleDir("/assets", "./backend/web/assets")
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error("err: ", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//5.注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productSerivce := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productSerivce)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//6.启动服务
	err1 := app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
	if err1 != nil {
		log.Error("err1:", err1)
	}

}
