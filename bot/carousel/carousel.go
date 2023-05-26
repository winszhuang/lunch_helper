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
