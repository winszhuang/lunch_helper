package util

import "testing"

func TestParseId(t *testing.T) {
	type args struct {
		key  string
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				key:  "restaurantmenu",
				data: "/restaurantmenu=12",
			},
			want:    12,
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				key:  "foodlike",
				data: "/foodlike=12",
			},
			want:    12,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseId(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseId() = %v, want %v", got, tt.want)
			}
		})
	}
}
