package food_deliver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFoodDeliverApi_GetDishesFromGoogleMap(t *testing.T) {
	type args struct {
		googleMapUrl string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name:      "test1",
			args:      args{googleMapUrl: "https://maps.google.com/?cid=6123919673254163962"},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:      "test2",
			args:      args{googleMapUrl: "https://maps.google.com/?cid=643071738901551210"},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:      "foodpanda",
			args:      args{googleMapUrl: "https://maps.google.com/?cid=10740885935407955428"},
			wantCount: 6,
			wantErr:   false,
		},
		{
			name:      "ubereats",
			args:      args{googleMapUrl: "https://maps.google.com/?cid=11815448699177978476"},
			wantCount: 19,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		deliverApi := NewFoodDeliverApi()
		t.Run(tt.name, func(t *testing.T) {
			got, err := deliverApi.GetDishesFromGoogleMap(tt.args.googleMapUrl)
			require.Equal(t, err != nil, tt.wantErr)
			require.Equal(t, tt.wantCount, len(got))
		})
	}
}
