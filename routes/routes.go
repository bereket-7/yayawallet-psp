package api

import (
	"yayawallet-psp/config"
	"yayawallet-psp/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
    s := service.NewYayaService(cfg)
    v1 := r.Group("/")
    {
        v1.POST("webhook", s.HandleWebhook)
        v1.POST("pay", s.CreatePaymentIntent)
    }
}
