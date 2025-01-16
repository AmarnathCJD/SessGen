package handler

import (
	"fmt"
	"net/http"
	"strconv"

	tg "github.com/amarnathcjd/gogram/telegram"
)

func SendCodeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	appId := r.Form.Get("appId")
	appHash := r.Form.Get("appHash")
	phone := r.Form.Get("phoneNumber")

	if appId == "" || appHash == "" || phone == "" {
		http.Error(w, `{"error":"appId, appHash and phoneNumber are required"}`, http.StatusBadRequest)
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

	hash, err := client.SendCode(phone)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"hash":"` + hash + `", "partial_session":"` + client.ExportSession() + `"}`))
}
