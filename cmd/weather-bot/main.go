package main

import (
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"weather-bot/pkg/config"
	"weather-bot/pkg/iqair"
	"weather-bot/pkg/rest"
	"weather-bot/pkg/template"
	"weather-bot/pkg/tgclient"

	log "github.com/sirupsen/logrus"
)

const maxGoroutines = 100

func main() {

	// Init service
	cfg := config.AppConfig

	tmpl, err := template.LoadTemplate(cfg.TemplatePath)
	if err != nil {
		log.Warn("Unable to load message template: %s", err.Error())
	}

	iqairClient := iqair.NewIqairClient(cfg.IqairUrl, cfg.IqairToken)
	tgClient, err := tgclient.NewTgClient(cfg.TgApiUrl, cfg.TgApiToken, iqairClient, tmpl, cfg)
	if err != nil {
		log.Fatalf("Unable to initialize tg client: %s", err.Error())
	}
	httpClient := rest.NewHttpClient(iqairClient, cfg)

	log.Info("Start listening for telegram events...")

	// Run Clients
	go httpClient.Listen()
	go tgClient.Listen()

	// Run healthchecks
	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		if checkLiveness() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("{}"))
		}
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if checkReadiness(cfg.IqairUrl, cfg.IqairToken) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("{}"))
		}
	})

	http.ListenAndServe(":8086", nil)
	return
}

func checkReadiness(apiUrl, apiKey string) bool {
	rawUrl, err := url.Parse(apiUrl)
	if err != nil {
		log.Error("Error parsing URL:", err.Error())
		return false
	}

	v := url.Values{}
	v.Add("key", apiKey)
	rawUrl.Path = "/v2/countries"
	rawUrl.RawQuery = v.Encode()

	resp, err := http.Get(rawUrl.String())
	if err != nil {
		log.Warn("Readiness check failed: ", err.Error())
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if strings.Contains(resp.Status, "429") {
			log.Warnf("Readiness check returned status [%s], take a look on issue: https://gitlab.com/ksxack/weather-bot/-/issues/1", resp.Status)
			return false
		}
		log.Warn("Readiness check returned non-200 status: ", resp.Status)
		return false
	}

	return true
}

func checkLiveness() bool {
	numGoroutines := runtime.NumGoroutine()
	if numGoroutines > maxGoroutines {
		log.Warnf("Liveness check failed: too many goroutines (%d)", numGoroutines)
		return false
	}
	return true
}
