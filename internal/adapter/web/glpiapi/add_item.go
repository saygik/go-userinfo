package glpiapi

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) AddItem(itemTypeUrl string, token string, body io.Reader, errMsg string) (ticketId int, err error) {
	ticketId = 0
	glpirest := r.url + itemTypeUrl
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	req, err := http.NewRequest("POST", glpirest, body)
	if err != nil {
		return ticketId, err
	}

	r.SetSessionHeader(req, token)

	resp, err := client.Do(req)
	if err != nil {
		return ticketId, err
	}
	defer resp.Body.Close()

	var glpiRes entity.GLPIApiResponse
	if resp.StatusCode == http.StatusCreated {
		err1 := json.NewDecoder(resp.Body).Decode(&glpiRes)
		fmt.Println(err1)
		if err1 != nil {
			return ticketId, err1
		}
		ticketId = glpiRes.Id
	} else {
		return 0, errors.New(errMsg)
	}
	return ticketId, err
}
