package binary

import (
	"go-AVM/avm/prefix"
	"unsafe"
)

func init() {
	b := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88}
	read64 := *(*uint64)(unsafe.Pointer(&b[0]))
	if read64 != 0x8877665544332211 {
		panic("platform not supported")
	}
	b = []byte{0x11, 0x22, 0x33, 0x44}
	read32 := *(*uint32)(unsafe.Pointer(&b[0]))
	if read32 != 0x44332211 {
		panic("platform not supported")
	}
}

func ReadInt64(src []byte, offset int64) int64 {
	_ = src[offset+7] // bounds check hint to compiler; see golang.org/issue/14808
	return *(*int64)(unsafe.Pointer(&src[offset]))
}

func PutInt64(dst []byte, offset int64, v int64) {
	_ = dst[offset+7]
	*(*int64)(unsafe.Pointer(&dst[offset])) = v
}

func ReadUint16(src []byte, offset int64) uint16 {
	_ = src[offset+1]
	return *(*uint16)(unsafe.Pointer(&src[offset]))
}

func ReadInt32(b []byte, offset int64) int32 {
	_ = b[offset+3]
	return *(*int32)(unsafe.Pointer(&b[offset]))
}

func readFloat64(b []byte, offset int64) float64 {
	_ = b[offset+7]
	return *(*float64)(unsafe.Pointer(&b[offset]))
}

func ReadIdentifier64(src []byte, offset int64) prefix.Identifier64 {
	_ = src[offset+7]
	return *(*prefix.Identifier64)(unsafe.Pointer(&src[offset]))
}

func Copy64(dst []byte, dstOffset int64, src []byte, srcOffset int64) {
	_ = src[srcOffset+7]
	_ = dst[dstOffset+7]
	*(*uint64)(unsafe.Pointer(&dst[dstOffset])) = *(*uint64)(unsafe.Pointer(&src[srcOffset]))
}

func Copy32(dst []byte, dstOffset int64, src []byte, srcOffset int64) {
	_ = src[srcOffset+3]
	_ = dst[dstOffset+3]
	*(*uint32)(unsafe.Pointer(&dst[dstOffset])) = *(*uint32)(unsafe.Pointer(&src[srcOffset]))
}

// CopyBytes always copies exactly n bytes
func CopyBytes(dst []byte, dstOffset int64, src []byte, srcOffset int64, n int64) {
	if n == 8 {
		Copy64(dst, dstOffset, src, srcOffset)
	} else {
		copy(dst[dstOffset:dstOffset+n], src[srcOffset:srcOffset+n])
	}
}
