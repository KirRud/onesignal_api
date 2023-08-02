package keitaro

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"onesignal_api/onesignal"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()
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
