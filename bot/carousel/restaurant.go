package carousel

import (
	db "lunch_helper/db/sqlc"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const DEFAULT_IMAGE_URL = "https://mnapoli.fr/images/posts/null.png"

// #TODO 要寫測試
func CreateRestaurantContainer(r db.Restaurant) *linebot.BubbleContainer {
	var image string
	if r.Image.Valid {
		image = r.Image.String
	} else {
		image = DEFAULT_IMAGE_URL
	}

	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeMicro,
		Hero: &linebot.ImageComponent{
			Type: linebot.FlexComponentTypeImage,
			// #NOTICE must be a valid image url
			URL:         image,
			Size:        "full",
			AspectMode:  "cover",
			AspectRatio: "320:213",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   r.Name,
					Weight: "bold",
					Size:   "sm",
					Wrap:   true,
				},
				&linebot.BoxComponent{
					Type:       linebot.FlexComponentTypeBox,
					Layout:     linebot.FlexBoxLayoutTypeHorizontal,
					PaddingTop: "5px",
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:         linebot.FlexComponentTypeBox,
							Layout:       linebot.FlexBoxLayoutTypeBaseline,
							PaddingStart: "2px",
							Contents: []linebot.FlexComponent{
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
							},
						},
						&linebot.TextComponent{
							Text:   r.Rating.String(),
							Type:   linebot.FlexComponentTypeText,
							Size:   "xs",
							Color:  "#8c8c8c",
							Margin: "md",
							Flex:   linebot.IntPtr(0),
						},
					},
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Margin:  "lg",
					Spacing: "sm",
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: "sm",
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  "text",
									Text:  "Place",
									Color: "#aaaaaa",
									Size:  "sm",
									Flex:  linebot.IntPtr(1),
								},
								&linebot.TextComponent{
									Type:  "text",
									Text:  r.Address,
									Color: "#666666",
									Wrap:  true,
									Size:  "sm",
									Flex:  linebot.IntPtr(5),
								},
							},
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: "sm",
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  "text",
									Text:  "Time",
									Color: "#aaaaaa",
									Size:  "sm",
									Flex:  linebot.IntPtr(1),
								},
								&linebot.TextComponent{
									Type:  "text",
									Text:  "no record",
									Color: "#666666",
									Wrap:  true,
									Size:  "sm",
									Flex:  linebot.IntPtr(5),
								},
							},
						},
					},
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: "xs",
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Height: "sm",
					Action: &linebot.PostbackAction{
						Label: "選擇餐廳",
						Data:  "&action=restaurant",
					},
					Margin: linebot.FlexComponentMarginTypeLg,
				},
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Height: "sm",
					Style:  linebot.FlexButtonStyleTypeLink,
					Action: &linebot.URIAction{
						Label: "詳細資料",
						URI:   "https://mileslin.github.io/2020/08/Golang/Live-Reload-For-Go/",
					},
				},
			},
		},
	}
}
