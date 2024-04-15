package main

import (
	"flag"
	"log"
	"net/http"
	"notify/pkg/application"
	"notify/pkg/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

var PathToConf = flag.String("path", "./config.yaml", "Path to config file")

func main() {
	flag.Parse()
	app := application.NewApp(*PathToConf)

	defer app.Close()

	if app.Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, *model.Pong())
	})

	router.POST("/send/:action", func(c *gin.Context) {
		action := c.Param("action")
		req := make(model.Request)
		c.BindJSON(&req)
		// app.Send(action, req)
		if err := app.SendMessage(action, req); err != nil {
			c.JSON(http.StatusBadRequest, *model.Error(err))
			return
		}

		c.JSON(http.StatusOK, *model.Success())
	})

	log.Fatal(router.Run(":" + strconv.Itoa(app.Config.Port)))
}
