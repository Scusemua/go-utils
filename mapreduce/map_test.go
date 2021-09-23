package mapreduce

import (
	"reflect"
	"strconv"
	"testing"
)

func TestParallelMap(t *testing.T) {
	type args struct {
		arr    interface{}
		mapper interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Argument is not an array",
			args:    args{arr: 1, mapper: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "mapper function is nil",
			args:    args{arr: []int{1, 2, 3}, mapper: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "mapper is not a function",
			args:    args{arr: []int{1, 2, 3}, mapper: 1},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Valid mapper",
			args: args{arr: []int{1, 2, 3}, mapper: func(num int) int {
				return num + 1
			}},
			want:    []int{2, 3, 4},
			wantErr: false,
		},
		{
			name: "Valid mapper",
			args: args{arr: []int{1, 2, 3}, mapper: func(num int) string {
				return strconv.Itoa(num)
			}},
			want:    []string{"1", "2", "3"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParallelMap(tt.args.mapper, tt.args.arr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Map() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
