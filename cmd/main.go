package main

import (
	"log"
	"net/http"
	"onesignal_api/keitaro"
	"os"
)

func main() {
	InitLogs()
	routes := keitaro.InitRoutes()
	log.Fatal(http.ListenAndServe(":8125", routes))
}

func InitLogs() {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}
