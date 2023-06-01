package service

// 這邊只負責處理grpc
type CrawlerService struct {
	// grpc實體

}

type CrawFood struct {
	Name        string
	Price       string
	Image       string
	Description string
	Category    string
}

func NewCrawlerService(restaurantService *RestaurantService) *CrawlerService {
	return &CrawlerService{}
}

// #TODO 抓取一個商家的餐點
func (s *CrawlerService) CrawlFoodsFromGoogleMap(googleMapUrl string) ([]CrawFood, error) {
	// grpc跟python爬蟲service要資料
	// ctx := context.Background()
	// 0. 確認資料庫該restaurant是否已經被爬取過(需判斷該restaurant是否在資料庫建檔)，有爬過就不爬取
	// 1. 呼叫grpc方法CrawlDishes取得商家的菜單
	// 2. 確認是否抓取到菜單，沒有抓到菜單就在
	// s.dbStore.
	return []CrawFood{}, nil
}
