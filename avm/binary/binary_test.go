package binary

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPutInt64(t *testing.T) {
	assert := assert.New(t)
	b := make([]byte, 15)

	PutInt64(b, 3, 0x1122334455667788)
	want := []byte{0x0, 0x0, 0x0, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x0, 0x0, 0x0, 0x0}
	assert.Equal(want, b)

	want = []byte{0x0, 0x0, 0x0, 0x88, 0x77, 0x66, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	PutInt64(b, 7, 0)
	assert.Equal(want, b)

	assert.Panics(func() { PutInt64(b, 8, 0) })

	assert.Panics(func() { PutInt64(b, -1, 0) })
}

func TestReadInt64(t *testing.T) {
	assert := assert.New(t)
	b := make([]byte, 12)

	var want int64 = 0
	assert.Equal(want, ReadInt64(b, 0))

	want = 123456789101112
	PutInt64(b, 0, want)
	assert.Equal(want, ReadInt64(b, 0))

	want = -25698745
	PutInt64(b, 4, want)
	assert.Equal(want, ReadInt64(b, 4))

	assert.Panics(func() { ReadInt64(b, 5) })

	assert.Panics(func() { ReadInt64(b, -1) })
}

func TestCopyBytes(t *testing.T) {

}

func BenchmarkCopy(b *testing.B) {
	src := make([]byte, 5*1024)
	dst := make([]byte, 5*1024)

	b.Run("unsafe", func(bb *testing.B) {
		bb.ResetTimer()
		for i := 0; i < bb.N; i++ {
			for i := 0; i < len(src)-20; i++ {
				Copy64(dst, int64(i), src, int64(i))
			}
		}
	})

	b.Run("go copy", func(bb *testing.B) {
		bb.ResetTimer()
		for i := 0; i < bb.N; i++ {
			for i := 0; i < len(src)-20; i++ {
				copy(dst[i:], src[i:i+8])
			}
		}
	})
}
