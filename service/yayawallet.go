package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

func (s *YayaService) CreatePaymentIntent(c *gin.Context) {
    var req model.PaymentIntentRequest

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
        return
    }

    payload, _ := json.Marshal(req)

    yayaReq, _ := http.NewRequest("POST",
        s.cfg.YayaBaseURL + "/payment-intent",
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

func (s *YayaService) HandleWebhook(c *gin.Context) {
    signature := c.GetHeader("X-Payment-Signature")
    if signature == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing signature"})
        return
    }

    var payload model.YayaCallbackRequest
    if err := c.BindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    // Step 1: Concatenate paymentId + paymentReference + amount
    raw := payload.PaymentId + payload.PaymentReference + formatAmount(payload.Amount)

    // Step 2: Generate HMAC SHA-256
    expected := computeHMAC(raw, s.cfg.ClientSecret)

    // Step 3: Compare signatures
    if !hmac.Equal([]byte(signature), []byte(expected)) {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error":       "invalid signature",
            "received":    signature,
            "expected":    expected,
            "stringified": raw,
        })
        return
    }

    // TODO: Store or forward to FenanPay
    // TODO: Implement idempotency using transactionId or paymentId

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func formatAmount(a float64) string {
    return fmt.Sprintf("%.0f", a) 
}

func computeHMAC(message, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return hex.EncodeToString(h.Sum(nil))
}