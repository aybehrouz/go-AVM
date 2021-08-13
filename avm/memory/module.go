package memory

import (
	"go-AVM/avm/binary"
	. "go-AVM/avm/prefix"
)

// Module error handling will be done by panicking instead of returning errors
type Module struct {
	chunks  map[Identifier64]map[Identifier64][]byte
	root    map[Identifier64][]byte
	current []byte
}

func (m *Module) LoadRoot(id Identifier64) *Module {
	// println("root changed:-> ", id)
	m.root = m.chunks[id]
	m.current = nil
	return m
}

func (m *Module) LoadChild(id Identifier64) *Module {
	// println("child changed:-> ", id)
	m.current = m.root[id]
	return nil
}

func (m *Module) UnLoadChild() *Module {
	return nil
}

func (m *Module) Load64(chunkOffset int64, dst []byte, dstOffset int64) {
	binary.Copy64(dst, dstOffset, m.current, chunkOffset)
}

func (m *Module) LoadUint16(offset int64) uint16 {
	return binary.ReadUint16(m.current, offset)
}

func (m *Module) StoreBytes8(offset int64, src []byte) {
}

func (m *Module) LoadByte(offset int64) byte {
	return m.current[offset]
}

func (m *Module) LoadInt64(offset int64) int64 {
	return 0
}

func (m *Module) StoreBytes(offset int64, num int, src []byte) {
}

func (m *Module) Restore() {
	println("Restore")

}

func (m *Module) Discard() {
	println("Discard")

}

func (m *Module) Save() {
	println("Save")
}

func NewMocker(chunks map[Identifier64]map[Identifier64][]byte) *Module {
	return &Module{chunks: chunks}
}
