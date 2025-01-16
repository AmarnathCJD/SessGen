package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	tg "github.com/amarnathcjd/gogram/telegram"
)

func HandleCodeLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	appId := r.Form.Get("appId")
	appHash := r.Form.Get("appHash")
	phone := r.Form.Get("phoneNumber")
	code := r.Form.Get("code")
	codeHash := r.Form.Get("codeHash")

	if appId == "" || appHash == "" || phone == "" || code == "" || codeHash == "" {
		http.Error(w, `{"error":"appId, appHash, phoneNumber, code and codeHash are required"}`, http.StatusBadRequest)
		return
	}

	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		http.Error(w, `{"error":"appId must be a number"}`, http.StatusBadRequest)
		return
	}

	client, err := tg.NewClient(tg.ClientConfig{
		AppID:         int32(appIdInt),
		AppHash:       appHash,
		MemorySession: true,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	err = client.Connect()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	signin, err := client.AuthSignIn(phone, codeHash, code, nil)
	if err != nil {
		if strings.Contains(err.Error(), "PHONE_CODE_INVALID") {
			http.Error(w, `{"error":"invalid code"}`, http.StatusBadRequest)
			return
		} else if strings.Contains(err.Error(), "SESSION_PASSWORD_NEEDED") {
			http.Error(w, `{"error":"password needed", "partial_session":"`+client.ExportSession()+`"}`, http.StatusBadRequest)
			return
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	}

	switch sg := signin.(type) {
	case *tg.AuthAuthorizationObj:
		w.Write([]byte(`{"session":"` + client.ExportSession() + `", "user":` + client.JSON(sg.User) + `}`))
	}

	w.Write([]byte(`{"session":"` + client.ExportSession() + `"}`))
}
