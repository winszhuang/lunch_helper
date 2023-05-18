package bot

import (
	"errors"
	"log"
	"lunch_helper/util"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type BotClient struct {
	*linebot.Client
}

func NewBotClient(channelSecret, channelToken string) (*BotClient, error) {
	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return nil, err
	}

	return &BotClient{bot}, nil
}

func (bc *BotClient) SetWebHookUrl(apiBaseUrl, endPoint string) error {
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

func (bc *BotClient) SendText(replyToken, text string) {
	_, err := bc.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do()
	if err != nil {
		log.Printf("SendText Error: %s", err)
	}
}

func (bc *BotClient) SendFlex(replyToken string, altText string, flexContainer linebot.FlexContainer) {
	_, err := bc.ReplyMessage(
		replyToken,
		linebot.NewFlexMessage(altText, flexContainer),
	).Do()
	if err != nil {
		log.Printf("SendText Error: %s", err)
	}
}
