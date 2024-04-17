package glpiapi

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/saygik/go-userinfo/internal/entity"
)

// Create comment ...
func (r *Repository) CreateSolution(form entity.NewCommentInputForm) (commentId int, err error) {

	var token string

	if len(form.Input.Token) == 0 {
		token, err = r.GLPIInitSession()
	} else {
		token, err = r.GLPIInitSessionUser(form.Input.Token)
	}

	if err != nil {
		return 0, errors.New(err.Error())
	}
	defer r.KillSession(token)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(form)

	return r.AddItem("itilsolution/", token, payloadBuf, "Невозможно создать решение GLPI")

}
