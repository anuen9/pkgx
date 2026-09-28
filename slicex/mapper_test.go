package slicex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type orderModel struct {
	ID int64
	A  int
}

type orderDomain struct {
	ID     int64
	Amount int
}

func TestMapper(t *testing.T) {
	toDomain := Func[orderModel, orderDomain](func(item orderModel) orderDomain {
		return orderDomain{ID: item.ID, Amount: item.A}
	})
	testCases := []struct {
		name string
		src  []orderModel
		want []orderDomain
	}{
		{
			name: "按字段转换成另一种对象并保留顺序",
			src: []orderModel{
				{ID: 1, A: 3},
				{ID: 2, A: 0},
				{ID: 4, A: -1},
			},
			want: []orderDomain{
				{ID: 1, Amount: 3},
				{ID: 2, Amount: 0},
				{ID: 4, Amount: -1},
			},
		},
		{
			name: "空切片得到空切片",
			src:  nil,
			want: []orderDomain{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Mapper(tc.src, toDomain)
			assert.Equal(t, tc.want, got)
			assert.NotNil(t, got)
		})
	}
}

func TestMapperDoesNotMutateSource(t *testing.T) {
	src := []orderModel{{ID: 1, A: 3}, {ID: 2, A: -1}}
	got := Mapper(src, func(item orderModel) orderDomain {
		return orderDomain{ID: item.ID, Amount: item.A}
	})
	assert.Equal(t, []orderDomain{{ID: 1, Amount: 3}, {ID: 2, Amount: -1}}, got)
	assert.Equal(t, []orderModel{{ID: 1, A: 3}, {ID: 2, A: -1}}, src)
}
