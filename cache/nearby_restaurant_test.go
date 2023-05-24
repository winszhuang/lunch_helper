package cache

import (
	"fmt"
	db "lunch_helper/db/sqlc"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func generateMockRestaurants(startIndex, endIndex int) []db.Restaurant {
	rand.Seed(time.Now().UnixNano())

	var restaurants []db.Restaurant

	for i := startIndex; i <= endIndex; i++ {
		restaurant := db.Restaurant{
			ID:   int32(i),
			Name: "Restaurant " + strconv.Itoa(i),
		}

		restaurants = append(restaurants, restaurant)
	}

	return restaurants
}

func initData() (*NearByRestaurantCache, LocationArgs) {
	nearByRestaurantCache := NewNearByRestaurantCache()
	// 建立一些測試資料
	args := LocationArgs{
		Lat:    24.1677759,
		Lng:    120.6654513,
		Radius: 500,
	}

	return nearByRestaurantCache, args
}

func TestNearByRestaurantCache_Append(t *testing.T) {
	t.Run("only one page when first nextPageToken given empty", func(t *testing.T) {
		nearByRestaurantCache, args := initData()
		lc := nearByRestaurantCache.checkLocationContext(args)

		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "",
			nextPageToken: "", // empty token mean there is no next page
			data:          generateMockRestaurants(1, 20),
		})
		require.Equal(t, generateMockRestaurants(1, 20), lc.listAll())

		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken2",
			nextPageToken: "",
			data:          generateMockRestaurants(21, 40),
		})
		require.Equal(t, generateMockRestaurants(1, 20), lc.listAll())

		require.Len(t, lc.listAll(), 20)
	})

	t.Run("default", func(t *testing.T) {
		nearByRestaurantCache, args := initData()
		lc := nearByRestaurantCache.checkLocationContext(args)

		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "",
			nextPageToken: "123", // empty token mean there is no next page
			data:          generateMockRestaurants(1, 20),
		})
		require.Equal(t, generateMockRestaurants(1, 20), lc.listAll())

		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "123",
			nextPageToken: "fakePageToken2",
			data:          generateMockRestaurants(21, 30),
		})
		require.Equal(t, generateMockRestaurants(1, 30), lc.listAll())

		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken2",
			nextPageToken: "",
			data:          generateMockRestaurants(31, 50),
		})
		require.Equal(t, generateMockRestaurants(1, 50), lc.listAll())

		require.Len(t, lc.listAll(), 50)
	})

	t.Run("concurrent: first append", func(t *testing.T) {
		nearByRestaurantCache, args := initData()
		lc := nearByRestaurantCache.checkLocationContext(args)

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			// time.Sleep(time.Microsecond * 1)
			nearByRestaurantCache.Append(args, PageDataOfPlaces{
				currentToken:  "",
				nextPageToken: "",
				data:          generateMockRestaurants(1, 20),
			})
			wg.Done()
		}()
		go func() {
			nearByRestaurantCache.Append(args, PageDataOfPlaces{
				currentToken:  "",
				nextPageToken: "",
				data:          generateMockRestaurants(1, 20),
			})
			wg.Done()
		}()
		wg.Wait()

		require.Equal(t, generateMockRestaurants(1, 20), lc.listAll())
		require.Len(t, lc.listAll(), 20)
	})

	t.Run("concurrent: append same token", func(t *testing.T) {
		nearByRestaurantCache, args := initData()
		lc := nearByRestaurantCache.checkLocationContext(args)

		// init append data
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "",
			nextPageToken: "fakePageToken1", // empty token mean there is no next page
			data:          generateMockRestaurants(1, 10),
		})
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken1",
			nextPageToken: "fakePageToken2", // empty token mean there is no next page
			data:          generateMockRestaurants(11, 20),
		})
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken2",
			nextPageToken: "fakePageToken3", // empty token mean there is no next page
			data:          generateMockRestaurants(21, 30),
		})
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken3",
			nextPageToken: "fakePageToken4", // empty token mean there is no next page
			data:          generateMockRestaurants(31, 40),
		})

		require.Len(t, lc.listAll(), 40)

		// test for concurrent append
		wg := sync.WaitGroup{}
		wg.Add(50)
		for i := 0; i < 50; i++ {
			if i == 25 {
				go func() {
					nearByRestaurantCache.Append(args, PageDataOfPlaces{
						currentToken:  "fakePageToken4",
						nextPageToken: "", // empty token mean there is no next page
						data:          generateMockRestaurants(41, 50),
					})
					wg.Done()
				}()
				return
			}
			go func(index int) {
				if index%3 == 0 {
					nearByRestaurantCache.Append(args, PageDataOfPlaces{
						currentToken:  "",
						nextPageToken: "fakePageToken1", // empty token mean there is no next page
						data:          generateMockRestaurants(1, 10),
					})
				} else if index%3 == 1 {
					nearByRestaurantCache.Append(args, PageDataOfPlaces{
						currentToken:  "fakePageToken1",
						nextPageToken: "fakePageToken2", // empty token mean there is no next page
						data:          generateMockRestaurants(11, 20),
					})
				} else if index%3 == 2 {
					nearByRestaurantCache.Append(args, PageDataOfPlaces{
						currentToken:  "fakePageToken2",
						nextPageToken: "fakePageToken3", // empty token mean there is no next page
						data:          generateMockRestaurants(21, 30),
					})
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		require.Len(t, lc.listAll(), 50)
	})

	t.Run("wrong append order", func(t *testing.T) {
		nearByRestaurantCache, args := initData()
		lc := nearByRestaurantCache.checkLocationContext(args)

		// init append data
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "",
			nextPageToken: "fakePageToken1", // empty token mean there is no next page
			data:          generateMockRestaurants(1, 10),
		})
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken1",
			nextPageToken: "fakePageToken2", // empty token mean there is no next page
			data:          generateMockRestaurants(11, 20),
		})

		// 先塞入一個更後面的數據
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken3",
			nextPageToken: "", // empty token mean there is no next page
			data:          generateMockRestaurants(31, 40),
		})
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken2",
			nextPageToken: "fakePageToken3", // empty token mean there is no next page
			data:          generateMockRestaurants(21, 30),
		})

		require.Equal(t, generateMockRestaurants(1, 30), lc.listAll())
		require.Len(t, lc.listAll(), 30)
	})

	t.Run("concurrent: random order append", func(t *testing.T) {
		nearByRestaurantCache, args := initData()
		lc := nearByRestaurantCache.checkLocationContext(args)

		// init append data
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "",
			nextPageToken: "fakePageToken1", // empty token mean there is no next page
			data:          generateMockRestaurants(1, 10),
		})

		wg := sync.WaitGroup{}
		wg.Add(200)
		for i := 0; i < 200; i++ {
			go func(index int) {
				if index%4 == 0 {
					nearByRestaurantCache.Append(args, PageDataOfPlaces{
						currentToken:  "",
						nextPageToken: "fakePageToken1", // empty token mean there is no next page
						data:          generateMockRestaurants(1, 10),
					})
				} else if index%4 == 1 {
					nearByRestaurantCache.Append(args, PageDataOfPlaces{
						currentToken:  "fakePageToken1",
						nextPageToken: "fakePageToken2", // empty token mean there is no next page
						data:          generateMockRestaurants(11, 20),
					})
				} else if index%4 == 2 {
					nearByRestaurantCache.Append(args, PageDataOfPlaces{
						currentToken:  "fakePageToken2",
						nextPageToken: "fakePageToken3", // empty token mean there is no next page
						data:          generateMockRestaurants(21, 30),
					})
				} else if index%4 == 3 {
					nearByRestaurantCache.Append(args, PageDataOfPlaces{
						currentToken:  "fakePageToken3",
						nextPageToken: "", // empty token mean there is no next page
						data:          generateMockRestaurants(31, 40),
					})

				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		require.Equal(t, generateMockRestaurants(1, 40), lc.listAll())
		require.Len(t, lc.listAll(), 40)
	})
}

