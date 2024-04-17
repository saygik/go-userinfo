package glpiapi

import "net/http"

func (r *Repository) SetSessionHeader(req *http.Request, token string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Session-Token", token)
	req.Header.Set("App-Token", r.token)
}
