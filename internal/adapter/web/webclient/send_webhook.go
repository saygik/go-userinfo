package webclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func (r *Repository) sendWebhook(data interface{}, url string) error {
	// Marshal the data into JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Prepare the webhook request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the webhook request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("Error closing response body:", err)
		}
	}(resp.Body)

	// Determine the status based on the response code
	status := "webhook failed"
	if resp.StatusCode == http.StatusOK {
		status = "webhook delivered"
	}

	r.log.Info(status)

	if status == "webhook failed" {
		return errors.New(status)
	}

	return nil
}
