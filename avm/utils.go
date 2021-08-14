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
	"go-AVM/avm/binary"
	"go-AVM/avm/prefix"
)

type dynamicArray struct {
	content []byte
	maxSize int64
}

func (da *dynamicArray) shrinkTo(length int64) {
	da.content = da.content[:length]
}

func (da *dynamicArray) ensureLen(length int64) *dynamicArray {
	if length <= int64(len(da.content)) {
		return da
	}
	if length <= int64(cap(da.content)) {
		da.content = da.content[:length]
		return da
	}
	if length > da.maxSize {
		panic(MemoryLimitExceeded)
	}

	b := make([]byte, min(da.maxSize, 2*length))
	copy(b, da.content)
	da.content = b[:length]
	return da
}

func (da *dynamicArray) length() int64 {
	return int64(len(da.content))
}

func newOperandStack() *dynamicArray {
	return &dynamicArray{
		content: make([]byte, 0, InitialOpStackSize),
		maxSize: MaxOpStackSize,
	}
}

func newLocalFrame() *dynamicArray {
	return &dynamicArray{
		content: make([]byte, InitialLocalFrameSize),
		maxSize: MaxLocalFrameSize,
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func (p *Processor) popIdentifier64() prefix.Identifier64 {
	top := p.current.operandStack.length()
	id := binary.ReadIdentifier64(p.current.operandStack.content, top-8)
	p.current.operandStack.shrinkTo(top - 8)
	return id
}

func (p *Processor) readConst16() uint16 {
	c := p.methodArea.LoadUint16(p.current.pc)
	p.current.pc += 2
	return c
}

func (p *Processor) peekInt64() (a int64, b int64, top int64) {
	top = p.current.operandStack.length()
	a = binary.ReadInt64(p.current.operandStack.content, top-8)
	b = binary.ReadInt64(p.current.operandStack.content, top-16)
	return
}

func (p *Processor) popInt64() int64 {
	return 0
}

func (p *Processor) pushInt64(v int64) {
}
