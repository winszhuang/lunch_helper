package richmenu

import (
	"lunch_helper/bot"
	"lunch_helper/constant"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type InputOptionType string

const (
	CloseRichMenu InputOptionType = "closeRichMenu"
	OpenRichMenu  InputOptionType = "openRichMenu"
	OpenKeyBoard  InputOptionType = "openKeyboard"
	OpenVoice     InputOptionType = "openVoice"
)

type RichMenuActionWithInputOption struct {
	Type            linebot.RichMenuActionType `json:"type"`
	URI             string                     `json:"uri,omitempty"`
	Text            string                     `json:"text,omitempty"`
	DisplayText     string                     `json:"displayText,omitempty"`
	Label           string                     `json:"label,omitempty"`
	Data            string                     `json:"data,omitempty"`
	Mode            string                     `json:"mode,omitempty"`
	Initial         string                     `json:"initial,omitempty"`
	Max             string                     `json:"max,omitempty"`
	Min             string                     `json:"min,omitempty"`
	RichMenuAliasID string                     `json:"richMenuAliasId,omitempty"`
	InputOption     string                     `json:"inputOption"`
}

func CreateSearchRichMenu() bot.RichMenu_New {
	return bot.RichMenu_New{
		Size: linebot.RichMenuSize{
			Width:  2500,
			Height: 1686,
		},
		Name:        "Search",
		ChatBarText: "搜尋功能",
		Areas: []bot.AreaDetail_New{
			// 左上
			{
				Bounds: linebot.RichMenuBounds{
					X:      0,
					Y:      0,
					Width:  833,
					Height: 843,
				},
				Action: bot.RichMenuAction_New{
					Type: "postback",
					Data: string(constant.Search),
				},
			},
			// 中間上
			{
				Bounds: linebot.RichMenuBounds{
					X:      833,
					Y:      0,
					Width:  833,
					Height: 843,
				},
				Action: bot.RichMenuAction_New{
					Type: "postback",
					Data: string(constant.SearchLocation),
				},
			},
			// 右上
			{
				Bounds: linebot.RichMenuBounds{
					X:      1666,
					Y:      0,
					Width:  833,
					Height: 843,
				},
				Action: bot.RichMenuAction_New{
					Type:        "postback",
					Data:        string(constant.SearchText),
					InputOption: string(OpenKeyBoard),
				},
			},
			// 左下
			{
				Bounds: linebot.RichMenuBounds{
					X:      0,
					Y:      843,
					Width:  833,
					Height: 843,
				},
				Action: bot.RichMenuAction_New{
					Type: "postback",
					Data: string(constant.SearchRadius),
				},
			},
			// 中間下
			{
				Bounds: linebot.RichMenuBounds{
					X:      833,
					Y:      843,
					Width:  833,
					Height: 843,
				},
				Action: bot.RichMenuAction_New{
					Type: "postback",
					Data: string(constant.SearchAI),
				},
			},
			// 右下
			{
				Bounds: linebot.RichMenuBounds{
					X:      1666,
					Y:      843,
					Width:  833,
					Height: 843,
				},
				Action: bot.RichMenuAction_New{
					Type: "postback",
					Data: string(constant.UserOption),
				},
			},
		},
	}
}
