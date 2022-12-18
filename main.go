package main

import (
	"log"
	"net/http"

	"github.com/ZAF07/go-loadbalancer/internal/config"
	"github.com/ZAF07/go-loadbalancer/internal/handler"
)

func main() {
	cfg := config.LoadConfigs()
	Serve(cfg)
}

func Serve(cfg *config.Config) {
	handler := handler.NewHandler(cfg)

	s := http.Server{
		Addr:    cfg.Proxy.Port,
		Handler: http.HandlerFunc(handler.LoadBalancerHandler),
	}
	log.Println("SERVER STARTED")
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
