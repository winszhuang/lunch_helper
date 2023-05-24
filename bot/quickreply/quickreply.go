package quickreply

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func QuickReplyLocation() *linebot.QuickReplyItems {
	return linebot.NewQuickReplyItems(linebot.NewQuickReplyButton("", linebot.NewLocationAction("發送位置")))
}
