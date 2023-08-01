package week01

import (
	"reflect"
	"testing"
)

var deleteTests = []struct {
	s    []int
	i, j int
	want []int
}{
	{
		[]int{1, 2, 3},
		0,
		0,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0,
		1,
		[]int{2, 3},
	},
	{
		[]int{1, 2, 3},
		3,
		3,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0,
		2,
		[]int{3},
	},
	{
		[]int{1, 2, 3},
		0,
		3,
		[]int{},
	},
}

func TestDelete(t *testing.T) {
	for _, test := range deleteTests {
		copyObj := Clone(test.s)
		if got := Delete(copyObj, test.i, test.j); !Equal(got, test.want) {
			t.Errorf("Delete(%v, %d, %d) = %v, want %v", test.s, test.i, test.j, got, test.want)
		}
	}
}

var deleteAndShrinkTests = []struct {
	s       []int
	i, j    int
	want    []int
	wantCap int
}{
	{
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		8,
		16,
		[]int{1, 2, 3, 4, 5, 6, 7, 8},
		16, // cap shrinks from 32 to 16, which is 2*len
	},
	{
		[]int{1, 2, 3, 4, 5, 6, 7, 8},
		0,
		1,
		[]int{2, 3, 4, 5, 6, 7, 8},
		8, // cap stays the same
	},
	{
		[]int{1, 2, 3},
		0,
		2,
		[]int{3},
		3, // cap stays the same
	},
	{
		[]int{1, 2, 3, 4},
		0,
		4,
		[]int{},
		4, // cap stays the same
	},
}

func TestDeleteAndShrink(t *testing.T) {
	for _, test := range deleteAndShrinkTests {
		copyObj := Clone(test.s)
		got := DeleteAndShrink(copyObj, test.i, test.j)
		if !reflect.DeepEqual(got, test.want) || cap(got) != test.wantCap {
			t.Errorf("DeleteAndShrink(%v, %d, %d) = %v, cap = %d, want %v, cap = %d", test.s, test.i, test.j, got, cap(got), test.want, test.wantCap)
		}
	}
}
