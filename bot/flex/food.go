package flex

import (
	"database/sql"
	"fmt"
	"lunch_helper/constant"
	db "lunch_helper/db/sqlc"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func CreateFoodItem(food db.Food) linebot.BubbleContainer {
	bubble := linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			URL:         checkImage(food.Image),
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
					Text:   food.Name,
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
									Text:   food.Price,
									Weight: linebot.FlexTextWeightTypeBold,
									Margin: linebot.FlexComponentMarginTypeSm,
									Flex:   linebot.IntPtr(0),
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "?kcl",
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
					Text:  checkDescription(food.Description),
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
					Style:  linebot.FlexButtonStyleTypePrimary,
					Color:  "#905c44",
					Margin: linebot.FlexComponentMarginTypeXxl,
					Action: &linebot.PostbackAction{
						Label: "加入收藏",
						Data:  fmt.Sprintf("/userlikefood=%d", food.ID),
					},
				},
			},
		},
	}

	return bubble
}

func checkImage(image sql.NullString) string {
	if image.Valid {
		return image.String
	}
	return constant.DEFAULT_IMAGE_URL
}

func checkDescription(description sql.NullString) string {
	if description.Valid {
		return description.String
	}
	return "No description"
}
