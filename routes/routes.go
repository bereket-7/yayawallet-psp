package api

import (
	"yayawallet-psp/config"
	"yayawallet-psp/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
    s := service.NewYayaService(cfg)
    v1 := r.Group("/api/v1")
    {
        v1.POST("/payment", s.CreatePayment)
        v1.POST("/webhook/yaya", s.HandleWebhook)
        v1.POST("/payment-intent", s.CreatePaymentIntent)
    }
}
