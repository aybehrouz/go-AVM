// Copyright (c) 2021 aybehrouz <behrouz_ayati@yahoo.com>. This file is
// part of the go-avm repository: the Go implementation of the Argennon
// Virtual Machine (AVM).
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program. If not, see <https://www.gnu.org/licenses/>.

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
	assert := assert.New(t)
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	dst := make([]byte, 15)

	CopyBytes(dst, 2, src, 3, 8)
	want := []byte{0x0, 0x0, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0x0, 0x0, 0x0, 0x0, 0x0}
	assert.Equal(want, dst)

	CopyBytes(dst, 0, src, 11, 1)
	want = []byte{0xc, 0x0, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0x0, 0x0, 0x0, 0x0, 0x0}
	assert.Equal(want, dst)

	CopyBytes(dst, 1, src, 11, 0)
	assert.Equal(want, dst)

	CopyBytes(dst, 0, src, 0, 12)
	want = []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0x0, 0x0, 0x0}
	assert.Equal(want, dst)

	CopyBytes(dst, 3, src, 0, 12)
	want = []byte{0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc}
	assert.Equal(want, dst)

	assert.Panics(func() {
		CopyBytes(dst, 0, src, 0, 15)
	})

	assert.Panics(func() {
		CopyBytes(dst, 8, src, 0, 8)
	})

	assert.Panics(func() {
		CopyBytes(dst, 0, src, 5, 8)
	})

	assert.Panics(func() {
		CopyBytes(dst, -1, src, 0, 8)
	})

	assert.Panics(func() {
		CopyBytes(dst, 15, src, 0, 1)
	})

	assert.Panics(func() {
		CopyBytes(dst, 0, src, -1, 1)
	})
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
