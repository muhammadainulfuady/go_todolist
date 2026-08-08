package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const resendEndpoint = "https://api.resend.com/emails"

type Mailer struct {
	apiKey string
	client *http.Client
}

type resendRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
}

func NewMailer(apiKey string) *Mailer {
	return &Mailer{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *Mailer) SendOTP(to, code string) error {
	payload := resendRequest{
		From:    "onboarding@resend.dev",
		To:      to,
		Subject: "Kode OTP Todo List",
		HTML:    fmt.Sprintf("<p>Kode OTP kamu adalah: <b>%s</b>. Kode berlaku 30 detik.</p>", code),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		logrus.WithField("detail", err.Error()).Error("Gagal membuat payload email OTP")
		return fmt.Errorf("gagal membuat payload email OTP: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resendEndpoint, bytes.NewReader(body))
	if err != nil {
		logrus.WithField("detail", err.Error()).Error("Gagal membuat request email OTP")
		return fmt.Errorf("gagal membuat request email OTP: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-todolist/1.0")

	resp, err := m.client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"to":     to,
			"detail": err.Error(),
		}).Error("Gagal terhubung ke API Resend")
		return fmt.Errorf("gagal mengirim email OTP: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithField("detail", err.Error()).Error("Gagal membaca response Resend")
		return fmt.Errorf("gagal membaca response Resend: %w", err)
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		logrus.WithFields(logrus.Fields{
			"to":     to,
			"status": resp.StatusCode,
		}).Info("Email OTP berhasil dikirim")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"to":       to,
		"status":   resp.StatusCode,
		"response": string(respBody),
	}).Error("Gagal mengirim email OTP via Resend")
	return fmt.Errorf("gagal mengirim email OTP: resend status %d: %s", resp.StatusCode, string(respBody))
}
