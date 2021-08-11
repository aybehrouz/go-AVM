package prefix

import "unsafe"

type Identifier64 uint64

func readIdentifier64(b []byte) (id Identifier64) {
	length := int(b[0])
	if length > 8 {
		panic("identifier is more than 8 bytes")
	}
	if len(b) > 8 {
		id = Identifier64(*(*uint64)(unsafe.Pointer(&b[1])))
	} else {
		padded := make([]byte, 9)
		copy(padded, b)
		id = Identifier64(*(*uint64)(unsafe.Pointer(&padded[1])))
	}
	id <<= (8 - length) << 8
	return
}
