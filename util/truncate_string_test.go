package util

import "testing"

func TestTruncateString(t *testing.T) {
	type args struct {
		s      string
		maxLen int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				s:      "金感蝦：評價第一無毒蝦專賣店 健康食材推薦 生態白蝦.活蝦仁.虱目魚肚.七星鱸魚排.黃金鯧.烏魚.一口烏魚子 生鮮.水產.海鮮 農遊券",
				maxLen: 17,
			},
			want: "金感蝦：評價第一無毒蝦專賣店 健康...",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateString(tt.args.s, tt.args.maxLen); got != tt.want {
				t.Errorf("TruncateString() = %v, want %v", got, tt.want)
			}
		})
	}
}
