package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weather-bot/pkg/config"
	"weather-bot/pkg/iqair"

	log "github.com/sirupsen/logrus"
)

type AirHandler struct {
	iqairClient *iqair.IqairClient
	cfg         *config.Config
}

func (ah AirHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("Receive http request on: %s", r.URL.Path)
	resp, err := ah.iqairClient.GetAirQuality(ah.cfg.Country, ah.cfg.State, ah.cfg.City)
	if err != nil {
		log.Errorf("unable to GetAirQuality(): %s", err.Error())
	}
	b, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("unable to json.Marshal() response: %s", err.Error())
		return
	}
	fmt.Fprintf(w, "%s", b)
}

