package food_deliver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFoodDeliverApi_CheckFoodDeliverFromGoogleMap(t *testing.T) {
	api := NewFoodDeliverApi()

	testCases := []struct {
		name         string
		googleMapURL string
		want         *FetchInfo
		wantErr      bool
	}{
		{
			name:         "no deliver case",
			googleMapURL: "https://maps.google.com/?cid=6123919673254163962",
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "foodpanda case",
			googleMapURL: "https://maps.google.com/?cid=10740885935407955428",
			want: &FetchInfo{
				DeliverName: FoodPanda,
				FetchLink:   "https://www.foodpanda.com.tw/restaurant/x1xg/luo-ji-niu-pai?utm_source\\\\u003dgoogle\\\\u0026utm_medium\\\\u003dorganic\\\\u0026utm_campaign\\\\u003dgoogle_reserve_place_order_action\\",
			},
			wantErr: false,
		},
		{
			name:         "ubereats case",
			googleMapURL: "https://maps.google.com/?cid=11815448699177978476",
			want: &FetchInfo{
				DeliverName: UberEats,
				FetchLink:   "https://www.ubereats.com/tw/store/大煙囪烤肉便當/PI7R2L-ZS_eLIQ7Y0Dndaw?utm_campaign\\\\u003dplace-action-link\\\\u0026utm_medium\\\\u003dorganic\\\\u0026utm_source\\\\u003dgoogle\\",
			},
			wantErr: false,
		},
		{
			name:         "InvalidURL",
			googleMapURL: "https://maps.google.com/invalid",
			want:         nil,
			wantErr:      true,
		},
		// 添加更多測試用例
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := api.CheckFoodDeliverFromGoogleMap(tc.googleMapURL)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestFoodDeliverApi_GetDishes(t *testing.T) {
	api := NewFoodDeliverApi()

	testCases := []struct {
		name      string    // 測試用例名稱
		args      FetchInfo // 測試用例參數
		wantCount int       // 期望的Dishes數量
		wantErr   bool      // 期望是否發生錯誤
	}{
		{
			name: "empty fetch info",
			args: FetchInfo{
				DeliverName: "",
				FetchLink:   "",
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "no deliverName",
			args: FetchInfo{
				DeliverName: "",
				FetchLink:   "https://www.foodpanda.com.tw/restaurant/x1xg/luo-ji-niu-pai?utm_source\\\\u003dgoogle\\\\u0026utm_medium\\\\u003dorganic\\\\u0026utm_campaign\\\\u003dgoogle_reserve_place_order_action\\",
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "fetch foodpanda but give ubereats deliver name",
			args: FetchInfo{
				DeliverName: UberEats,
				FetchLink:   "https://www.foodpanda.com.tw/restaurant/x1xg/luo-ji-niu-pai?utm_source\\\\u003dgoogle\\\\u0026utm_medium\\\\u003dorganic\\\\u0026utm_campaign\\\\u003dgoogle_reserve_place_order_action\\",
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "fetch ubereats but give foodpanda deliver name",
			args: FetchInfo{
				DeliverName: FoodPanda,
				FetchLink:   "https://www.ubereats.com/tw/store/大煙囪烤肉便當/PI7R2L-ZS_eLIQ7Y0Dndaw?utm_campaign\\\\u003dplace-action-link\\\\u0026utm_medium\\\\u003dorganic\\\\u0026utm_source\\\\u003dgoogle\\",
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "ubereats fetch",
			args: FetchInfo{
				DeliverName: UberEats,
				FetchLink:   "https://www.ubereats.com/tw/store/大煙囪烤肉便當/PI7R2L-ZS_eLIQ7Y0Dndaw?utm_campaign\\\\u003dplace-action-link\\\\u0026utm_medium\\\\u003dorganic\\\\u0026utm_source\\\\u003dgoogle\\",
			},
			wantCount: 19,
			wantErr:   false,
		},
		{
			name: "foodpanda fetch",
			args: FetchInfo{
				DeliverName: FoodPanda,
				FetchLink:   "https://www.foodpanda.com.tw/restaurant/u0gl/ming-zhi-gao-xian-xian-chao-guan?utm_source=google&utm_medium=organic&utm_campaign=google_reserve_place_order_action",
			},
			wantCount: 39,
			wantErr:   false,
		},
		{
			name: "InvalidURL",
			args: FetchInfo{
				DeliverName: FoodPanda,
				FetchLink:   "https://maps.google.com/invalid",
			},
			wantCount: 0,
			wantErr:   true,
		},
		// 添加更多測試用例
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := api.GetDishes(&tc.args)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, got, tc.wantCount)
		})
	}
}
