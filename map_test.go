package constrainediterable

import (
	"reflect"
	"testing"
	"time"
)

func TestNewMap(t *testing.T) {
	type args struct {
		size int
		age  time.Duration
	}
	tests := []struct {
		name string
		args args
		want *Map
	}{{
		name: "BasicConstructor",
		args: args{size: 10, age: 1 * time.Minute},
		want: &Map{ageLimit: 1 * time.Minute, sizeLimit: 10, innerMap: map[string]*internalElement{}, sortedKeys: []string{}},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMap(tt.args.size, tt.args.age); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMapFromMap(t *testing.T) {
	type args struct {
		size int
		age  time.Duration
		m    map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *Map
	}{
		{
			name: "BasicConstructor",
			args: args{size: 10, age: 1 * time.Minute, m: map[string]interface{}{"akey": "avalue"}},
			want: &Map{
				ageLimit:   1 * time.Minute,
				sizeLimit:  10,
				innerMap:   map[string]*internalElement{"akey": {value: "avalue"}},
				sortedKeys: []string{"akey"}},
		}, {
			name: "ConstructWithOverflow",
			args: args{size: 2, age: 1 * time.Minute, m: map[string]interface{}{
				"akey":  "avalue",
				"akey1": "avalue1",
				"akey2": "avalue2"}},
			want: &Map{
				ageLimit:  1 * time.Minute,
				sizeLimit: 2,
				innerMap: map[string]*internalElement{
					"akey":  {value: "avalue"},
					"akey1": {value: "avalue1"},
				},
				sortedKeys: []string{"akey"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMapFromMap(tt.args.size, tt.args.age, tt.args.m)
			if len(got.innerMap) != len(tt.want.innerMap) {
				t.Errorf("internal map differs from passed in lenght got: %d want: %d", len(got.innerMap), len(tt.want.innerMap))
			}
			for k := range got.innerMap {
				if got.innerMap[k].value != tt.want.innerMap[k].value {
					t.Errorf("values in Map differ from map key: %s got: %s want: %s", k, got.innerMap[k].value, tt.want.innerMap[k].value)
				}
			}
		})
	}
}

func TestMap_Add(t *testing.T) {
	type fields struct {
		innerMap   map[string]*internalElement
		sizeLimit  int
		ageLimit   time.Duration
		sortedKeys []string
	}
	type args struct {
		k string
		v interface{}
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		expectedSortedKeys []string
	}{
		{name: "Add Basic",
			fields: fields{
				innerMap:   map[string]*internalElement{},
				sortedKeys: []string{},
				ageLimit:   0,
				sizeLimit:  10,
			},
			args:               args{k: "1", v: "1"},
			expectedSortedKeys: []string{"1"},
		},
		{name: "Add Triggers Evict N",
			fields: fields{
				innerMap: map[string]*internalElement{"0": &internalElement{
					value:        "0",
					creationTime: time.Now().Add(-10 * time.Second),
				}},
				sortedKeys: []string{"0"},
				ageLimit:   0,
				sizeLimit:  1,
			},
			args:               args{k: "1", v: "1"},
			expectedSortedKeys: []string{"1"},
		},
		{name: "Add Triggers Evict Old",
			fields: fields{
				innerMap: map[string]*internalElement{
					"0": &internalElement{
						value:        "0",
						creationTime: time.Now().Add(-10 * time.Minute),
					},
					"1": &internalElement{
						value:        "1",
						creationTime: time.Now().Add(-5 * time.Minute),
					},
					"2": &internalElement{
						value:        "2",
						creationTime: time.Now().Add(-10 * time.Second),
					},
				},
				sortedKeys: []string{"0", "1", "2"},
				ageLimit:   1 * time.Minute,
				sizeLimit:  3,
			},
			args:               args{k: "3", v: "3"},
			expectedSortedKeys: []string{"3", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				innerMap:   tt.fields.innerMap,
				sizeLimit:  tt.fields.sizeLimit,
				ageLimit:   tt.fields.ageLimit,
				sortedKeys: tt.fields.sortedKeys,
			}
			m.Add(tt.args.k, tt.args.v)
			if !reflect.DeepEqual(m.sortedKeys, tt.expectedSortedKeys) {
				t.Errorf("keys in store are not what we expected got: %v, want: %v", m.sortedKeys, tt.expectedSortedKeys)
			}
		})
	}
}

func TestMap_evictOld(t *testing.T) {
	type fields struct {
		innerMap   map[string]*internalElement
		sizeLimit  int
		ageLimit   time.Duration
		sortedKeys []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				innerMap:   tt.fields.innerMap,
				sizeLimit:  tt.fields.sizeLimit,
				ageLimit:   tt.fields.ageLimit,
				sortedKeys: tt.fields.sortedKeys,
			}
			m.evictOld()
		})
	}
}

func TestMap_evictN(t *testing.T) {
	type fields struct {
		innerMap   map[string]*internalElement
		sizeLimit  int
		ageLimit   time.Duration
		sortedKeys []string
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				innerMap:   tt.fields.innerMap,
				sizeLimit:  tt.fields.sizeLimit,
				ageLimit:   tt.fields.ageLimit,
				sortedKeys: tt.fields.sortedKeys,
			}
			m.evictN(tt.args.n)
		})
	}
}

func TestMap_evict(t *testing.T) {
	type fields struct {
		innerMap   map[string]*internalElement
		sizeLimit  int
		ageLimit   time.Duration
		sortedKeys []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				innerMap:   tt.fields.innerMap,
				sizeLimit:  tt.fields.sizeLimit,
				ageLimit:   tt.fields.ageLimit,
				sortedKeys: tt.fields.sortedKeys,
			}
			m.evict()
		})
	}
}
