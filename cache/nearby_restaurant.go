package cache

import (
	"fmt"
	db "lunch_helper/db/sqlc"
	"lunch_helper/util"
	"sync"
)

// #TODO refactor to use redis
type NearByRestaurantCache struct {
	sync.Map
	mu sync.Mutex
}

type LocationContext struct {
	list      []db.Restaurant
	nextToken string
}

type LocationArgs struct {
	Lat    float64
	Lng    float64
	Radius int
}

const (
	LAST_TOKEN  = "end"
	FIRST_TOKEN = "start"
)

// e.g.
// key := generateKey(args{24.1677759, 120.6654513, 500})
// log.Println(key) // output: wsmc3z9gk_500
func generateKey(args LocationArgs) string {
	geoHash := util.ToGeoHash(args.Lat, args.Lng)
	return fmt.Sprintf("%s_%d", geoHash, args.Radius)
}

func NewNearByRestaurantCache() *NearByRestaurantCache {
	return &NearByRestaurantCache{}
}

func newLocationContext() *LocationContext {
	return &LocationContext{
		list:      []db.Restaurant{},
		nextToken: FIRST_TOKEN,
	}
}

func (nb *NearByRestaurantCache) checkLocationContext(args LocationArgs) *LocationContext {
	key := generateKey(args)
	value, _ := nb.LoadOrStore(key, newLocationContext())
	return value.(*LocationContext)
}

func (nb *NearByRestaurantCache) RemoveLocationContext(args LocationArgs) {
	key := generateKey(args)
	nb.Delete(key)
}

func (nb *NearByRestaurantCache) Append(args LocationArgs, list []db.Restaurant, pageToken string) {
	lc := nb.checkLocationContext(args)
	nb.mu.Lock()
	if pageToken == lc.nextToken {
		nb.mu.Unlock()
		return
	}
	if pageToken == "" {
		lc.nextToken = LAST_TOKEN
		nb.mu.Unlock()
		return
	}
	lc.list = append(lc.list, list...)
	nb.mu.Unlock()
}

func (nb *NearByRestaurantCache) GetRestaurantListByPagination(args LocationArgs, pageIndex, pageSize int) ([]db.Restaurant, bool) {
	lc := nb.checkLocationContext(args)
	result := util.Paginate(lc.list, pageIndex, pageSize)
	if nb.isRestaurantListFull(args) {
		return result, true
	}
	if len(result) == pageSize {
		return result, true
	}
	return result, false
}

// 該定位的附近店家資訊是否已滿
func (nb *NearByRestaurantCache) isRestaurantListFull(args LocationArgs) bool {
	lc := nb.checkLocationContext(args)
	return lc.nextToken == LAST_TOKEN
}
