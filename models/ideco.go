package models

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// Sharepoint Model ...
type IDECOModel struct{}

type WList struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Urls        []string `json:"urls"`
}

// GLPI User find by Mail ...
func (m IDECOModel) GetWiteList() (wlist []WList, err error) {
	urls1 := "https://iutm.brnv.rw/web/auth/login"
	urls2 := "https://iutm.brnv.rw/content-filter/users_categories"

	var jsonStr = []byte(`{ "login": "` + os.Getenv("IDECO_SERVER_LOGIN") + `", "password": "` + os.Getenv("IDECO_SERVER_PASS") + `", "rest_path": "/" }`)
	req, err := http.NewRequest("POST", urls1, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wlist, err
	}

	defer resp.Body.Close()
	cookies := resp.Cookies()

	req2, err := http.NewRequest("GET", urls2, nil)
	req2.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		return wlist, err
	}
	for _, v := range cookies {
		req2.AddCookie(&http.Cookie{Name: v.Name, Value: v.Value})
	}

	resp2, err := client.Do(req2)
	if err != nil {
		return wlist, err
	}
	defer resp2.Body.Close()
	body, _ := io.ReadAll(resp2.Body)
	err = json.Unmarshal(body, &wlist)
	if err != nil {
		return wlist, err
	}

	return wlist, err
}
