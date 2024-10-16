package tgclient

import (
	"strconv"
	"weather-bot/pkg/model"
)

func getToday(resp *model.IqairResponse) (string, string, error) {

	return "-", strconv.Itoa(resp.Data.Current.Pollution.Aqius), nil
}

