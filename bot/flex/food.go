package flex

import (
	db "lunch_helper/db/sqlc"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func CreateFoodListContainer(foods []db.Food, r db.Restaurant) linebot.BubbleContainer {
	bubble := linebot.BubbleContainer{
		Type: "bubble",
		Body: &linebot.BoxComponent{
			Type:   "box",
			Layout: "vertical",
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   "text",
					Text:   "菜單",
					Weight: "bold",
					Color:  "#1DB446",
					Size:   "sm",
				},
				&linebot.TextComponent{
					Type:   "text",
					Text:   r.Name,
					Weight: "bold",
					Size:   "xxl",
					Margin: "md",
				},
				&linebot.TextComponent{
					Type:  "text",
					Text:  r.Address,
					Size:  "xs",
					Color: "#aaaaaa",
					Wrap:  true,
				},
				&linebot.SeparatorComponent{
					Type:   "separator",
					Margin: "xxl",
				},
				&linebot.BoxComponent{
					Type:     "box",
					Layout:   "vertical",
					Margin:   "xxl",
					Spacing:  "sm",
					Contents: createFoodsContent(foods),
				},
				&linebot.SeparatorComponent{
					Type:   "separator",
					Margin: "xxl",
				},
				&linebot.BoxComponent{
					Type:   "box",
					Layout: "horizontal",
					Margin: "md",
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  "text",
							Text:  "電話",
							Size:  "xs",
							Color: "#aaaaaa",
							Flex:  linebot.IntPtr(0),
						},
						&linebot.TextComponent{
							Type:  "text",
							Text:  r.PhoneNumber,
							Color: "#aaaaaa",
							Size:  "xs",
							Align: "end",
						},
					},
				},
			},
		},
		Styles: &linebot.BubbleStyle{
			Footer: &linebot.BlockStyle{
				Separator: true,
			},
		},
	}

	return bubble
}

func createFoodsContent(foods []db.Food) []linebot.FlexComponent {
	var contents []linebot.FlexComponent
	for _, food := range foods {
		contents = append(contents, createFoodContent(food))
	}
	return contents
}

func createFoodContent(food db.Food) linebot.FlexComponent {
	return &linebot.BoxComponent{
		Type:   "box",
		Layout: "horizontal",
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:  "text",
				Text:  food.Name,
				Size:  "sm",
				Color: "#555555",
				Flex:  linebot.IntPtr(0),
			},
			&linebot.TextComponent{
				Type:  "text",
				Text:  food.Price,
				Size:  "sm",
				Color: "#111111",
				Align: "end",
			},
		},
	}
}
