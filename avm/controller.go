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
	"go-AVM/avm/memory"
	"go-AVM/avm/prefix"
	"log"
	"reflect"
)

type Opcode byte

type Controller struct {
	processor           Processor
	instructionRoutines []func()
}

// run this after changing instructionRoutines array
//go:generate /bin/sh awk.sh

func NewController() (c *Controller) {
	c = &Controller{}
	c.instructionRoutines = []func(){
		0x00: c.processor.noOp,
		0x01: c.processor.invokeDispatcher,
		0x02: c.processor.indInvokeDispatcher,
		0x03: c.processor.spawnDispatcher,
		0x04: c.processor.invokeInternal,
		0x05: c.processor.indInvokeInternal,
		0x08: c.processor.ret0,
		0x09: c.processor.ret64,
		0x0b: c.processor.throw,
		0x0c: c.processor.enter,
		0x10: c.processor.pushC64,
		0x11: c.processor.pop,
		0x12: c.processor.iAdd,
		0x13: c.processor.iSub,
		0x14: c.processor.argC16,
		0x15: c.processor.lfLoadC16,
		0x16: c.processor.lfStoreC16,
		0x17: c.processor.jmpEqC16,
	}
	return
}

func (c *Controller) SetupNewSession(calledApp prefix.Identifier64, argumentBuffer []byte,
	methodArea, heap *memory.Module) *Controller {
	c.processor = *newProcessor(&dynamicArray{
		content: argumentBuffer,
		maxSize: MaxLocalFrameSize,
	}, heap, methodArea)
	c.processor.callMethod(calledApp, calledApp, DispatcherID)
	c.processor.current.isIndependent = true
	c.processor.heap.Save()
	return c
}

func (c *Controller) Emulate() ([]byte, ErrorCode) {
	eof := false
	for !eof {
		eof = c.EmulateNextInstruction()
	}
	return c.processor.returnData, c.processor.errorStatus
}

func (c *Controller) EmulateNextInstruction() (eof bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("panic:", r)
			c.processor.throwBytes(0, convertToErrorCode(r))
			eof = false
		}
	}()
	opcode, eof := c.processor.nextOpcode()
	if eof {
		return true
	}
	c.instructionRoutines[opcode]()
	return false
}

func convertToErrorCode(r interface{}) ErrorCode {
	switch reflect.TypeOf(r).String() {
	case "runtime.boundsError":
		return InvalidReference
	case "avm.ErrorCode":
		return r.(ErrorCode)
	default:
		return RuntimeError
	}
}
