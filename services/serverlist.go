package services

import (
	"encoding/json"
	"os"
	"strings"

	"maxion-zone4/config"
)

type ServerEntry struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

type ServerList struct {
	Servers []ServerEntry `json:"servers"`
}

// SERVER_IP no .env define o IP enviado ao client (default 127.0.0.1)
func serverIP() string {
	ip := strings.TrimSpace(os.Getenv("SERVER_IP"))
	if ip == "" {
		ip = "127.0.0.1"
	}
	return ip
}

func BuildServerList() ServerList {
	ip := serverIP()
	ports := config.ChannelPorts()

	list := ServerList{Servers: make([]ServerEntry, 0, len(ports))}
	for _, p := range ports {
		name := "Canal"
		if p == 9090 {
			name = "Canal Principal"
		}
		list.Servers = append(list.Servers, ServerEntry{
			ID:     name,
			Name:   name,
			IP:     ip,
			Port:   p,
			Status: "online",
		})
	}
	return list
}

func BuildServerListJSON() []byte {
	b, _ := json.Marshal(BuildServerList())
	return b
}
