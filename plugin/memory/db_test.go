package main

import (
	"reflect"
	"testing"
)

func Test_db_getFromDb(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		d       db
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "get key",
			args:    args{id: "test"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFromDb() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFromDb() got = %v, want %v", got, tt.want)
			}
		})
	}
}
