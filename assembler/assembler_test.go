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

package assembler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssembleString1(t *testing.T) {
	tests := []struct {
		name    string
		program string
		want    []byte
	}{
		{
			name:    "1 byte positive",
			program: "pushC64 1d23 ret64",
			want:    []byte{0x10, 0x17, 0x9},
		},
		{
			name:    "normal positive",
			program: "pushC64 23 ret64",
			want:    []byte{0x10, 0x17, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9},
		},
		{
			name:    "hex",
			program: "pushC64 0xd1d2d3 ret64",
			want:    []byte{0x10, 0xd3, 0xd2, 0xd1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9},
		},
		{
			name:    "2 byte positive",
			program: "pushC64 2d23 ret64",
			want:    []byte{0x10, 0x17, 0x0, 0x9},
		},
		{
			name:    "2 byte positive big",
			program: "pushC64 2d65535 ret64",
			want:    []byte{0x10, 0xff, 0xff, 0x9},
		},
		{
			name:    "3 byte positive big",
			program: "pushC64 3d65535 ret64",
			want:    []byte{0x10, 0xff, 0xff, 0x0, 0x9},
		},
		{
			name:    "2 byte negative",
			program: "pushC64 2d-2 ret64",
			want:    []byte{0x10, 0xfe, 0xff, 0x9},
		},
		{
			name:    "2 byte small negative",
			program: "pushC64 2d-32768 ret64",
			want:    []byte{0x10, 0x0, 0x80, 0x9},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := AssembleString(testCase.program)
			assert.Equal(t, testCase.want, got)
		})
	}
}
