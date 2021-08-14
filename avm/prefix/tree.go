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

package prefix

import "unsafe"

type Identifier64 uint64

func readIdentifier64(b []byte) (id Identifier64) {
	length := int(b[0])
	if length > 8 {
		panic("identifier is more than 8 bytes")
	}
	if len(b) > 8 {
		id = Identifier64(*(*uint64)(unsafe.Pointer(&b[1])))
	} else {
		padded := make([]byte, 9)
		copy(padded, b)
		id = Identifier64(*(*uint64)(unsafe.Pointer(&padded[1])))
	}
	id <<= (8 - length) << 8
	return
}
