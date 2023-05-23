package bot

import (
	"errors"
	"log"
	"lunch_helper/util"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type BotClient interface {
	ParseRequest(r *http.Request) ([]*linebot.Event, error)
	SendText(replyToken, text string)
	SendFlex(replyToken string, altText string, flexContainer linebot.FlexContainer)
	GetProfile(userID string) *linebot.GetProfileCall
}

type LineBotClient struct {
	*linebot.Client
}

func NewBotClient(channelSecret, channelToken string) (*LineBotClient, error) {
	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return nil, err
	}

	return &LineBotClient{bot}, nil
}

func (bc *LineBotClient) SetWebHookUrl(apiBaseUrl, endPoint string) error {
	url, err := util.BindUrl(apiBaseUrl, endPoint)
	if err != nil {
		return errors.New("bind url error: " + err.Error())
	}
	_, err = bc.Client.SetWebhookEndpointURL(url).Do()
	if err != nil {
		return errors.New("SetWebHookUrl Error: " + err.Error())
	}
	return nil
}

func (bc *LineBotClient) SendText(replyToken, text string) {
	_, err := bc.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do()
	if err != nil {
		log.Printf("SendText Error: %s", err)
	}
}

func (bc *LineBotClient) SendFlex(replyToken string, altText string, flexContainer linebot.FlexContainer) {
	_, err := bc.ReplyMessage(
		replyToken,
		linebot.NewFlexMessage(altText, flexContainer),
	).Do()
	if err != nil {
		log.Printf("SendText Error: %s", err)
	}
}
