package bp128

import (
	"reflect"
	"unsafe"
)

func isAligned(intSize int, addr uintptr, index int) bool {
	addr += uintptr(index * (intSize / 8))
	return addr&15 == 0
}

func makeAlignedBytes(length int) []byte {
	const alignment = 16
	bytes := make([]byte, length+alignment)

	idx := 0
	addr := unsafe.Pointer(&bytes[0])
	for !isAligned(8, uintptr(addr), idx) {
		idx++
	}

	return bytes[idx : idx+length]
}

func makeAlignedSlice(slice interface{}, n int) reflect.Value {
	intSize := 0
	switch slice.(type) {
	case []int, []uint, []int64, []uint64:
		intSize = 64
	case []int32, []uint32:
		intSize = 32
	case []int16, []uint16:
		intSize = 16
	case []int8, []uint8:
		intSize = 8
	default:
		panic("bp128: input is not an integer slice")
	}

	const alignment = 16
	offset := (alignment * 8) / intSize

	c := n + offset
	vslice := reflect.MakeSlice(reflect.TypeOf(slice), c, c)

	idx := 0
	addr := unsafe.Pointer(vslice.Pointer())
	for !isAligned(intSize, uintptr(addr), idx) {
		idx++
	}

	return vslice.Slice(idx, idx+n)
}

func alignSlice(intSize int, v reflect.Value) reflect.Value {
	const alignment = 16
	offset := (alignment * 8) / intSize

	nslice := v
	length := v.Len() + offset
	if v.Cap() < length {
		nslice = reflect.MakeSlice(v.Type(), length, length)
	}

	idx := 0
	addr := unsafe.Pointer(nslice.Pointer())
	for !isAligned(intSize, uintptr(addr), idx) {
		idx++
	}

	return reflect.AppendSlice(nslice.Slice(idx, idx), v)
}

func convertToBytes(intSize int, v reflect.Value) []byte {
	if !v.IsValid() {
		return nil
	}

	nbytes := intSize / 8
	sh := &reflect.SliceHeader{}
	sh.Cap = v.Cap() * nbytes
	sh.Len = v.Len() * nbytes
	sh.Data = v.Pointer()
	return *(*[]uint8)(unsafe.Pointer(sh))
}

func appendBytes(intSize int, v reflect.Value, b []byte) reflect.Value {
	length := (len(b) * 8) / intSize

	sh := &reflect.SliceHeader{}
	sh.Cap = length
	sh.Len = length
	sh.Data = uintptr(unsafe.Pointer(&b[0]))
	nslice := reflect.NewAt(v.Type(), unsafe.Pointer(sh)).Elem()

	return reflect.AppendSlice(v, nslice)
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
