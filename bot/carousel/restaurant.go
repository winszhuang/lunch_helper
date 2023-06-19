package carousel

import (
	"fmt"
	"lunch_helper/constant"
	db "lunch_helper/db/sqlc"
	"lunch_helper/util"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/shopspring/decimal"
)

func CreateRestaurantCarouselItem(r db.Restaurant, postbackContents func(r db.Restaurant) []linebot.FlexComponent) *linebot.BubbleContainer {
	var image string
	if r.Image.Valid {
		image = r.Image.String
	} else {
		image = constant.DEFAULT_IMAGE_URL
	}

	var userRatingTotal string
	if r.UserRatingsTotal.Valid {
		userRatingTotal = strconv.Itoa(int(r.UserRatingsTotal.Int32))
	} else {
		userRatingTotal = "no data"
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
					Text:   util.NoEmptyString(r.Name),
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
							Contents:     CreateRatingComponent(r.Rating),
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
									Text:  util.NoEmptyString(r.Address),
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
									Text:  "Phone",
									Color: "#aaaaaa",
									Size:  "sm",
									Flex:  linebot.IntPtr(1),
								},
								&linebot.TextComponent{
									Type:  "text",
									Text:  util.NoEmptyString(r.PhoneNumber),
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
									Text:  "評論數",
									Color: "#aaaaaa",
									Size:  "sm",
									Flex:  linebot.IntPtr(1),
								},
								&linebot.TextComponent{
									Type:  "text",
									Text:  userRatingTotal,
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
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Spacing:  "xs",
			Contents: postbackContents(r),
		},
	}
}

func PostBackContentsWithShowMenuAndLikeAndViewOnMap(r db.Restaurant) []linebot.FlexComponent {
	return []linebot.FlexComponent{
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
			Action: &linebot.PostbackAction{
				Label: "加入收藏",
				Data:  fmt.Sprintf("/userlikerestaurant=%d", r.ID),
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
	}
}

func PostBackContentsWithShowMenuAndUnLikeAndViewOnMap(r db.Restaurant) []linebot.FlexComponent {
	return []linebot.FlexComponent{
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
			Action: &linebot.PostbackAction{
				Label: "取消收藏",
				Data:  fmt.Sprintf("/userunlikerestaurant=%d", r.ID),
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
	}
}

func CreateRatingComponent(rating decimal.Decimal) []linebot.FlexComponent {
	roundedRating := rating.Floor()

	goldStarURL := "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
	grayStarURL := "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gray_star_28.png"

	components := []linebot.FlexComponent{}
	for i := 0; i < 5; i++ {
		var starURL string
		if i < int(roundedRating.IntPart()) {
			starURL = goldStarURL
		} else {
			starURL = grayStarURL
		}

		components = append(components, &linebot.IconComponent{
			Type: linebot.FlexComponentTypeIcon,
			URL:  starURL,
			Size: linebot.FlexIconSizeTypeSm,
		})
	}

	return components
}
