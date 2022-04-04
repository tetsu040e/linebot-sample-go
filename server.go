package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Message struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Text string `json:"text"`
}

type Source struct {
	Type   string `json:"type"`
	UserId string `json:"userId"`
}

type Event struct {
	Type       string  `json:"type"`
	Message    Message `json:"message"`
	Timestamp  uint64  `json:"timestamp"`
	Sourcee    Source  `json:"source"`
	ReplyToken string  `json:"replyToken"`
	Mode       string  `json:"mode"`
}

type EventsObject struct {
	Destination string  `json:"destination"`
	Events      []Event `json:"events"`
}

var secret string = ""
var token string = ""

func init() {
	secret = os.Getenv("LINEBOT_CHANNEL_SECRET")
	token = os.Getenv("LINEBOT_CHANNEL_TOKEN")
	if secret == "" || token == "" {
		log.Fatal("Environment variable \"LINEBOT_CHANNEL_SECRET\" and \"LINEBOT_CHANNEL_TOKEN\" needs to be set.")
	}
}

func main() {
	var addr = flag.String("a", "127.0.0.1:8000", "bind address")
	e := echo.New()
	e.Use(validator)
	e.Use(parser)
	e.POST("/", func(c echo.Context) error {
		events, ok := c.Get("events").(EventsObject)
		if !ok {
			return c.String(http.StatusBadRequest, "NG")
		}
		bot, err := linebot.New(secret, token)
		if err != nil {
			return c.String(http.StatusBadRequest, "NG")
		}
		for _, event := range events.Events {
			if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(event.Message.Text)).Do(); err != nil {
				return c.String(http.StatusBadRequest, "NG")
			}
		}
		log.Printf("%v", events)
		return c.String(http.StatusOK, "Hello!!")
	})
	e.Logger.Fatal(e.Start(*addr))
}

func validator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusBadRequest, "NG")
		}
		log.Print(string(body))
		signature := c.Request().Header.Get("x-line-signature")
		decoded, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			return c.String(http.StatusBadRequest, "NG")
		}
		hash := hmac.New(sha256.New, []byte(secret))
		hash.Write(body)

		if !hmac.Equal(decoded, hash.Sum(nil)) {
			return c.String(http.StatusForbidden, "NG")
		}

		c.Set("rawbody", body)
		return next(c)
	}
}

// XXX ToDo: SDK 使うように修正する
func parser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, ok := c.Get("rawbody").([]byte)
		if !ok {
			return c.String(http.StatusBadRequest, "NG")
		}

		var events EventsObject
		if err := json.Unmarshal(body, &events); err != nil {
			return c.String(http.StatusBadRequest, "NG")
		}
		c.Set("events", events)
		return next(c)
	}
}
