package glpiapi

import (
	"crypto/tls"
	"net/http"
)

func (r *Repository) KillSession(token string) (err error) {
	glpirest := r.url + "killSession/"
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	req, err := http.NewRequest("GET", glpirest, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Session-Token", token)
	req.Header.Set("App-Token", r.token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
