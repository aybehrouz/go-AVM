package avm

import (
	"fmt"
	"go-AVM/avm/binary"
	"math"
	"testing"
)

/*
func TestProcessor_iAdd64(t *testing.T) {
	p := Processor{
		operandStack: newOperandStack(),
	}
	p.operandStack.ensureLen(16)
	p.operandStack.content[0] = 3
	p.operandStack.content[8] = 4
	p.iAdd64()
	fmt.Printf("%v", p.operandStack)
	println(cap(p.operandStack.content))
}
*/

const Log10 = 3.321928

func TestFloat(t *testing.T) {
	f1 := 100000000000.123
	f2 := 125.33333333333333333
	a := math.Float64bits(f1)
	b := math.Float64bits(f2)
	bt := truncate(b, int(math.Ceil(Log10*4)))
	fmt.Printf("ft: %v \n", math.Float64frombits(bt))
	fmt.Printf("%v\n", add(a, bt))
	fmt.Printf("%v", f1+f2)
}
func BenchmarkProcessor_iAdd64(b *testing.B) {
	p := Processor{}
	p.current = &CallInfo{
		operandStack: newOperandStack().ensureLen(64 * 1024),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.iAdd64()
	}
}

func BenchmarkCopy(b *testing.B) {
	bytes := make([][8]byte, 500)
	dst := make([]byte, 8)
	for i := 0; i < b.N; i++ {
		for j := 0; j < 490; j++ {
			binary.CopyBytes(dst, 0, bytes[j][:], 0, 4)
			// copy(dst, bytes[j][:])
		}
	}
}
