package richmenu

import (
	"lunch_helper/constant"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func CreateSearchRichMenu() linebot.RichMenu {
	return linebot.RichMenu{
		Size: linebot.RichMenuSize{
			Width:  2500,
			Height: 1686,
		},
		Name:        "Search",
		ChatBarText: "搜尋功能",
		Areas: []linebot.AreaDetail{
			// 左上
			{
				Bounds: linebot.RichMenuBounds{
					X:      0,
					Y:      0,
					Width:  833,
					Height: 843,
				},
				Action: linebot.RichMenuAction{
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
				Action: linebot.RichMenuAction{
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
				Action: linebot.RichMenuAction{
					Type: "postback",
					Data: string(constant.SearchText),
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
				Action: linebot.RichMenuAction{
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
				Action: linebot.RichMenuAction{
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
				Action: linebot.RichMenuAction{
					Type: "postback",
					Data: string(constant.UserOption),
				},
			},
		},
	}
}
