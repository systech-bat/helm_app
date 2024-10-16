package tgclient

import (
	"strings"
	"weather-bot/pkg/config"
	"weather-bot/pkg/iqair"
	"weather-bot/pkg/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type TgClient struct {
	apiURL      string
	bot         *tgbotapi.BotAPI
	iqairClient *iqair.IqairClient
	tmpl        *template.TemplateRaw
	cfg         *config.Config
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/quality_now"),
	),
)

func (tg *TgClient) Listen() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tg.bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		if update.Message != nil && !update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			msg.ReplyMarkup = numericKeyboard

			if _, err := tg.bot.Send(msg); err != nil {
				log.Error("unable to send message tg.bot.Send(msg): ", err.Error())
			}
		} else if update.Message.IsCommand() {
			log.Infof("Receive tg command: %s", update.Message.Command())
			var pm10, pm25 string

			resp, err := tg.iqairClient.GetAirQuality(tg.cfg.Country, tg.cfg.State, tg.cfg.City)
			if err != nil {
				log.Error("unable to GetAirQuality(): ", err.Error())
			}

			//log.Debug(resp.Data.City.URL)

			if update.Message.Command() == "quality_now" {
				pm10, pm25, err = getToday(resp)
				if err != nil {
					log.Error("unable to GetAirQuality(): ", err.Error())
				}
			}

			log.Debugf("pm10: %s, pm25: %s", pm10, pm25)

			t := &template.Template{
				PM25: pm25,
				PM10: pm10,
				URL:  "https://www.iqair.com/" + strings.ToLower(tg.cfg.Country) + "/" + strings.ToLower(tg.cfg.State) + "/" + strings.ToLower(tg.cfg.City),
			}
			tg.tmpl.Tmpl = t
			payload, err := tg.tmpl.Parse()
			if err != nil {
				log.Error("unable to convert payload: ", err.Error())
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, payload)

			msg.ReplyMarkup = numericKeyboard
			msg.ParseMode = tgbotapi.ModeHTML

			if _, err := tg.bot.Send(msg); err != nil {
				log.Error("unable to send message tg.bot.Send(msg) 2: ", err.Error())
			}
		}
	}
}

func NewTgClient(tgUrl, tgToken string, iqairClient *iqair.IqairClient, tmpl *template.TemplateRaw, cfg *config.Config) (*TgClient, error) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		return nil, err
	}
	return &TgClient{
		apiURL:      tgUrl,
		bot:         bot,
		iqairClient: iqairClient,
		tmpl:        tmpl,
		cfg:         cfg,
	}, nil
}
