package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntervalArray(t *testing.T) {
	iv := PartitionList[int]{
		{0, 1},
		{10, 2},
		{50, 3},
	}
	assert.Equal(t, 1, iv.GetItem(0))
	assert.Equal(t, 1, iv.GetItem(9))
	assert.Equal(t, 2, iv.GetItem(10))
	assert.Equal(t, 2, iv.GetItem(49))
	assert.Equal(t, 3, iv.GetItem(50))
	assert.Equal(t, 3, iv.GetItem(100000))

	iv.AddRecord(5, 4)
	assert.Equal(t, 1, iv.GetItem(0))
	assert.Equal(t, 1, iv.GetItem(4))
	assert.Equal(t, 4, iv.GetItem(5))
	assert.Equal(t, 4, iv.GetItem(9))
	assert.Equal(t, 2, iv.GetItem(10))
	assert.Equal(t, 2, iv.GetItem(49))
	assert.Equal(t, 3, iv.GetItem(50))
	assert.Equal(t, 3, iv.GetItem(100000))

	iv.AddRecord(20, 5)
	assert.Equal(t, 1, iv.GetItem(0))
	assert.Equal(t, 1, iv.GetItem(4))
	assert.Equal(t, 4, iv.GetItem(5))
	assert.Equal(t, 4, iv.GetItem(9))
	assert.Equal(t, 2, iv.GetItem(10))
	assert.Equal(t, 2, iv.GetItem(19))
	assert.Equal(t, 5, iv.GetItem(20))
	assert.Equal(t, 5, iv.GetItem(49))
	assert.Equal(t, 3, iv.GetItem(50))
	assert.Equal(t, 3, iv.GetItem(100000))
}
