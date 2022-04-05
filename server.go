package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

var secret string = ""
var token string = ""
var bot *linebot.Client

func init() {
	secret = os.Getenv("LINEBOT_CHANNEL_SECRET")
	token = os.Getenv("LINEBOT_CHANNEL_TOKEN")
	if secret == "" || token == "" {
		log.Fatal("Environment variable \"LINEBOT_CHANNEL_SECRET\" and \"LINEBOT_CHANNEL_TOKEN\" needs to be set.")
	}
	var err error
	bot, err = linebot.New(secret, token)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var addr = flag.String("a", "127.0.0.1:8000", "bind address")
	e := echo.New()
	e.Use(parser)
	e.POST("/", func(c echo.Context) error {
		events, ok := c.Get("events").([]*linebot.Event)
		if !ok {
			return c.NoContent(http.StatusBadRequest)
		}
		for _, event := range events {
			err := dispachEvent(event)
			if err != nil {
				c.Logger().Info(err)
				return c.NoContent(http.StatusBadRequest)
			}
		}
		return c.NoContent(http.StatusOK)
	})
	e.Logger.Fatal(e.Start(*addr))
}

func parser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		events, err := bot.ParseRequest(c.Request())
		if err != nil {
			c.Logger().Info(err)
			return c.NoContent(http.StatusBadRequest)
		}
		c.Set("events", events)
		return next(c)
	}
}

func dispachEvent(event *linebot.Event) error {
	var err error
	switch event.Type {
	case linebot.EventTypeMessage:
		err = dispatchMessage(event)
	}
	return err
}

func dispatchMessage(event *linebot.Event) error {
	var err error
	switch event.Message.(type) {
	case *linebot.TextMessage:
		err = handleTextMessage(event.Message.(*linebot.TextMessage), event.ReplyToken, event.Source)
	}
	return err
}

func handleTextMessage(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
		return err
	}
	return nil
}
