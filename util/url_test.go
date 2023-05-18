package util

import "testing"

func Test_BindUrl(t *testing.T) {
	type args struct {
		apiBaseUrl string
		endPoint   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				apiBaseUrl: "http://localhost:8080",
				endPoint:   "callback",
			},
			want:    "http://localhost:8080/callback",
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				apiBaseUrl: "http://localhost:8080/",
				endPoint:   "/callback",
			},
			want:    "http://localhost:8080/callback",
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				apiBaseUrl: "http://localhost:8080/",
				endPoint:   "callback",
			},
			want:    "http://localhost:8080/callback",
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				apiBaseUrl: "http://localhost:8080",
				endPoint:   "/callback",
			},
			want:    "http://localhost:8080/callback",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BindUrl(tt.args.apiBaseUrl, tt.args.endPoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("bindUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bindUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
