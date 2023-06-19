package cache

import (
	"lunch_helper/constant"
	"sync"

	"googlemaps.github.io/maps"
)

type Merchant struct {
	Name  string
	Phone string
}

type MessageCache struct {
	sync.Map
}

type UserMessageContext struct {
	userId              string
	conversation        *UserConversation
	currentLocation     *UserLocation
	currentSearchRadius int
	currentSearchText   string
	searchTextList      []string
	locationList        []UserLocation
}

type UserConversation struct {
	Mode     constant.Directive
	Question constant.Question
	Data     *Merchant
}

type UserLocation struct {
	ID      string
	Title   string
	Address string
	LatLng  maps.LatLng
}

const (
	DEFAULT_SEARCH_RADIUS = 500
	DEFAULT_SEARCH_TEXT   = ""
)

func NewMessageCache() *MessageCache {
	return &MessageCache{}
}

func newUserMessageCache(userId string) *UserMessageContext {
	return &UserMessageContext{
		userId:              userId,
		conversation:        &UserConversation{Data: &Merchant{}},
		locationList:        []UserLocation{},
		currentLocation:     nil,
		currentSearchRadius: DEFAULT_SEARCH_RADIUS,
		currentSearchText:   DEFAULT_SEARCH_TEXT,
		searchTextList:      []string{},
	}
}

func (mc *MessageCache) checkUserMessageCache(userId string) *UserMessageContext {
	value, _ := mc.LoadOrStore(userId, newUserMessageCache(userId))
	return value.(*UserMessageContext)
}

func (mc *MessageCache) RemoveUserMessage(userId string) {
	mc.Delete(userId)
}

func (mc *MessageCache) ResetUserMessage(userId string) {
	singleUserCache := mc.checkUserMessageCache(userId)
	singleUserCache.locationList = []UserLocation{}
	singleUserCache.currentLocation = nil
	singleUserCache.conversation = &UserConversation{Data: &Merchant{}}
}

func (mc *MessageCache) SetMode(userId string, mode constant.Directive) {
	singleUserCache := mc.checkUserMessageCache(userId)
	singleUserCache.conversation.Mode = mode
}

func (mc *MessageCache) SetQuestion(userId string, question constant.Question) {
	singleUserCache := mc.checkUserMessageCache(userId)
	singleUserCache.conversation.Question = question
}

func (mc *MessageCache) SetData(userId string, fn func(*Merchant) *Merchant) {
	singleUserCache := mc.checkUserMessageCache(userId)
	singleUserCache.conversation.Data = fn(singleUserCache.conversation.Data)
}

func (mc *MessageCache) UpdateLocation(userId string, location *UserLocation) {
	singleUserCache := mc.checkUserMessageCache(userId)
	singleUserCache.currentLocation = location
	for _, l := range singleUserCache.locationList {
		if l.ID == location.ID {
			return
		}
	}
	singleUserCache.locationList = append(singleUserCache.locationList, *location)
}

func (mc *MessageCache) UpdateSearchText(userId string, text string) {
	singleUserCache := mc.checkUserMessageCache(userId)
	singleUserCache.currentSearchText = text
	for _, t := range singleUserCache.searchTextList {
		if t == text {
			return
		}
	}
	singleUserCache.searchTextList = append(singleUserCache.searchTextList, text)

}

func (mc *MessageCache) UpdateSearchRadius(userId string, radius int) {
	singleUserCache := mc.checkUserMessageCache(userId)
	singleUserCache.currentSearchRadius = radius
}

func (mc *MessageCache) GetUserMode(userId string) constant.Directive {
	singleUserCache := mc.checkUserMessageCache(userId)
	return singleUserCache.conversation.Mode
}

func (mc *MessageCache) GetCurrentLocation(userId string) (*UserLocation, bool) {
	singleUserCache := mc.checkUserMessageCache(userId)
	if singleUserCache.currentLocation == nil {
		return nil, false
	}

	return singleUserCache.currentLocation, true
}

func (mc *MessageCache) GetLocationList(userId string) []UserLocation {
	singleUserCache := mc.checkUserMessageCache(userId)
	return singleUserCache.locationList
}

func (mc *MessageCache) GetCurrentSearchText(userId string) string {
	singleUserCache := mc.checkUserMessageCache(userId)
	return singleUserCache.currentSearchText
}

func (mc *MessageCache) GetSearchTextList(userId string) []string {
	singleUserCache := mc.checkUserMessageCache(userId)
	return singleUserCache.searchTextList
}

func (mc *MessageCache) GetCurrentRadius(userId string) int {
	singleUserCache := mc.checkUserMessageCache(userId)
	return singleUserCache.currentSearchRadius
}
