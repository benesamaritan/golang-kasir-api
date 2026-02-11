package handlers

import (
	"encoding/json"
	"net/http"
)

type RouteInfo struct {
	Path string `json:"path"`
}

type APIHandler struct {
	Routes *[]RouteInfo
}

func NewAPIHandler(routes *[]RouteInfo) *APIHandler {
	return &APIHandler{Routes: routes}
}

func (h *APIHandler) ListRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Routes)
}

func (h *APIHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Selamat datang di Kasir API",
		"version": "0.1",
		"docs":    "/api",
	})
}
