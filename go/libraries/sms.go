package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var apiKey = "01ZtycjQXUQFlarNuAVGMRmaJHFQilUrKSbGeIUBaeD2ZI6Q"

type smsPayload struct {
	Mobile     string `json:"mobile"`
	TemplateId int    `json:"templateId"`
	Parameters []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"parameters"`
}

func SendSmsLogin(code, to string) (bool, error) {
	if os.Getenv("SMS_SANDBOX") == "true" {
		fmt.Printf("📨 [SANDBOX] Fake SMS sent to %s with code %s\n", to, code)
		return true, nil
	}
	payload := smsPayload{
		Mobile:     to,
		TemplateId: 763111,
		Parameters: []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{
			{Name: "CODE", Value: code},
		},
	}
	body, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", "https://api.sms.ir/v1/send/verify", bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("x-api-key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}
	if status, ok := res["status"].(float64); ok && int(status) == 1 {
		return true, nil
	}
	fmt.Println("SMS API response:", res)
	return false, nil
}
