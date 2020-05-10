package app

import (
	"github.com/StarsPoker/loginBackend/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	logger.Info("about to start the application ...")
	mapUrls()
	router.Run(":8080")
}
