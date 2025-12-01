package iutm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) List() []entity.IutmCategoryList {
	var categories []entity.IutmCategoryList
	loginURL := r.url + "/web/auth/login"
	dataURL := r.url + "/content-filter/users_categories"
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// Авторизация
	loginJSON := []byte(`{"login": "userinfo-bot", "password": "UHg&-nMA?c1G", "rest_path": "/"}`)
	reqLogin, _ := http.NewRequest("POST", loginURL, bytes.NewBuffer(loginJSON))
	reqLogin.Header.Set("Content-Type", "application/json")

	respLogin, err := client.Do(reqLogin)
	if err != nil {
		return categories
	}
	defer respLogin.Body.Close()

	if respLogin.StatusCode != 200 {
		fmt.Println("Authorization failed with status:", respLogin.Status)
		return categories
	}

	// Отложенный logout, вызовется при выходе из main()
	defer func() {
		reqLogout, _ := http.NewRequest("DELETE", loginURL, nil)
		respLogout, errLogout := client.Do(reqLogout)
		if errLogout != nil {
			fmt.Println("Logout error:", errLogout)
			return
		}
		defer respLogout.Body.Close()
		if respLogout.StatusCode == 200 {
			//fmt.Println("Successfully logged out")
		} else {
			//fmt.Println("Logout failed with status:", respLogout.Status)
		}
	}()

	// Получение данных
	reqData, _ := http.NewRequest("GET", dataURL, nil)
	respData, errData := client.Do(reqData)
	if errData != nil {
		return categories
	}
	defer respData.Body.Close()

	body, _ := io.ReadAll(respData.Body)
	json.Unmarshal(body, &categories)

	return categories
}
