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
			name:    "positive number",
			program: "pushC64 23 ret64",
			want:    []byte{0x10, 0x17, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9},
		},
		{
			name:    "negative number",
			program: "pushC64 -254 ret64",
			want:    []byte{0x10, 0x2, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x9},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := AssembleString(testCase.program)
			assert.Equal(t, testCase.want, got)
		})
	}
}
