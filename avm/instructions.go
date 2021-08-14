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
)

func (p *Processor) noOp() {}

// invokeDispatcher invokes the dispatcher method of another application
//
// Format:
//		invokeDispatcher
// OperandStack:
// 		[..., id64 ->
// 		[... <-
// Description:
//
// `id64` is the 64-bit representation of the applicationID of the called
// application. This identifier value is popped from the stack and nothing
// is pushed onto the stack.
func (p *Processor) invokeDispatcher() {
	appID := p.popIdentifier64()
	p.callMethod(appID, appID, DispatcherID, false)
}

func (p *Processor) indInvokeDispatcher() {
	p.invokeDispatcher()
	p.current.isIndependent = true
	p.heap.Save()
}

func (p *Processor) spawnDispatcher() {
	if p.findIndependentCaller() != 0 {
		panic(InvalidSpawnState)
	}
	appID := p.popIdentifier64()
	p.callMethod(appID, appID, DispatcherID, true)
}

func (p *Processor) invokeInternal() {
	p.callMethod(p.current.context, p.current.methodID.appID, p.popIdentifier64(), false)
}

func (p *Processor) indInvokeInternal() {
	p.invokeInternal()
	p.current.isIndependent = true
	p.heap.Save()
}

func (p *Processor) ret0() {
	p.returnBytes(0, NoError)
}

func (p *Processor) ret64() {
	p.returnBytes(8, NoError)
}

func (p *Processor) throw() {
	top := p.current.operandStack.length()
	n := binary.ReadUint16(p.current.operandStack.content, top-2)
	p.throwBytes(int64(n)+2, SoftwareError)
}

func (p *Processor) enter() {
	// we should use context here
	if lock := p.entranceLocks[p.current.context]; lock != nil && *lock {
		panic(Reentrancy)
	}
	p.current.entranceLock = new(bool)
	*p.current.entranceLock = true
	p.entranceLocks[p.current.context] = p.current.entranceLock
}

func (p *Processor) pushC64() {
	top := p.current.operandStack.length()
	p.current.operandStack.ensureLen(top + 8)
	p.methodArea.Load64(p.current.pc, p.current.operandStack.content, top)
	p.current.pc += 8
}

func (p *Processor) pop64() {
	p.current.operandStack.shrinkTo(p.current.operandStack.length() - 8)
}

func (p *Processor) iAdd() {
	a, b, top := p.peekInt64()
	binary.PutInt64(p.current.operandStack.content, top-16, a+b)
	p.current.operandStack.shrinkTo(top - 8)
}

func (p *Processor) iSub() {
	a, b, top := p.peekInt64()
	binary.PutInt64(p.current.operandStack.content, top-16, b-a)
	p.current.operandStack.shrinkTo(top - 8)
}

func (p *Processor) argC16() {
	offset := int64(p.readConst16())
	top := p.current.operandStack.length()
	binary.Copy64(p.nextLocalFrame.content, offset, p.current.operandStack.content, top-8)
	p.current.operandStack.shrinkTo(top - 8)
}

// lfLoadC16 loads 64 bits from the local frame using a 16-bit unsigned
// constant index
//
// Format:
//		lfLoadC16 2bIndex
// OperandStack:
// 		[... ->
// 		[..., value <-
// Description:
//
// The `Index` is an unsigned 16-bit integer that must be an index into the
// current local frame. Eight bytes from the position `Index` to `Index+7`
// (inclusive) of the local frame is considered as a single `value` and is
// pushed onto the operand stack.
func (p *Processor) lfLoadC16() {
	index := int64(p.readConst16())
	top := p.current.operandStack.length()
	p.current.operandStack.ensureLen(top + 8)
	binary.Copy64(p.current.operandStack.content, top, p.current.localFrame.content, index)
}

func (p *Processor) lfStoreC16() {
	index := int64(p.readConst16())
	top := p.current.operandStack.length()
	p.current.localFrame.ensureLen(index + 8)
	binary.Copy64(p.current.localFrame.content, index, p.current.operandStack.content, top-8)
	p.current.operandStack.shrinkTo(top - 8)
}

func (p *Processor) jmpEqC16() {
	a, b, _ := p.peekInt64()
	if a == b {
		offset := int64(int16(p.readConst16()))
		p.current.pc += offset
	} else {
		p.current.pc += 2
	}
}

/*
// we will only have 64bit offset smaller integers will be used for push

func (p *Processor) hLoadLocal() {
	p.heap.LoadChild(0).LoadChild(readIdentifier64(p.operandStack.content, p.operandStack.top-8))
	p.operandStack.top -= 8
}

func (p *Processor) hUnLoadLocal() {
	p.heap.UnLoadChild().UnLoadChild()
	// check errors
}

func (p *Processor) hLoad64() {
	offset := readInt64(p.operandStack.content, p.operandStack.top-8)
	err := p.heap.LoadBytes8(
		offset,
		p.operandStack.content[p.operandStack.top-8:])
	if err != nil {
		panic("ijkhjknt")
	}
}
*/
