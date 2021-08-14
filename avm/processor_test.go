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

package avm

import (
	"fmt"
	"go-AVM/avm/binary"
	"math"
	"testing"
)

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
		operandStack: newOperandStack(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.current.operandStack.ensureLen(64 * 1024)
		for j := 0; j < 1024*4; j++ {
			p.iAdd()
		}
	}
}

func BenchmarkProcessor_iAdd64NoFunc(b *testing.B) {
	p := Processor{}
	p.current = &CallInfo{
		operandStack: newOperandStack(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.current.operandStack.ensureLen(64 * 1024)
		for j := 0; j < 1024*4; j++ {
			p.iAddNoFunc()
		}
	}
}

func (p *Processor) iAddNoFunc() {
	top := p.current.operandStack.length()
	binary.PutInt64(p.current.operandStack.content,
		top-16, binary.ReadInt64(p.current.operandStack.content, top-8)+binary.ReadInt64(p.current.operandStack.content, top-16))
	p.current.operandStack.shrinkTo(top - 8)
}
