package vos

import "testing"

func Test_home(t *testing.T) {
	type args struct {
		user string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "user-root",
			args: args{
				user: "root",
			},
			want: "/root",
		},
		{
			name: "user-sakura",
			args: args{
				user: "sakura",
			},
			want: "/home/sakura",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := home(tt.args.user); got != tt.want {
				t.Errorf("home() = %v, want %v", got, tt.want)
			}
		})
	}
}
