package util

import (
	"github.com/mmcloughlin/geohash"
)

// 精確度到9位數(4.8m x 4.8m)
const PRECISION_LENGTH = 9

func ToGeoHash(lat, lng float64) string {
	return geohash.EncodeWithPrecision(lat, lng, PRECISION_LENGTH)
}