func TestNearByRestaurantCache_GetRestaurantListByPagination(t *testing.T) {
	t.Run("no next page then append", func(t *testing.T) {
		nearByRestaurantCache, args := initData()

		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "",
			nextPageToken: "", // empty token mean there is no next page
			data:          generateMockRestaurants(1, 20),
		})
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "",
			nextPageToken: "", // empty token mean there is no next page
			data:          generateMockRestaurants(1, 20),
		})

		list, isEnough := nearByRestaurantCache.GetRestaurantListByPagination(args, 1, 20)
		require.True(t, isEnough)
		require.Equal(t, generateMockRestaurants(1, 20), list)
	})

	t.Run("enough", func(t *testing.T) {
		nearByRestaurantCache, args := initData()

		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "",
			nextPageToken: "fakePageToken1", // empty token mean there is no next page
			data:          generateMockRestaurants(1, 20),
		})
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken1",
			nextPageToken: "fakePageToken2", // empty token mean there is no next page
			data:          generateMockRestaurants(21, 40),
		})
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken2",
			nextPageToken: "", // empty token mean there is no next page
			data:          generateMockRestaurants(41, 60),
		})

		list, isEnough := nearByRestaurantCache.GetRestaurantListByPagination(args, 4, 18)
		require.True(t, isEnough)
		require.Equal(t, generateMockRestaurants(55, 60), list)
	})

	t.Run("not enough", func(t *testing.T) {
		nearByRestaurantCache, args := initData()

		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "",
			nextPageToken: "fakePageToken1", // empty token mean there is no next page
			data:          generateMockRestaurants(1, 20),
		})
		nearByRestaurantCache.Append(args, PageDataOfPlaces{
			currentToken:  "fakePageToken1",
			nextPageToken: "fakePageToken2", // empty token mean there is no next page
			data:          generateMockRestaurants(21, 40),
		})

		list, isEnough := nearByRestaurantCache.GetRestaurantListByPagination(args, 3, 18)
		require.False(t, isEnough)
		require.Equal(t, generateMockRestaurants(37, 40), list)
	})

	t.Run("for loop to append then get", func(t *testing.T) {
		nearByRestaurantCache, args := initData()

		count := 1
		var list []db.Restaurant
		var isEnough bool
		for {
			list, isEnough = nearByRestaurantCache.GetRestaurantListByPagination(args, 3, 18)
			start := (count-1)*20 + 1
			end := count * 20
			if !isEnough {
				currentToken := ""
				if count > 1 {
					currentToken = fmt.Sprintf("fakePageToken%d", count-1)
				}
				nextPageToken := fmt.Sprintf("fakePageToken%d", count)
				nearByRestaurantCache.Append(args, PageDataOfPlaces{
					currentToken:  currentToken,
					nextPageToken: nextPageToken, // empty token mean there is no next page
					data:          generateMockRestaurants(start, end),
				})
				count++
			} else {
				break
			}
		}

		require.Equal(t, 4, count)
		require.Equal(t, generateMockRestaurants(37, 54), list)
	})

	// #FIX panic: test timed out after 30s
	t.Run("concurrent: for loop to append then get", func(t *testing.T) {
		nearByRestaurantCache, args := initData()

		const LENGTH = 3
		wg := sync.WaitGroup{}
		wg.Add(LENGTH)
		for i := 0; i < LENGTH; i++ {
			go func() {
				count := 1
				var isEnough bool
				for {
					_, isEnough = nearByRestaurantCache.GetRestaurantListByPagination(args, 3, 18)
					start := (count-1)*20 + 1
					end := count * 20
					if !isEnough {
						currentToken := ""
						if count > 1 {
							currentToken = fmt.Sprintf("fakePageToken%d", count-1)
						}
						nextPageToken := fmt.Sprintf("fakePageToken%d", count)
						nearByRestaurantCache.Append(args, PageDataOfPlaces{
							currentToken:  currentToken,
							nextPageToken: nextPageToken, // empty token mean there is no next page
							data:          generateMockRestaurants(start, end),
						})
						count++
					} else {
						break
					}
				}

				require.Equal(t, 4, count)
				wg.Done()
			}()
		}

		wg.Wait()
		result, _ := nearByRestaurantCache.GetRestaurantListByPagination(args, 3, 18)
		require.Equal(t, generateMockRestaurants(37, 54), result)
	})
}
