package main

import (
	"fmt"
	"encoding/json"
	"log"
	"time"
	"errors"

	"gopkg.in/resty.v0"
	"github.com/tucnak/telebot"
)


func main() {
	bot, err := telebot.NewBot("TOKEN")
    if err != nil {
        log.Fatalln(err)
    }

    messages := make(chan telebot.Message, 100)
    bot.Listen(messages, 1*time.Second)

    for message := range messages {
        if message.Text == "/hi" {
		    weather, err := getWeather()
			if err != nil {
				bot.SendMessage(message.Chat, "Sorry, the problem is on our side", nil)
			} else {
				answer := fmt.Sprintf("\n\xE2\x99\xA8 Температура воздуха: *%s* °C", 	weather["outTemp"].(string)) +
						  fmt.Sprintf("\n\xF0\x9F\x92\xAA Давление: *%s* мм рт.ст",  	weather["pressure"].(string)) +
						  fmt.Sprintf("\n\xF0\x9F\x92\xA6 Влажность: *%s* %%", 		 	weather["outHumidity"].(string)) +
						  fmt.Sprintf("\n\xF0\x9F\x92\xA8 Ветер: *%s* м/с", 		 	weather["windSpeed"].(string)) +
						  fmt.Sprintf("\n\xE2\x98\x94 Дождь: *%s* мм/ч", 		 	 	weather["deltarain"].(string)) + 
						  fmt.Sprintf("\n\xF0\x9F\x92\xA1 Освещенность: *%s* люкс",  	weather["illumination"].(string))
				bot.SendMessage(message.Chat, answer, nil)
			}
        }
    }
}


func getWeather() (map[string]interface{}, error) {
	URL := "http://weather.thirdpin.ru/api/get"
	code, body, err := request(URL)

	var weather map[string]interface{}

	if err != nil {
		return weather, err
	}

	if code != 200 {
		return weather, errors.New(string(code))
	}

	if json.Unmarshal(body, &weather) != nil {
		panic(err)
	}

	return weather, nil
}

func request(URL string) (int, []byte, error) {

	var resp_code int
	var resp_body []byte

	resp, err := resty.R().Get(URL)
	if err != nil {
		return resp_code, resp_body, err
	}

	resp_code = resp.StatusCode()
	resp_body = resp.Body()

	return resp_code, resp_body, err
}

func parse(body []byte) map[string]interface{} {
	var weather map[string]interface{}
	if err := json.Unmarshal(body, &weather); err != nil {
        panic(err)
    }
    return weather
}

type ParseMode string
const (
        ModeDefault  ParseMode = ModeMarkdown
        ModeMarkdown ParseMode = "Markdown"
        ModeHTML     ParseMode = "HTML"
)
