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

func (m *Module) Load64(loadIndex int64, dst []byte, writeIndex int64) {
	binary.Copy64(dst, writeIndex, m.current, loadIndex)
}

func (m *Module) LoadUint16(index int64) uint16 {
	return binary.ReadUint16(m.current, index)
}

func (m *Module) StoreBytes8(offset int64, src []byte) {
}

func (m *Module) LoadByte(index int64) byte {
	return m.current[index]
}

func (m *Module) LoadInt64(offset int64) int64 {
	return 0
}

func (m *Module) StoreBytes(offset int64, num int, src []byte) {
}

func (m *Module) Restore() {
}

func (m *Module) Discard() {
}

func (m *Module) Save() {
}

func NewMocker(chunks map[Identifier64]map[Identifier64][]byte) *Module {
	return &Module{chunks: chunks}
}
