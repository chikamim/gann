package index

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/mathetake/gann/item"
)

func TestInitializeWithNormalize(t *testing.T) {
	rawItems := [][]float32{
		{2, 0},
		{0, 2},
	}
	d := 2
	idx, err := Initialize(rawItems, d, 1, 1, true)
	if err != nil {
		panic(idx)
	}

	assert.Equal(t, 2, len(idx.itemIDToItem))
	assert.Equal(t, item.Item{
		ID:  0,
		Vec: []float32{1, 0},
	}, idx.itemIDToItem[0])
	assert.Equal(t, item.Item{
		ID:  1,
		Vec: []float32{0, 1},
	}, idx.itemIDToItem[1])
}

func TestInitializeWithoutNormalize(t *testing.T) {
	rawItems := [][]float32{
		{2, 0},
		{0, 2},
	}
	d := 2
	idx, err := Initialize(rawItems, d, 1, 1, false)
	if err != nil {
		panic(idx)
	}

	assert.Equal(t, 2, len(idx.itemIDToItem))
	assert.Equal(t, item.Item{
		ID:  0,
		Vec: []float32{2, 0},
	}, idx.itemIDToItem[0])
	assert.Equal(t, item.Item{
		ID:  1,
		Vec: []float32{0, 2},
	}, idx.itemIDToItem[1])
}
