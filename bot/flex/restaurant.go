package flex

import (
	"database/sql"
	"fmt"
	"lunch_helper/constant"
	db "lunch_helper/db/sqlc"
	"lunch_helper/util"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func CreateRestaurantItem(restaurant db.Restaurant, isUserRestaurant bool) linebot.BubbleContainer {
	var labelText string
	var postData string
	if isUserRestaurant {
		labelText = "取消收藏"
		postData = fmt.Sprintf("/userunlikerestaurant=%d", restaurant.ID)
	} else {
		labelText = "加入收藏"
		postData = fmt.Sprintf("/userlikerestaurant=%d", restaurant.ID)
	}
	bubble := linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			URL:         checkRestaurantImage(restaurant.Image),
			Size:        linebot.FlexImageSizeTypeFull,
			AspectRatio: linebot.FlexImageAspectRatioType20to13,
			AspectMode:  linebot.FlexImageAspectModeTypeCover,
		},
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   restaurant.Name,
					Size:   linebot.FlexTextSizeTypeXl,
					Weight: linebot.FlexTextWeightTypeBold,
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeBaseline,
							Contents: []linebot.FlexComponent{
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/restaurant_regular_32.png",
								},
								&linebot.TextComponent{
									Type:   linebot.FlexComponentTypeText,
									Text:   restaurant.PhoneNumber,
									Weight: linebot.FlexTextWeightTypeBold,
									Margin: linebot.FlexComponentMarginTypeSm,
									Flex:   linebot.IntPtr(0),
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "phone",
									Size:  linebot.FlexTextSizeTypeSm,
									Align: linebot.FlexComponentAlignTypeEnd,
									Color: "#aaaaaa",
								},
							},
						},
					},
				},
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  util.NoEmptyString(restaurant.Address),
					Wrap:  true,
					Color: "#aaaaaa",
					Size:  linebot.FlexTextSizeTypeXxs,
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypeLink,
					Margin: linebot.FlexComponentMarginTypeSm,
					Action: &linebot.PostbackAction{
						Label: "查看菜單",
						Data:  fmt.Sprintf("/restaurantmenu=%d", restaurant.ID),
					},
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypeLink,
					Action: &linebot.URIAction{
						Label: "地圖上查看",
						URI:   restaurant.GoogleMapUrl,
					},
				},
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypePrimary,
					Color:  "#905c44",
					Margin: linebot.FlexComponentMarginTypeXxl,
					Action: &linebot.PostbackAction{
						Label: labelText,
						Data:  postData,
					},
				},
			},
		},
	}

	return bubble
}

func checkRestaurantImage(image sql.NullString) string {
	if image.Valid {
		return image.String
	}
	return constant.DEFAULT_IMAGE_URL
}
