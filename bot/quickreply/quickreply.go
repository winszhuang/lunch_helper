package quickreply

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func QuickReplyLocation() *linebot.QuickReplyItems {
	return linebot.NewQuickReplyItems(linebot.NewQuickReplyButton("", linebot.NewLocationAction("發送位置")))
}

func QuickReplyRadiusOptions() *linebot.QuickReplyItems {
	return linebot.NewQuickReplyItems(
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("100", "100")),
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("500", "500")),
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("1000", "1000")),
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("2000", "2000")),
	)
}
