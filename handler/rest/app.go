package rest

import (
	"fmt"
	"vocatrueid/helpers"
	"vocatrueid/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var port = ":" + helpers.GetEnv("PORT")

func StartApp() {
	var l *launcher.Launcher
	var browser *rod.Browser
	if helpers.GetEnv("mode") == "production" {
		l = launcher.MustNewManaged("")
		l.Headless(true).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")
		browser = rod.New().Client(l.MustClient()).MustConnect()
		defer l.Cleanup()
		defer browser.MustClose()
	} else {
		l = launcher.New().Headless(false)
		url := l.MustLaunch()
		browser = rod.New().ControlURL(url).MustConnect()
		defer l.Cleanup()
		defer browser.MustClose()
	}

	route := gin.Default()

	utilsRepository := utils.NewRepository(browser)
	utilsService := utils.NewService(utilsRepository)
	utilsHandler := NewService(utilsService)

	route.POST("/checknickname/:region/:code", utilsHandler.CheckNickname)

	fmt.Println("Server running on PORT =>", port)
	route.Run(port)
}
