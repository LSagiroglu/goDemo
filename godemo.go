package godemo

import "github.com/kataras/iris"

const Version string = "v1.0"

func HomePage(ctx iris.Context) {
	ctx.Application().Logger().Info("Request path: " + ctx.Path())
	Ip := ctx.RemoteAddr()
	ctx.HTML("<b>Go APP !</b><br/> Ip : " + Ip)

}

func InfoPage(ctx iris.Context) {
	ctx.Application().Logger().Info("Request path: " + ctx.Path())
	ctx.HTML("<b>Go APP Info!</b><br/> set DOMAIN and EMAIL environment variables and restart.")

}

func PingPage(ctx iris.Context) {
	Ip := ctx.RemoteAddr()
	ctx.HTML("<b>Go APP !</b><br/> Ip : " + Ip)

}
