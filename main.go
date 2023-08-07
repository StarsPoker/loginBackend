package main

import (
	"github.com/StarsPoker/loginBackend/app"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	app.StartApplication()
}
