package avm

import (
	"go-AVM/avm/binary"
	"go-AVM/avm/prefix"
)

func (p *Processor) noOp() {}

func (p *Processor) invokeDispatcher() {
	appID := p.popIdentifier64()
	p.callMethod(appID, appID, DispatcherID)
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
	p.callStackQueue = append(p.callStackQueue, []*CallInfo{&CallInfo{
		pc:      0,
		context: appID,
		methodID: struct {
			appID   prefix.Identifier64
			localID prefix.Identifier64
		}{appID, DispatcherID},
		isIndependent: false,
		operandStack:  newOperandStack(),
		localFrame:    p.nextLocalFrame,
	}})
	p.nextLocalFrame = newLocalFrame()
}

func (p *Processor) invokeInternal() {
	p.callMethod(p.current.context, p.current.methodID.appID, p.popIdentifier64())
}

func (p *Processor) indInvokeInternal() {
	p.invokeInternal()
	p.current.isIndependent = true
	p.heap.Save()
}

func (p *Processor) ret() {
	p.returnBytes(0, NoError)
}

func (p *Processor) ret64() {
	p.returnBytes(8, NoError)
}

func (p *Processor) throw() {
	top := p.current.operandStack.length()
	n := p.current.operandStack.content[top-1]
	p.throwBytes(int64(n)+1, SoftwareError)
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

func (p *Processor) iAdd64() {
	top := p.current.operandStack.length()
	a := binary.ReadInt64(p.current.operandStack.content, top-8)
	b := binary.ReadInt64(p.current.operandStack.content, top-16)
	binary.PutInt64(p.current.operandStack.content, top-16, a+b)
	p.current.operandStack.shrinkTo(top - 8)
}

func (p *Processor) popIdentifier64() prefix.Identifier64 {
	top := p.current.operandStack.length()
	id := binary.ReadIdentifier64(p.current.operandStack.content, top-8)
	p.current.operandStack.shrinkTo(top - 8)
	return id
}

func (p *Processor) popInt64() int64 {
	return 0
}

func (p *Processor) pushInt64(v int64) {
}

/*
// we will only have 64bit offset smaller integers will be used for push
func (p *Processor) lFrameLoad64() {
	top := p.operandStack.length()
	offset := binary.ReadInt64(p.operandStack.content, top-8)
	binary.CopyBytes(p.operandStack.content, top-8, p.localFrame.content, offset, 8)
}

func (p *Processor) push64() {
	top := p.operandStack.length()
	p.operandStack.ensureLen(top + 8)
	p.methodArea.LoadBytes(pc)

}

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


func (p *Processor) invokeInternal() {
	// callInfo {opstack, localFram, pc, methodID}
	p.operandStack = new
	p.localFrame = new
}
*/
