package carousel

import (
	"fmt"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func CreateCarouselWithNext[T any](
	list []T,
	createBubbleFunc func(T) *linebot.BubbleContainer,
	nextPageIndex int,
	lat, lng float64,
) *linebot.CarouselContainer {
	bubble := CreateCarousel(list, createBubbleFunc)
	if nextPageIndex > 1 {
		bubble.Contents = append(bubble.Contents, createNext(nextPageIndex, lat, lng))
	}
	return bubble
}

func CreateCarousel[T any](list []T, createBubbleFunc func(T) *linebot.BubbleContainer) *linebot.CarouselContainer {
	containerList := []*linebot.BubbleContainer{}

	for _, item := range list {
		containerList = append(containerList, createBubbleFunc(item))
	}

	return &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: containerList,
	}
}

func createNext(nextPageIndex int, lat, lng float64) *linebot.BubbleContainer {
	nextData := fmt.Sprintf(
		"lat=%s,lng=%s,pageIndex=%d",
		strconv.FormatFloat(lat, 'f', 6, 64),
		strconv.FormatFloat(lng, 'f', 6, 64),
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
