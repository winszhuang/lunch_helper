package thirdparty

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommonFoodDeliverApi_GetDishes(t *testing.T) {
	deliverApi := NewCommonFoodDeliverApi()

	testCases := []struct {
		name                 string
		pageUrl              string
		hasError             bool
		isLenGreaterThanZero bool
	}{
		{
			name:                 "ubereats url",
			pageUrl:              "https://www.ubereats.com/tw/store/%E5%BA%B7%E5%AF%B6%E8%97%A5%E7%87%89%E6%8E%92%E9%AA%A8-%E7%BE%8E%E5%BE%B7%E7%B8%BD%E5%BA%97/-NQAJwmsRJGIo-7WUwHfHQ?utm_campaign=place-action-link&utm_medium=organic&utm_source=google",
			hasError:             false,
			isLenGreaterThanZero: true,
		},
		{
			name:                 "foodpanda url",
			pageUrl:              "https://www.foodpanda.com.tw/restaurant/gbnl/chen-xi-liang-mian-tai-zhong-bei-ping-dian",
			hasError:             false,
			isLenGreaterThanZero: true,
		},
		{
			name:                 "not ubereats or foodpanda",
			pageUrl:              "https://bistro-257.business.site/",
			hasError:             true,
			isLenGreaterThanZero: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := deliverApi.GetDishes(tc.pageUrl)
			require.Equal(t, tc.hasError, err != nil)
			require.Equal(t, tc.isLenGreaterThanZero, len(result) > 0)
		})
	}
}
