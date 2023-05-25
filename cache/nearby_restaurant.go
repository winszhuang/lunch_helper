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
}

type LocationContext struct {
	mu    sync.RWMutex
	pages []PageDataOfPlaces
}

type PageDataOfPlaces struct {
	currentToken  string
	nextPageToken string
	data          []db.Restaurant
}

type LocationArgs struct {
	Lat    float64
	Lng    float64
	Radius int
}

func NewPageDataOfPlaces(currentToken, nextPageToken string, data []db.Restaurant) PageDataOfPlaces {
	return PageDataOfPlaces{currentToken, nextPageToken, data}
}

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

func (nb *NearByRestaurantCache) checkLocationContext(args LocationArgs) *LocationContext {
	key := generateKey(args)
	value, _ := nb.LoadOrStore(key, &LocationContext{pages: []PageDataOfPlaces{}})
	return value.(*LocationContext)
}

func (nb *NearByRestaurantCache) RemoveLocationContext(args LocationArgs) {
	key := generateKey(args)
	nb.Delete(key)
}

func (nb *NearByRestaurantCache) Append(args LocationArgs, pageDataOfPlaces PageDataOfPlaces) {
	lc := nb.checkLocationContext(args)
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if lc.isFull() || lc.isExist(pageDataOfPlaces) {
		return
	}

	if lc.canAdd(pageDataOfPlaces) {
		lc.addPage(pageDataOfPlaces)
	}
}

// if the second return value is false, it is mean amount of data is not enough, need call api to get more
func (nb *NearByRestaurantCache) GetRestaurantListByPagination(args LocationArgs, pageIndex, pageSize int) ([]db.Restaurant, bool) {
	lc := nb.checkLocationContext(args)
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	list := lc.listAll()
	result := util.Paginate(list, pageIndex, pageSize)
	isEnough := lc.isFull() || len(result) == pageSize

	return result, isEnough
}

func (lc *LocationContext) addPage(p PageDataOfPlaces) {
	lc.pages = append(lc.pages, p)
}

func (lc *LocationContext) listAll() []db.Restaurant {
	result := []db.Restaurant{}
	for _, page := range lc.pages {
		result = append(result, page.data...)
	}
	return result
}

func (lc *LocationContext) isExist(p PageDataOfPlaces) bool {
	exist := false
	for _, page := range lc.pages {
		if page.currentToken == p.currentToken && page.nextPageToken == p.nextPageToken {
			exist = true
			break
		}
	}
	return exist
}

func (lc *LocationContext) isFull() bool {
	for _, page := range lc.pages {
		if page.nextPageToken == "" {
			return true
		}
	}
	return false
}

func (lc *LocationContext) canAdd(p PageDataOfPlaces) bool {
	// 首筆資料一定能新增
	if p.currentToken == "" && len(lc.pages) == 0 {
		return true
	}

	len := len(lc.pages)
	return p.currentToken == lc.pages[len-1].nextPageToken
}
