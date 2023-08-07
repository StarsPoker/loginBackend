package app

import (
	"os"

	"github.com/StarsPoker/loginBackend/logger"
	"github.com/gin-gonic/gin"
)

var (
	router    = gin.Default()
	stars_pem = os.Getenv(stars_pem_key)
	stars_key = os.Getenv(stars_key_key)
	use_cert  = os.Getenv(use_cert_key)
)

const (
	stars_pem_key = "stars_pem"
	stars_key_key = "stars_key"
	use_cert_key  = "use_cert"
)

func StartApplication() {
	logger.Info("about to start the application ...")
	mapUrls()

	if use_cert == "true" {
		router.RunTLS(":8079", stars_pem, stars_key)
	} else {
		router.Run(":8079")
	}
}
