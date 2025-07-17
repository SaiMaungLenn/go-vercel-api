package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"success": false, "error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req ContactRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"success": false, "error": "Invalid body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"success": false, "error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Message == "" {
		http.Error(w, `{"success": false, "error": "Missing required fields"}`, http.StatusBadRequest)
		return
	}

	text := fmt.Sprintf("📩 New Contact Form Submission\n\n👤 Name: %s\n📧 Email: %s\n", req.Name, req.Email)
	if req.Subject != "" {
		text += fmt.Sprintf("📝 Subject: %s\n", req.Subject)
	}
	text += fmt.Sprintf("💬 Message:\n%s", req.Message)

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, `{"success": false, "error": "Failed to send to Telegram"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
