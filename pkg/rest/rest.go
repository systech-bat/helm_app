package rest

import (
	"fmt"
	"net/http"
	"weather-bot/pkg/config"
	"weather-bot/pkg/iqair"
)

type HttpClient struct {
	iqairClient *iqair.IqairClient
	cfg         *config.Config
}

func (h *HttpClient) Listen() {
	airHandler := &AirHandler{
		iqairClient: h.iqairClient,
		cfg:         h.cfg,
	}
	http.Handle("/air", airHandler)
	port := h.cfg.HttpPort
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func NewHttpClient(iqairClient *iqair.IqairClient, cfg *config.Config) *HttpClient {

	return &HttpClient{
		iqairClient: iqairClient,
		cfg:         cfg,
	}
}

