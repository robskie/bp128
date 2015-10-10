package bp128

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxBits128_32(t *testing.T) {
	in := makeAlignedSlice([]uint32{}, 128)
	for i := 0; i < 32; i++ {
		in.Index(127).SetUint(1 << uint(i))
		bs := maxBits128_32(in.Pointer(), 0, new(byte))

		if !assert.EqualValues(t, i+1, bs) {
			break
		}
	}
}

func TestMaxBits128_64(t *testing.T) {
	in := makeAlignedSlice([]uint64{}, 128)
	for i := 0; i < 64; i++ {
		in.Index(127).SetUint(1 << uint(i))
		bs := maxBits128_64(in.Pointer(), 0, new(byte))

		if !assert.EqualValues(t, i+1, bs) {
			break
		}
	}
}

func TestDMaxBits128_32(t *testing.T) {
	offset := 3
	in := makeAlignedSlice([]uint32{}, 128)
	for i := 0; i < 32-offset+1; i++ {
		delta := 1 << uint(i)

		v := 0
		for i := 0; i < 128; i++ {
			v += delta
			in.Index(i).SetUint(uint64(v))
		}

		seed := makeAlignedBytes(16)
		copy(seed, convertToBytes(32, in.Slice(0, 4)))

		bs := dmaxBits128_32(in.Pointer(), 0, &seed[0])
		if !assert.EqualValues(t, i+offset, bs) {
			break
		}
	}
}

func TestDMaxBits128_64(t *testing.T) {
	offset := 2
	in := makeAlignedSlice([]uint64{}, 128)
	for i := 0; i < 64-offset+1; i++ {
		delta := 1 << uint(i)

		v := 0
		for i := 0; i < 128; i++ {
			v += delta
			in.Index(i).SetUint(uint64(v))
		}

		seed := makeAlignedBytes(16)
		copy(seed, convertToBytes(64, in.Slice(0, 2)))

		bs := dmaxBits128_64(in.Pointer(), 0, &seed[0])
		if !assert.EqualValues(t, i+offset, bs) {
			break
		}
	}
}
