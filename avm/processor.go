package avm

import (
	"avm/binary"
	"avm/memory"
	"avm/prefix"
)

const DispatcherID = 0

const (
	InitialOpStackSize    = 512
	MaxOpStackSize        = 128 * 1024
	InitialLocalFrameSize = 4 * 1024
	MaxLocalFrameSize     = 256 * 1024
)

type ErrorCode int

const (
	NoError ErrorCode = iota
	InvalidOperands
	InvalidSpawnState
	SoftwareError
	OutOfBounds
	OutOfMemory
	OverFlow
	UnderFlow
	PrecisionLoss
	Reentrancy
)

type CallInfo struct {
	pc       int64
	context  prefix.Identifier64
	methodID struct {
		appID   prefix.Identifier64
		localID prefix.Identifier64
	}
	isIndependent bool
	entranceLock  *bool
	operandStack  *dynamicArray
	localFrame    *dynamicArray
}

type Processor struct {
	callStackQueue [][]*CallInfo
	current        *CallInfo
	errorStatus    ErrorCode
	entranceLocks  map[prefix.Identifier64]*bool
	nextLocalFrame *dynamicArray
	returnData     []byte
	heap           *memory.Module
	methodArea     *memory.Module
}

func newProcessor(nextLocalFrame *dynamicArray, heap, methodArea *memory.Module) *Processor {
	return &Processor{
		callStackQueue: [][]*CallInfo{{}},
		errorStatus:    NoError,
		entranceLocks:  map[prefix.Identifier64]*bool{},
		nextLocalFrame: nextLocalFrame,
		heap:           heap,
		methodArea:     methodArea,
	}
}

func (p *Processor) callMethod(context, app, method prefix.Identifier64) {
	p.callStackQueue[0] = append(p.callStackQueue[0], &CallInfo{
		context: context,
		methodID: struct {
			appID   prefix.Identifier64
			localID prefix.Identifier64
		}{app, method},
		operandStack: newOperandStack(),
		localFrame:   p.nextLocalFrame,
	})
	p.updateCurrentCallContext()
	p.errorStatus = NoError
}

func (p *Processor) updateCurrentCallContext() {
	if p.callStackQueue == nil {
		p.current = nil
		p.nextLocalFrame = nil
		return
	}
	p.current = p.callStackQueue[0][len(p.callStackQueue[0])-1]
	p.nextLocalFrame = newLocalFrame()
	p.heap.LoadRoot(p.current.context)
	p.methodArea.LoadRoot(p.current.methodID.appID).LoadChild(p.current.methodID.localID)
}

func (p *Processor) returnBytes(n int64, status ErrorCode) {
	// important: this function may panic when n > 0, but it should not panic when n == 0. When it panics
	// it must not change any state or have any effects
	if top := len(p.callStackQueue[0]); top <= 1 {
		if n > 0 {
			l := p.current.operandStack.length()
			p.returnData = append(p.returnData, p.current.operandStack.content[l-n:l]...)
		}
		// update the call stack queue
		if len(p.callStackQueue) <= 1 || status != NoError {
			p.callStackQueue = nil
		} else {
			p.callStackQueue[0] = nil
			p.callStackQueue = p.callStackQueue[1:]
		}
	} else {
		nextCallInfo := p.callStackQueue[0][top-2]
		if n > 0 {
			// `ensureLen()` changes the state of the caller's stack, so before calling it, first we should make
			// sure that binary.CopyBytes won't panic
			if p.current.operandStack.length() < n {
				panic(InvalidOperands)
			}
			callerStackTop := nextCallInfo.operandStack.length()
			nextCallInfo.operandStack.ensureLen(callerStackTop + n)
			binary.CopyBytes(
				nextCallInfo.operandStack.content, callerStackTop,
				p.current.operandStack.content, p.current.operandStack.length()-n, n)
		}
		// update the call stack queue
		p.callStackQueue[0][top-1] = nil
		p.callStackQueue[0] = p.callStackQueue[0][:top-1]
	}
	if status != NoError {
		p.heap.Restore()
	} else if p.current.isIndependent || p.callStackQueue == nil {
		p.heap.Discard()
	}
	if p.current.entranceLock != nil {
		*p.current.entranceLock = false
	}
	p.errorStatus = status
	p.updateCurrentCallContext()
}

func (p *Processor) throwBytes(n int64, code ErrorCode) {
	ic := p.findIndependentCaller()
	for i := ic + 1; i < len(p.callStackQueue[0]); i++ {
		*p.callStackQueue[0][i].entranceLock = false
		p.callStackQueue[0][i] = nil
	}
	p.callStackQueue[0] = p.callStackQueue[0][:ic+1]
	p.returnBytes(n, code)
}

func (p *Processor) findIndependentCaller() int {
	for i := len(p.callStackQueue[0]) - 1; i > 0; i-- {
		if p.callStackQueue[0][i].isIndependent {
			return i
		}
	}
	return 0
}

func (p *Processor) nextOpcode() (Opcode, bool) {
	if p.current == nil {
		return 0, true
	}
	opcode := Opcode(p.methodArea.LoadByte(p.current.pc))
	p.current.pc++
	return opcode, false
}

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
		panic("max size exceeded")
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
