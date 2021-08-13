package avm

import (
	"go-AVM/avm/binary"
)

func (p *Processor) noOp() {}

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

func (p *Processor) lfLoadC16() {
	offset := int64(p.readConst16())
	top := p.current.operandStack.length()
	p.current.operandStack.ensureLen(top + 8)
	binary.Copy64(p.current.operandStack.content, top, p.current.localFrame.content, offset)
}

func (p *Processor) lfStoreC16() {
	offset := int64(p.readConst16())
	top := p.current.operandStack.length()
	binary.Copy64(p.current.localFrame.content, offset, p.current.operandStack.content, top-8)
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
