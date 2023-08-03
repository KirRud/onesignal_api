package keitaro

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"onesignal_api/onesignal"
	"runtime"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(panicRecovery)
	routes := r.PathPrefix("/postback").Subrouter()
	routes.HandleFunc("", handlePostback).Methods("GET")

	return r
}

func handlePostback(w http.ResponseWriter, r *http.Request) {
	// Разбираем URL запроса
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем значения параметров
	queryParams := u.Query()
	appID := queryParams.Get("app_id")
	userID := queryParams.Get("external_user_id")
	data := queryParams.Get("tags")

	var temp map[string]interface{}
	json.Unmarshal([]byte(data), &temp)
	newT := temp["status"].(int)
	log.Printf("Check %v", newT)

	// Отправляем в ответе
	err, status := sendToOneSignal(data, appID, userID)
	if err != nil {
		log.Printf("Error while updated: %v", err)
	}
	log.Printf("Status:%d app_id:%s external_user_id:%s tags:%s", status, appID, userID, data)
	w.Write([]byte(http.StatusText(status)))
}

func sendToOneSignal(tags, appID, userID string) (error, int) {
	client := http.Client{}
	osConn := onesignal.NewOneSignal(&client)

	err, status := osConn.UpdateUserTag(tags, appID, userID)
	return err, status
}

func panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				log.Printf("recovering from err: %v\non url: %s\n%s", err, r.URL, buf)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`server got panic`))
			}
		}()

		h.ServeHTTP(w, r)
	})
}
