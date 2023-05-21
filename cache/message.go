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
	userId       string
	conversation UserConversation
	location     []UserLocation
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

func NewMessageCache() *MessageCache {
	return &MessageCache{}
}

func newUserMessageCache(userId string) *UserMessageContext {
	return &UserMessageContext{
		userId:       userId,
		conversation: UserConversation{Data: &Merchant{}},
		location:     []UserLocation{},
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
	singleUserCache.location = []UserLocation{}
	singleUserCache.conversation = UserConversation{Data: &Merchant{}}
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

func (mc *MessageCache) GetUserMode(userId string) constant.Directive {
	singleUserCache := mc.checkUserMessageCache(userId)
	return singleUserCache.conversation.Mode
}

func (mc *MessageCache) GetLocation(userId string) []UserLocation {
	singleUserCache := mc.checkUserMessageCache(userId)
	return singleUserCache.location
}
