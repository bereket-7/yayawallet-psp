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
func (s *YayaService) getAccessToken() (string, error) {
    payload := map[string]string{
        "client_id":     s.cfg.YayaClientID,
        "client_secret": s.cfg.YayaClientSecret,
        "grant_type":    "client_credentials",
    }

    body, _ := json.Marshal(payload)

    req, err := http.NewRequest(
        "POST",
        s.cfg.YayaBaseURL+"/api/auth/token",
        bytes.NewReader(body),
    )
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("auth failed: %s", string(b))
    }

    var tokenResp struct {
        AccessToken string `json:"access_token"`
        ExpiresIn   int    `json:"expires_in"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
        return "", err
    }

    return tokenResp.AccessToken, nil
}

func (s *YayaService) CreatePaymentIntent(c *gin.Context) {
    var req model.PaymentIntentRequest

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := s.getAccessToken()
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    payload, _ := json.Marshal(req)

    yayaReq, _ := http.NewRequest(
        "POST",
        s.cfg.YayaBaseURL+"/api/payment-intent",
        bytes.NewReader(payload),
    )

    yayaReq.Header.Set("Authorization", "Bearer "+token)
    yayaReq.Header.Set("Content-Type", "application/json")
    yayaReq.Header.Set("Accept", "application/json")
    yayaReq.Header.Set("User-Agent", "yayawallet-psp/1.0")

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

    raw := payload.PaymentId + payload.PaymentReference + formatAmount(payload.Amount)

    expected := computeHMAC(raw, s.cfg.YayaClientSecret)

    if !hmac.Equal([]byte(signature), []byte(expected)) {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error":       "invalid signature",
            "received":    signature,
            "expected":    expected,
            "stringified": raw,
        })
        return
    }

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