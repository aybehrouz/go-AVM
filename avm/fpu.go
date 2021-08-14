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
	"fmt"
	"math"
)

const MaxAdditionLoss = 0x100000

func add(a, b uint64) float64 {
	e1 := extractExp(a)
	e2 := extractExp(b)

	if e1 > e2 {
		fmt.Printf("%v %v %v %b ", e1, e2, e1-e2, (uint64(math.MaxUint64)>>(64-e1+e2))&b)
		if (math.MaxUint64>>(64-e1+e2))&b > MaxAdditionLoss {
			panic("under")
		}
	} else if e2 < e1 {
		if (math.MaxUint64>>(64-e2+e1))&a > MaxAdditionLoss {
			panic("under")
		}
	}
	return math.Float64frombits(a) + math.Float64frombits(b)
}

func truncate(f uint64, n int) uint64 {
	shift := (52 + 1023) - (int(extractExp(f)) + n)
	fmt.Printf("sh:%v ", shift)
	if shift <= 0 {
		return f
	}
	if shift > 52 {
		return 0
	}
	fmt.Printf("%b", uint(math.MaxUint64<<shift))
	return f & (math.MaxUint64 << shift)
}

func extractExp(f uint64) uint64 {
	return f >> 52 & 0x7ff
}
