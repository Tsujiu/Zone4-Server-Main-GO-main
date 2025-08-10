package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Channel struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Port  int    `json:"port"`
	RunCmd string `json:"run_cmd"`
}

var (
	channelsMu sync.RWMutex
	channels   []Channel

	procMu   sync.RWMutex
	procMap  = map[string]*exec.Cmd{}
)

func loadChannelsJSON(path string) error {
	channelsMu.Lock()
	defer channelsMu.Unlock()

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var arr []Channel
	if err := json.Unmarshal(b, &arr); err != nil {
		return err
	}
	channels = arr
	return nil
}

func getChannels() []Channel {
	channelsMu.RLock()
	defer channelsMu.RUnlock()
	cp := make([]Channel, len(channels))
	copy(cp, channels)
	return cp
}
func getChannelByID(id string) *Channel {
	channelsMu.RLock()
	defer channelsMu.RUnlock()
	for _, c := range channels {
		if c.ID == id {
			cc := c
			return &cc
		}
	}
	return nil
}

// ---- process helpers ----

func isRunning(id string) bool {
	procMu.RLock()
	defer procMu.RUnlock()
	cmd, ok := procMap[id]
	if !ok || cmd == nil || cmd.Process == nil {
		return false
	}
	return true
}
func startProcess(id, rawCmd string) error {
	procMu.Lock()
	defer procMu.Unlock()

	if _, ok := procMap[id]; ok {
		return fmt.Errorf("already running")
	}
	// sh -c para permitir flags --port/--index
	cmd := exec.Command("sh", "-c", rawCmd)
	cmd.Dir = "." // raiz do projeto
	if err := cmd.Start(); err != nil {
		return err
	}
	procMap[id] = cmd
	go func(id string, cmd *exec.Cmd) {
		_ = cmd.Wait()
		procMu.Lock()
		delete(procMap, id)
		procMu.Unlock()
	}(id, cmd)
	return nil
}
func stopProcess(id string) error {
	procMu.Lock()
	defer procMu.Unlock()

	cmd, ok := procMap[id]
	if !ok || cmd == nil || cmd.Process == nil {
		return fmt.Errorf("not running")
	}
	err := cmd.Process.Kill()
	delete(procMap, id)
	return err
}
func stopAll() {
	procMu.Lock()
	defer procMu.Unlock()
	for id, cmd := range procMap {
		if cmd != nil && cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		delete(procMap, id)
	}
}

// ---- HTTP Handlers ----

func startHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/start/")
	ch := getChannelByID(id)
	if ch == nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}
	if err := startProcess(ch.ID, ch.RunCmd); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.Write([]byte("Started " + ch.ID))
}
func stopHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/stop/")
	if err := stopProcess(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Write([]byte("Stopped " + id))
}
func statusHandler(w http.ResponseWriter, r *http.Request) {
	type status struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Running bool   `json:"running"`
		Port    int    `json:"port"`
	}
	var result []status
	for _, ch := range getChannels() {
		result = append(result, status{
			ID: ch.ID, Name: ch.Name, Running: isRunning(ch.ID), Port: ch.Port,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
func startAllHandler(w http.ResponseWriter, r *http.Request) {
	for _, ch := range getChannels() {
		_ = startProcess(ch.ID, ch.RunCmd)
	}
	w.Write([]byte("Started all channels"))
}
func stopAllHandler(w http.ResponseWriter, r *http.Request) {
	stopAll()
	w.Write([]byte("Stopped all channels"))
}

// ---- Boot Manager ----

func runManager() {
	// carrega manager/config/channels.json
	cfg := filepath.Join("manager", "config", "channels.json")
	if err := loadChannelsJSON(cfg); err != nil {
		log.Fatalf("Failed to load %s: %v", cfg, err)
	}

	http.HandleFunc("/start/", startHandler)
	http.HandleFunc("/stop/", stopHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/start-all", startAllHandler)
	http.HandleFunc("/stop-all", stopAllHandler)

	// logs
	http.Handle("/logs/", http.StripPrefix("/logs/", http.FileServer(http.Dir("./logs"))))
	// painel web estático (opcional)
	http.Handle("/", http.FileServer(http.Dir("./manager/web")))

	log.Println("Channel Manager running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
