package handler

import (
	"fmt"
	"net/http"
	"strings"

	tg "github.com/amarnathcjd/gogram/telegram"
)

func HandlePasswLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sess := r.Form.Get("partial_session")
	pass := r.Form.Get("password")

	if sess == "" || pass == "" {
		http.Error(w, `{"error":"partial_session and password are required"}`, http.StatusBadRequest)
		return
	}

	client, err := tg.NewClient(tg.ClientConfig{
		StringSession: sess,
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

	accPassword, err := client.AccountGetPassword()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	inputPassword, err := tg.GetInputCheckPassword(pass, accPassword)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = client.AuthCheckPassword(inputPassword)
	if err != nil {
		if strings.Contains(err.Error(), "PASSWORD_HASH_INVALID") {
			http.Error(w, `{"error":"invalid password"}`, http.StatusBadRequest)
			return
		}

		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"session":"` + client.ExportSession() + `"}`))
}
