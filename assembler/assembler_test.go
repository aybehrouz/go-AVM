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
