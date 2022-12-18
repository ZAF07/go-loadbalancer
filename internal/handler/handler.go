package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/ZAF07/go-loadbalancer/internal/config"
)

const (
	FAVICON = "/favicon.ico"
	SERVER  = "/"
)

type Handler struct {
	Id     int
	Mu     *sync.Mutex
	Config *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		Id:     0,
		Mu:     &sync.Mutex{},
		Config: cfg,
	}
}

func (h *Handler) LoadBalancerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case FAVICON:
		faviconHandler(w, r)
	case SERVER:
		h.serveProxy(w, r)
		return
	}
}

func (h *Handler) serveProxy(w http.ResponseWriter, r *http.Request) {
	maxLen := len(h.Config.Backends)
	h.Mu.Lock()
	serverIsDead := h.Config.Backends[h.Id%maxLen].GetStatus()
	if serverIsDead {
		h.setId()
	}

	targetURL, err := url.Parse(h.Config.Backends[h.Id%maxLen].URL)
	if err != nil {
		log.Fatalf("error getting url : %+v", err)
	}

	h.setId()
	h.Mu.Unlock()

	w.Header().Add("node", h.Config.Backends[h.Id%maxLen].URL)
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.ServeHTTP(w, r)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

func (h *Handler) setId() {
	if h.Id == 2 {
		h.Id = 0
		return
	}
	h.Id++
}
