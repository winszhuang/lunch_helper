package carousel

import (
	"fmt"
	"lunch_helper/constant"
	db "lunch_helper/db/sqlc"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func CreateRestaurantNextPageContainer(nextPageIndex int, lat, lng float64, radius int) *linebot.BubbleContainer {
	nextData := fmt.Sprintf(
		"lat=%s,lng=%s,radius=%d,pageIndex=%d",
		strconv.FormatFloat(lat, 'f', 6, 64),
		strconv.FormatFloat(lng, 'f', 6, 64),
		radius,
		nextPageIndex,
	)

	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeMicro,
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: "xs",
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Height: "sm",
					Action: &linebot.PostbackAction{
						Label: "下一頁資料",
						Data:  nextData,
					},
					Margin: linebot.FlexComponentMarginTypeLg,
				},
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Height: "sm",
					Style:  linebot.FlexButtonStyleTypeLink,
					Action: &linebot.URIAction{
						Label: "地圖上查看",
						URI:   "https://mileslin.github.io/2020/08/Golang/Live-Reload-For-Go/",
					},
				},
			},
		},
	}
}

// #TODO 要寫測試
func CreateRestaurantContainer(r db.Restaurant) *linebot.BubbleContainer {
	var image string
	if r.Image.Valid {
		image = r.Image.String
	} else {
		image = constant.DEFAULT_IMAGE_URL
	}

	if r.GoogleMapUrl == "" {
		fmt.Sprintf("%s找不到google map url", r.Name)
		r.GoogleMapUrl = "https://www.google.com.tw/maps"
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
						Label: "查看菜單",
						Data:  fmt.Sprintf("/restaurantmenu=%d", r.ID),
					},
					Margin: linebot.FlexComponentMarginTypeLg,
				},
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Height: "sm",
					Style:  linebot.FlexButtonStyleTypeLink,
					Action: &linebot.URIAction{
						Label: "地圖上查看",
						URI:   r.GoogleMapUrl,
					},
				},
			},
		},
	}
}
