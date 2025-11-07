package main

import (
	"log"
	"yayawallet-psp/config"
	api "yayawallet-psp/routes"

	"github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    r := gin.Default()
    api.RegisterRoutes(r, cfg)
    log.Printf("yayawallet-psp running on %s", cfg.Port)
    r.Run(":" + cfg.Port)
}