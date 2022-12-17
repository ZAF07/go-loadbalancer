package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/ZAF07/go-loadbalancer/internal/config"
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

	maxLen := len(h.Config.Backends)
	w.Header().Add("node", h.Config.Backends[h.Id%maxLen].URL)
	h.Mu.Lock()
	currentBackend := h.Config.Backends[h.Id%maxLen].URL
	// targetURL, err := url.Parse(h.Config.Backends[h.Id%maxLen].URL)
	targetURL, err := url.Parse(currentBackend)
	if err != nil {
		log.Fatalf("error getting url : %+v", err)
	}

	// h.Id++
	h.setId()
	h.Mu.Unlock()
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.ServeHTTP(w, r)
}

func (h *Handler) setId() {
	if h.Id == 2 {
		h.Id = 0
		return
	}
	h.Id++
}
