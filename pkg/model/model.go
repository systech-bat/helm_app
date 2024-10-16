package model

import "time"

type IqairResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	City     string   `json:"city"`
	State    string   `json:"state"`
	Country  string   `json:"country"`
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Current struct {
	Pollution Pollution `json:"pollution"`
	Weather   Weather   `json:"weather"`
}

type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Pollution struct {
	Ts     time.Time `json:"ts"`
	Aqius  int       `json:"aqius"`
	Mainus string    `json:"mainus"`
	Aqicn  int       `json:"aqicn"`
	Maincn string    `json:"maincn"`
}

type Weather struct {
	Ts time.Time `json:"ts"`
	Tp int       `json:"tp"`
	Pr int       `json:"pr"`
	Hu int       `json:"hu"`
	Ws float64   `json:"ws"`
	Wd int       `json:"wd"`
	Ic string    `json:"ic"`
}

