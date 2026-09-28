package slicex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFilterUnique(t *testing.T) {
	testCases := []struct {
		name string
		src  []int64
		want []int64
	}{
		{
			name: "过滤非正数并按首次出现去重",
			src:  []int64{3, 0, 3, -1, 2, 2, 1},
			want: []int64{3, 2, 1},
		},
		{
			name: "空切片得到空切片",
			src:  nil,
			want: []int64{},
		},
		{
			name: "全部被过滤仍是空切片",
			src:  []int64{0, -2},
			want: []int64{},
		},
		{
			name: "原顺序保留",
			src:  []int64{8, 1, 8, 4},
			want: []int64{8, 1, 4},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := New(tc.src).Filter(func(id int64) bool { return id > 0 }).Unique().Values()
			assert.Equal(t, tc.want, got)
			assert.NotNil(t, got)
		})
	}
}

func TestListDoesNotMutateSource(t *testing.T) {
	src := []int64{1, 0, 1, -1, 2}
	got := New(src).Filter(func(id int64) bool { return id > 0 }).Unique().Values()
	assert.Equal(t, []int64{1, 2}, got)
	assert.Equal(t, []int64{1, 0, 1, -1, 2}, src)
}

// order 用可比较字段演示按对象字段过滤。
type order struct {
	ID int64
	A  int
}

func TestListFilterStruct(t *testing.T) {
	testCases := []struct {
		name string
		src  []order
		want []order
	}{
		{
			name: "只保留 A 大于等于 0 并按整个对象去重",
			src: []order{
				{ID: 1, A: 3},
				{ID: 2, A: -1},
				{ID: 1, A: 3},
				{ID: 3, A: 0},
				{ID: 4, A: 2},
				{ID: 4, A: 2},
				{ID: 1, A: 1},
			},
			want: []order{
				{ID: 1, A: 3},
				{ID: 3, A: 0},
				{ID: 4, A: 2},
				{ID: 1, A: 1},
			},
		},
		{
			name: "空切片得到空切片",
			src:  nil,
			want: []order{},
		},
		{
			name: "全部被过滤仍是空切片",
			src: []order{
				{ID: 1, A: -1},
				{ID: 2, A: -2},
			},
			want: []order{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src := append([]order(nil), tc.src...)
			got := New(tc.src).
				Filter(func(item order) bool { return item.A >= 0 }).
				Unique().
				Values()
			assert.Equal(t, tc.want, got)
			assert.NotNil(t, got)
			assert.Equal(t, src, tc.src)
		})
	}
}
