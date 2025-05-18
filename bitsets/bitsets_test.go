package bitsets_test

import (
	"testing"

	"github.com/josestg/dsa/bitsets"
	"github.com/stretchr/testify/assert"
)

func TestBitSet_BasicOps(t *testing.T) {
	b := bitsets.New(128)

	for i := 0; i < b.Len(); i++ {
		assert.False(t, b.Exists(i))
	}
	assert.Equal(t, 0, b.Count())

	b.Add(1)
	b.Add(64)
	b.Add(127)

	assert.True(t, b.Exists(1))
	assert.True(t, b.Exists(64))
	assert.True(t, b.Exists(127))
	assert.False(t, b.Exists(0))

	assert.Equal(t, 3, b.Count())

	b.Toggle(64)
	assert.False(t, b.Exists(64))
	assert.Equal(t, 2, b.Count())

	b.Toggle(64)
	assert.True(t, b.Exists(64))
	assert.Equal(t, 3, b.Count())

	b.Del(1)
	assert.False(t, b.Exists(1))
	assert.Equal(t, 2, b.Count())

	b.Reset()
	for i := 0; i < b.Len(); i++ {
		assert.False(t, b.Exists(i))
	}
	assert.Equal(t, 0, b.Count())
}

func TestBitSet_StringShape(t *testing.T) {
	b := bitsets.New(16)
	b.Add(0)
	b.Add(2)
	b.Add(15)

	s := b.String()
	assert.Equal(t, "1010000000000001000000000000000000000000000000000000000000000000", s)
}

func TestBitSet_PanicOutOfRange(t *testing.T) {
	b := bitsets.New(64)

	assert.Panics(t, func() {
		b.Add(-1)
	})
	assert.Panics(t, func() {
		b.Del(64)
	})
	assert.Panics(t, func() {
		b.Exists(999)
	})
}

func TestBitSet_Len(t *testing.T) {
	b := bitsets.New(512)
	assert.Equal(t, 512, b.Len())

	b = bitsets.New(64)
	assert.Equal(t, 64, b.Len())
}
