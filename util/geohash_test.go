package util

import "testing"

func TestToGeoHash(t *testing.T) {
	type args struct {
		lat float64
		lng float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				lat: 24.179697021653315,
				lng: 120.68577441170953,
			},
			// https://www.movable-type.co.uk/scripts/geohash.html
			want: "wsmcd3b8y",
		},
		{
			name: "test2",
			args: args{
				lat: 24.1677759,
				lng: 120.6654513,
			},
			want: "wsmc3z9gk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToGeoHash(tt.args.lat, tt.args.lng); got != tt.want {
				t.Errorf("ToGeoHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
