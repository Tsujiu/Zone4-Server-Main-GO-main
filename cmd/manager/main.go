package main

import (
	"log"
	"net/http"

	mcfg "maxion-zone4/manager/config"
	"maxion-zone4/manager/controllers"
)

func main() {
	if err := mcfg.LoadChannels("manager/config/channels.json"); err != nil {
		log.Fatalf("Failed to load channels: %v", err)
	}

	http.HandleFunc("/start/", controllers.StartChannelHandler)
	http.HandleFunc("/stop/", controllers.StopChannelHandler)
	http.HandleFunc("/status", controllers.StatusHandler)
	http.HandleFunc("/start-all", controllers.StartAllHandler)
	http.HandleFunc("/stop-all", controllers.StopAllHandler)

	http.Handle("/logs/", http.StripPrefix("/logs/", http.FileServer(http.Dir("./logs"))))
	http.Handle("/", http.FileServer(http.Dir("./manager/web")))

	log.Println("Channel Manager running on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
