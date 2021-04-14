package db

import (
	"reflect"
	"testing"
)

func TestAttrToQuery(t *testing.T) {
	type args struct {
		attr Attr
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []interface{}
	}{
		{
			name: "Attr转换",
			args: args{
				attr: Attr{
					OrderBy: "id desc",
				},
			},
			want:  "order by ?",
			want1: []interface{}{"id desc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := attrToQuery(tt.args.attr)
			if got != tt.want {
				t.Errorf("attrToQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("attrToQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFilterToQuery(t *testing.T) {
	type args struct {
		filter Filter
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []interface{}
	}{
		{
			name: "where转换",
			args: args{
				filter: Filter{
					"name": "daodao",
				},
			},
			want:  "where `name` = ?",
			want1: []interface{}{"daodao"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := filterToQuery(tt.args.filter)
			if got != tt.want {
				t.Errorf("filterToQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("filterToQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInsertRecordToQuery(t *testing.T) {
	type args struct {
		record Record
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []interface{}
	}{
		{
			name: "Insert转换",
			args: args{
				record: Record{
					"name": "daodao",
				},
			},
			want:  "insert into %s (`name`) values (?)",
			want1: []interface{}{"daodao"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := insertRecordToQuery(tt.args.record)
			if got != tt.want {
				t.Errorf("insertRecordToQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("insertRecordToQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
