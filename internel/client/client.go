package client

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/big"
	"net/http"
)

type MessageBody struct {
	MobileNetworkType string `json:"mobile_network_type"`
	Content           string `json:"content"`
	Nonce             string `json:"nonce"`
	TTS               bool   `json:"tts"`
	Flags             int    `json:"flags"`
}

func generateNumericNonce(length int) (string, error) {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func RequestDcMessage(authorization string, content string) error {
	url := "https://discord.com/api/v9/channels/1255514876072427582/messages"
	nonce, err := generateNumericNonce(19)
	data := MessageBody{
		MobileNetworkType: "unknown",
		Content:           "$request " + content,
		Nonce:             nonce,
		TTS:               false,
		Flags:             0,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error marshalling data: ", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating request: ", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body: ", "err", err)
			return fmt.Errorf("received non-OK response code %d with no body", resp.StatusCode)
		}
		slog.Error("Error reading response body: ", "body", bodyBytes)
		return fmt.Errorf("received non-OK response code %d with body: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
	// Handle response here...
}
