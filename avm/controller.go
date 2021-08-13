package avm

import (
	"go-AVM/avm/memory"
	"go-AVM/avm/prefix"
	"log"
	"reflect"
)

//go:generate awk -f prog.awk controller.go

type Opcode byte

type Controller struct {
	processor           Processor
	instructionRoutines []func()
}

func NewController() (c *Controller) {
	c = &Controller{}
	c.instructionRoutines = []func(){
		0x00: c.processor.noOp,
		0x01: c.processor.invokeDispatcher,
		0x02: c.processor.indInvokeDispatcher,
		0x03: c.processor.spawnDispatcher,
		0x04: c.processor.invokeInternal,
		0x05: c.processor.indInvokeInternal,
		0x06: nil,
		0x07: nil,
		0x08: c.processor.ret0,
		0x09: c.processor.ret64,
		0x0a: nil,
		0x0b: c.processor.throw,
		0x0c: c.processor.enter,
		0x10: c.processor.pushC64,
		0x11: c.processor.iAdd,
		0x12: c.processor.argC16,
		0x13: c.processor.lfLoadC16,
		0x14: c.processor.lfStoreC16,
		0x15: c.processor.jmpEqC16,
	}
	return
}

func (c *Controller) SetupNewSession(calledApp prefix.Identifier64, argumentBuffer []byte,
	methodArea, heap *memory.Module) *Controller {
	c.processor = *newProcessor(&dynamicArray{
		content: argumentBuffer,
		maxSize: MaxLocalFrameSize,
	}, heap, methodArea)
	c.processor.callMethod(calledApp, calledApp, DispatcherID, false)
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
			log.Println(r)
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
