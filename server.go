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
var zero int = 0
var one int = 1
var five int = 5

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
		err = handleTextMessage(event.ReplyToken, event.Message.(*linebot.TextMessage), event.Source)
	}
	return err
}

func handleTextMessage(replyToken string, message *linebot.TextMessage, source *linebot.EventSource) error {
	if source.Type != linebot.EventSourceTypeUser {
		// Room やグループでのメッセージはスルーする
		return nil
	}

	var err error
	switch message.Text {
	case "flex":
		err = replyFlexSample(replyToken)
	default:
		err = replyText(replyToken, message.Text)
	}
	return err
}

func replyText(token, text string) error {
	_, err := bot.ReplyMessage(token, linebot.NewTextMessage(text)).Do()
	return err
}

func replyFlexSample(token string) error {
	container := &linebot.BubbleContainer{
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			URL:         "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_1_cafe.png",
			Size:        linebot.FlexImageSizeTypeFull,
			AspectRatio: linebot.FlexImageAspectRatioType20to13,
			Action: &linebot.URIAction{
				URI: "http://linecorp.com/",
			},
		},
		Body: &linebot.BoxComponent{
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Text:   "Brown Cafe",
					Size:   linebot.FlexTextSizeTypeXl,
					Weight: linebot.FlexTextWeightTypeBold,
				},
				&linebot.BoxComponent{
					Layout: linebot.FlexBoxLayoutTypeBaseline,
					Margin: linebot.FlexComponentMarginTypeMd,
					Contents: []linebot.FlexComponent{
						&linebot.IconComponent{
							URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
							Size: linebot.FlexIconSizeTypeSm,
						},
						&linebot.IconComponent{
							URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
							Size: linebot.FlexIconSizeTypeSm,
						},
						&linebot.IconComponent{
							URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
							Size: linebot.FlexIconSizeTypeSm,
						},
						&linebot.IconComponent{
							URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
							Size: linebot.FlexIconSizeTypeSm,
						},
						&linebot.IconComponent{
							URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gray_star_28.png",
							Size: linebot.FlexIconSizeTypeSm,
						},
						&linebot.TextComponent{
							Text:  "4.0",
							Flex:  &zero,
							Size:  linebot.FlexTextSizeTypeSm,
							Color: "#999999",
						},
					},
				},
				&linebot.BoxComponent{
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Margin:  linebot.FlexComponentMarginTypeLg,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Text:  "Place",
									Flex:  &one,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#aaaaaa",
								},
								&linebot.TextComponent{
									Text:  "Miraina Tower, 4-1-6 Shinjuku, Tokyo",
									Flex:  &five,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#666666",
									Wrap:  true,
								},
							},
						},
						&linebot.BoxComponent{
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Text:  "Time",
									Flex:  &one,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#aaaaaa",
								},
								&linebot.TextComponent{
									Text:  "10:00 - 23:00",
									Flex:  &five,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#666666",
									Wrap:  true,
								},
							},
						},
					},
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Flex:    &zero,
			Spacing: linebot.FlexComponentSpacingTypeSm,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Height: linebot.FlexButtonHeightTypeSm,
					Style:  linebot.FlexButtonStyleTypeLink,
					Action: &linebot.URIAction{
						Label: "CALL",
						URI:   "https://linecorp.com",
					},
				},
				&linebot.ButtonComponent{
					Height: linebot.FlexButtonHeightTypeSm,
					Style:  linebot.FlexButtonStyleTypeLink,
					Action: &linebot.URIAction{
						Label: "WEBSITE",
						URI:   "https://linecorp.com",
					},
				},
			},
		},
	}
	message := linebot.NewFlexMessage("flex", container)
	_, err := bot.ReplyMessage(token, message).Do()
	return err
}
