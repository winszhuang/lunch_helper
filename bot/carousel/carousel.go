package carousel

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func CreateCarouselWithNext[T any](
	list []T,
	createBubbleFunc func(T) *linebot.BubbleContainer,
	createNextFunc func() *linebot.BubbleContainer,
) *linebot.CarouselContainer {
	carouselBubble := CreateCarousel(list, createBubbleFunc)
	nextBubble := createNextFunc()
	if nextBubble != nil {
		carouselBubble.Contents = append(carouselBubble.Contents, nextBubble)
	}
	return carouselBubble
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

func CreateNextPageContainer(data string) *linebot.BubbleContainer {
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
						Data:  data,
					},
					Margin: linebot.FlexComponentMarginTypeLg,
				},
			},
		},
	}
}
