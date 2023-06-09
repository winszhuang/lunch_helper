package richmenu

import (
	"lunch_helper/bot"
	"lunch_helper/constant"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func CreateUserRichMenu() bot.RichMenu_New {
	return bot.RichMenu_New{
		Size: linebot.RichMenuSize{
			Width:  2500,
			Height: 1686,
		},
		Name:        "richmenu-user",
		ChatBarText: "使用者選單",
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
					Data: string(constant.FavoriteRestaurants),
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
					Data: string(constant.FavoriteFoods),
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
					Type: "postback",
					Data: string(constant.PickRestaurant),
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
					Data: string(constant.NotificationSetting),
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
					// #TODO 增加趨勢功能
					Data: "/trend",
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
					Type:            "richmenuswitch",
					RichMenuAliasID: "richmenu-search",
					Data:            string(constant.SearchOption),
				},
			},
		},
	}
}
