package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"yayawallet-psp/config"
	"yayawallet-psp/model"

	"github.com/gin-gonic/gin"
)

type YayaService struct {
    cfg *config.Config
}

func NewYayaService(cfg *config.Config) *YayaService {
    return &YayaService{cfg: cfg}
}

func (s *YayaService) CreatePayment(c *gin.Context) {
    var req map[string]interface{}
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    data, _ := json.Marshal(req)
    httpReq, _ := http.NewRequest("POST", s.cfg.YayaBaseURL+"/payments", bytes.NewReader(data))
    httpReq.Header.Add("Authorization", "Bearer "+s.cfg.YayaApiKey)
    httpReq.Header.Add("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(httpReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    c.Data(resp.StatusCode, "application/json", body)
}

func (s *YayaService) HandleWebhook(c *gin.Context) {
    var payload map[string]interface{}
    if err := c.BindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "received", "data": payload})
}


func (s *YayaService) CreatePaymentIntent(c *gin.Context) {
    var req model.PaymentIntentRequest

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
        return
    }

    payload, _ := json.Marshal(req)

    yayaReq, _ := http.NewRequest("POST",
        s.cfg.YayaBaseURL+"/api/payment-intent",
        bytes.NewReader(payload),
    )

    yayaReq.Header.Add("Authorization", "Bearer "+s.cfg.YayaApiKey)
    yayaReq.Header.Add("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(yayaReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    c.Data(resp.StatusCode, "application/json", body)
}
