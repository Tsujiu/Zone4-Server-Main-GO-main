package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	mcfg "maxion-zone4/manager/config"
	"maxion-zone4/manager/process"
)

// ===== helpers =====

type apiResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func enableCORS(w http.ResponseWriter, r *http.Request) bool {
	// CORS básico para permitir chamadas do painel web local
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func allowMethods(w http.ResponseWriter, r *http.Request, methods ...string) bool {
	if r.Method == http.MethodOptions {
		return true
	}
	for _, m := range methods {
		if r.Method == m {
			return true
		}
	}
	writeJSON(w, http.StatusMethodNotAllowed, apiResp{OK: false, Message: "method not allowed"})
	return false
}

// ===== handlers =====

func StartChannelHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if !allowMethods(w, r, http.MethodPost, http.MethodGet) {
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/start/")
	channel := mcfg.GetChannelByID(id)
	if channel == nil {
		writeJSON(w, http.StatusNotFound, apiResp{OK: false, Message: "channel not found"})
		return
	}
	if err := process.StartProcess(channel.ID, channel.RunCmd); err != nil {
		writeJSON(w, http.StatusConflict, apiResp{OK: false, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResp{OK: true, Message: "started " + channel.ID})
}

func StopChannelHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if !allowMethods(w, r, http.MethodPost, http.MethodGet) {
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/stop/")
	if err := process.StopProcess(id); err != nil {
		writeJSON(w, http.StatusNotFound, apiResp{OK: false, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResp{OK: true, Message: "stopped " + id})
}

func RestartChannelHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if !allowMethods(w, r, http.MethodPost, http.MethodGet) {
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/restart/")
	_ = process.StopProcess(id)
	ch := mcfg.GetChannelByID(id)
	if ch == nil {
		writeJSON(w, http.StatusNotFound, apiResp{OK: false, Message: "channel not found"})
		return
	}
	if err := process.StartProcess(ch.ID, ch.RunCmd); err != nil {
		writeJSON(w, http.StatusConflict, apiResp{OK: false, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResp{OK: true, Message: "restarted " + id})
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if !allowMethods(w, r, http.MethodGet) {
		return
	}

	channels := mcfg.GetChannels()
	type status struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Running bool   `json:"running"`
		Port    int    `json:"port"`
	}
	var result []status
	for _, ch := range channels {
		result = append(result, status{
			ID:      ch.ID,
			Name:    ch.Name,
			Running: process.IsRunning(ch.ID),
			Port:    ch.Port,
		})
	}
	writeJSON(w, http.StatusOK, result)
}

func StartAllHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if !allowMethods(w, r, http.MethodPost, http.MethodGet) {
		return
	}

	for _, ch := range mcfg.GetChannels() {
		_ = process.StartProcess(ch.ID, ch.RunCmd)
	}
	writeJSON(w, http.StatusOK, apiResp{OK: true, Message: "started all channels"})
}

func StopAllHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if !allowMethods(w, r, http.MethodPost, http.MethodGet) {
		return
	}

	process.StopAll()
	writeJSON(w, http.StatusOK, apiResp{OK: true, Message: "stopped all channels"})
}
