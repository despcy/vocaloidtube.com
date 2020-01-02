package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/didip/tollbooth"
	"github.com/iris-contrib/middleware/csrf"
	"github.com/iris-contrib/middleware/secure"
	"github.com/iris-contrib/middleware/tollboothic"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/host"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"github.com/kataras/neffos"
	"github.com/yangchenxi/VOCALOIDTube/config"
	"github.com/yangchenxi/VOCALOIDTube/controller"
	"github.com/yangchenxi/VOCALOIDTube/model"
	"github.com/yangchenxi/VOCALOIDTube/model/authentication"
	"github.com/yangchenxi/VOCALOIDTube/security"
)

func main() {
	f := newLogFile()
	defer f.Close()
	app := iris.New()

	app.HandleDir("/assets", "./view/templates/assets")
	hometmp := iris.HTML("view/templates", ".html")
	app.Logger().SetLevel("info") //TODO:comment when release
	app.Logger().SetOutput(f)     //TODO:uncomment when release
	app.Logger().SetTimeFormat("2006-01-02 15:04:05.999999")
	hometmp.Reload(true) //TODO: turn false when release
	app.RegisterView(hometmp)
	app.OnErrorCode(404, notFound)
	limiter := tollbooth.NewLimiter(config.RateLimit, nil)
	protect := csrf.Protect([]byte(config.CSRFKEY), csrf.Secure(false)) // Defaults to true, but pass `false` while no https (devmode).
	//Now this web only has get method, so csrf has no effect here,only a 摆设,面试用
	s := secure.New(secure.Options{
		BrowserXSSFilter: true, // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		IsDevelopment:    true, // TODO:Change to false
	})
	app.Use(s.Serve)
	app.UseGlobal(beforeMiddleware)
	app.DoneGlobal(afterMiddleware)
	sessManager := authentication.NewSessionsManager()
	app.Use(sessManager.Handler())
	mvc.Configure(app.Party("/", protect, tollboothic.LimitHandler(limiter)), videos)
	mvc.Configure(app.Party("/auth"), auth)
	mvc.Configure(app.Party("/status", tollboothic.LimitHandler(limiter)), status)
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: controller.NativeMessageHandler,
	})
	ws.OnConnect = controller.ConnectHandler
	ws.OnDisconnect = controller.DisconnectHandler
	controller.WebsocketChan = make(chan []byte, 1000) //TODO: determine the approperate size of the channel buffer
	go webSocketBroadCastListener(ws)
	app.Get("/websocket", websocket.Handler(ws))
	// app.Get("/websoctest", func(ctx iris.Context) { //TODO:only for test purpose
	// 	ctx.View("wsclient.html")
	// })
	//websocAPI.Get("/", websocket.Handler(ws))

	// app.Run(
	// 	// Start the web server at localhost:8080
	// 	iris.Addr(":8080"),
	// 	// skip err server closed when CTRL/CMD+C pressed:
	// 	iris.WithoutServerError(iris.ErrServerClosed),
	// 	// enables faster json serialization and more:
	// 	iris.WithOptimizations,
	// )
	target, _ := url.Parse("https://vocaloidtube.com:443")
	go host.NewRedirection(":80", target, iris.StatusMovedPermanently).ListenAndServe()
	app.Run(iris.TLS(":443", "server.crt", "server.key"))

}

func beforeMiddleware(ctx iris.Context) {
	//Rate limit
	//block ip

	log.Printf("middleware" + ctx.Path())
	ip := ctx.RemoteAddr()
	if val, ok := security.BlackList[ip]; ok {

		ctx.StatusCode(403)
		ctx.Write([]byte("your ip has been blocked since " + val.TimeBlocked + " because illegal attack is detected"))
		ctx.EndRequest()
		return
	}
	security.RequestIpAntiRobot(ip, ctx.Path())
	ctx.Next()
}

func afterMiddleware(ctx iris.Context) {
	log.Printf("after middleware" + ctx.RemoteAddr())
	logMiddleware(ctx)
}

func webSocketBroadCastListener(ws *neffos.Server) {
	for {
		select {
		case c := <-controller.WebsocketChan:

			ws.Broadcast(nil, neffos.DeserializeMessage(nil, c, true, false))
		}
	}
}

func status(app *mvc.Application) {
	userDBService := model.NewUserDBService()
	app.Register(userDBService)
	app.Handle(new(controller.StatusController))
}

func auth(app *mvc.Application) {

	authService := model.NewAuthService()
	app.Register(authService)
	app.Handle(new(controller.AuthController))
}

func videos(app *mvc.Application) {

	videoDBService := model.NewVideoDBService()
	app.Register(videoDBService)
	//code before handle are executed when init the whole program
	app.Handle(new(controller.VideoController))

}
func notFound(ctx iris.Context) {

	ctx.StatusCode(404)

	logMiddleware(ctx)
	ctx.View("page-404.html")
}

//------------log 小本本拉清单

func logMiddleware(ctx iris.Context) {
	defer ctx.Next()
	ip := ctx.RemoteAddr()
	path := ctx.RequestPath(false)
	Referer := ctx.GetReferrer().URL
	UserAgent := ctx.GetHeader("User-Agent")
	RespCode := strconv.Itoa(ctx.GetStatusCode())
	Location := security.GetLocationFromIP(ip)

	ctx.Application().Logger().Infof("Request ip: %s |Location: %s| Request url: %s| Referer: %s | UserAgent: %s| RespCode: %s", ip, Location, path, Referer, UserAgent, RespCode)
	var secritylevel string
	if RespCode == "200" || RespCode == "304" {
		secritylevel = "Pass"
	} else if RespCode == "403" {
		secritylevel = "Blocked"
	} else {
		secritylevel = "Warning"
	}
	model.VisitCountMutex.Lock()
	model.VisitCount++
	model.VisitCountMutex.Unlock()
	t := time.Now()

	controller.SendWebSocData(strconv.FormatUint(model.VisitCount, 10), t.Format("15:04:05"), security.FormatDisplayIP(ip), Location, path, secritylevel)

}

func newLogFile() *os.File {
	filename := time.Now().Format("Jan 02 2006") + "serverlog.txt"
	// Open the file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}
