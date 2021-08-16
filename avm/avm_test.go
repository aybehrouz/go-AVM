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

package avm_test

import (
	"github.com/stretchr/testify/assert"
	"go-AVM/assembler"
	"go-AVM/avm"
	"go-AVM/avm/memory"
	"go-AVM/avm/prefix"
	"testing"
)

func TestController_Emulate(t *testing.T) {
	tests := []struct {
		name              string
		methodArea        *memory.Module
		heap              *memory.Module
		calledApp         prefix.Identifier64
		arguments         []byte
		wantOutput        []byte
		wantMethodAreaLog string
		wantHeapLog       string
		wantError         avm.ErrorCode
	}{
		{
			name: "empty",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0:  []byte{},
					10: []byte{},
				},
			}),
			calledApp:  17,
			wantOutput: nil,
			wantError:  avm.InvalidReference,
		},
		{
			name: "simple return",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0:  assembler.AssembleString("pushC64 2 pushC64 3 iAdd ret64"),
					10: []byte{},
				},
			}),
			calledApp:  17,
			wantOutput: []byte{0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantError:  avm.NoError,
		},
		{
			name: "nonexistent App",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				0x11: {
					0: assembler.AssembleString("pushC64 443 invokeDispatcher ret0"),
				},
			}),
			calledApp:   0x11,
			wantOutput:  nil,
			wantHeapLog: "root<-11   Save   root<-1bb   Restore",
			wantError:   avm.InvalidReference,
		},
		{
			name: "catch nonexistent App error",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				0x11: {
					0: assembler.AssembleString("pushC64 443 indInvokeDispatcher pushC64 7 ret64"),
				},
			}),
			calledApp:   0x11,
			wantOutput:  []byte{0x7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantHeapLog: "root<-11   Save   root<-1bb   Save   Restore   root<-11   Discard",
			wantError:   avm.NoError,
		},
		{
			name: "multiple internal return",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0: assembler.AssembleString("pushC64 1 invokeInternal pushC64 1 iAdd ret64"),
					1: assembler.AssembleString("pushC64 2 invokeInternal pushC64 2 iAdd ret64"),
					2: assembler.AssembleString("pushC64 3 invokeInternal pushC64 3 iAdd ret64"),
					3: assembler.AssembleString("pushC64 10 ret64"),
				},
			}),
			calledApp:  17,
			wantOutput: []byte{16, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantError:  avm.NoError,
		},
		{
			name: "parameter passing",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0: assembler.AssembleString("pushC64 777777 argC16 2d0 pushC64 5 invokeInternal ret64"),
					5: assembler.AssembleString("lfLoadC16 2d0 pushC64 10 iAdd ret64"),
				},
			}),
			calledApp:  17,
			wantOutput: []byte{0x3b, 0xde, 0xb, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantError:  avm.NoError,
		},

		// Throw instruction tests:
		{
			name: "simple throw",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0: assembler.AssembleString("pushC64 566265685016576 throw"),
				},
			}),
			calledApp:  17,
			wantOutput: []byte{0x4, 0x3, 0x2, 0x0},
			wantError:  avm.SoftwareError,
		},
		{
			name: "failed throw",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				0x11: {
					0:    assembler.AssembleString("pushC64 5 pushC64 0x12 invokeInternal ret64"),
					0x12: assembler.AssembleString("pushC64 0x0700000000000000 throw"),
				},
			}),
			calledApp:   0x11,
			wantOutput:  nil,
			wantHeapLog: "root<-11   Save   root<-11   Restore",
			wantError:   avm.InvalidReference,
		},
		{
			name: "catch failed throw",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				0x11: {
					0:    assembler.AssembleString("pushC64 5 pushC64 0x12 indInvokeInternal ret64"),
					0x12: assembler.AssembleString("pushC64 0x0700000000000000 throw"),
				},
			}),
			calledApp:   0x11,
			wantOutput:  []byte{0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantHeapLog: "root<-11   Save   root<-11   Save   Restore   root<-11   Discard",
			wantError:   avm.NoError,
		},
		{
			name: "multiple throw",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0: assembler.AssembleString("pushC64 1 invokeInternal pushC64 1 iAdd ret64"),
					1: assembler.AssembleString("pushC64 2 invokeInternal pushC64 2 iAdd ret64"),
					2: assembler.AssembleString("pushC64 3 invokeInternal pushC64 3 iAdd ret64"),
					3: assembler.AssembleString("pushC64 566265685016576 throw iAdd"),
				},
			}),
			calledApp:  17,
			wantOutput: []byte{0x4, 0x3, 0x2, 0x0},
			wantError:  avm.SoftwareError,
		},
		{
			name: "multiple external throw",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				0x11: {
					0: assembler.AssembleString("pushC64 6 pushC64 0x12 invokeDispatcher iAdd ret64"),
				},
				0x12: {
					0: assembler.AssembleString("pushC64 0x13 invokeDispatcher pushC64 20 iAdd ret64"),
				},
				0x13: {
					0: assembler.AssembleString("pushC64 0x14 invokeDispatcher pushC64 30 iAdd ret64"),
				},
				0x14: {
					0: assembler.AssembleString("pushC64 0x0006000000000001 throw ret0"),
				},
			}),
			calledApp:   0x11,
			wantOutput:  []byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6, 0x0},
			wantHeapLog: "root<-11   Save   root<-12   root<-13   root<-14   Restore",
			wantError:   avm.SoftwareError,
		},
		{
			name: "multi throw and catch",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				0x11: {
					0: assembler.AssembleString("pushC64 6 pushC64 0x12 indInvokeDispatcher iAdd ret64"),
				},
				0x12: {
					0: assembler.AssembleString("pushC64 0x13 invokeDispatcher pushC64 20 iAdd ret64"),
				},
				0x13: {
					0: assembler.AssembleString("pushC64 0x14 invokeDispatcher pushC64 30 iAdd ret64"),
				},
				0x14: {
					0: assembler.AssembleString("pushC64 0x0006000000000001 throw ret0"),
				},
			}),
			calledApp:   0x11,
			wantOutput:  []byte{0x7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6, 0x0},
			wantHeapLog: "root<-11   Save   root<-12   Save   root<-13   root<-14   Restore   root<-11   Discard",
			wantError:   avm.NoError,
		},

		// Entrance lock tests:
		{
			name: "simple lock",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				0x11: {
					0: assembler.AssembleString("pushC64 6 pushC64 0x2 invokeInternal iAdd ret64"),
					2: assembler.AssembleString("enter pushC64 4 pushC64 0x12 invokeDispatcher ret64"),
				},
				0x12: {
					0: assembler.AssembleString("pushC64 0x11 invokeDispatcher ret64"),
				},
			}),
			calledApp:   0x11,
			wantOutput:  nil,
			wantHeapLog: "root<-11   Save   root<-11   root<-12   root<-11   root<-11   Restore",
			wantError:   avm.Reentrancy,
		},
		{
			name: "simple lock catch",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				0x11: {
					0: assembler.AssembleString("pushC64 6 pushC64 0x2 invokeInternal iAdd ret64"),
					2: assembler.AssembleString("enter pushC64 4 pushC64 0x12 indInvokeDispatcher ret64"),
				},
				0x12: {
					0: assembler.AssembleString("pushC64 0x11 invokeDispatcher ret64"),
				},
			}),
			calledApp:  0x11,
			wantOutput: []byte{10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantHeapLog: "root<-11   Save   root<-11   root<-12   Save   root<-11   root<-11   Restore   " +
				"root<-11   root<-11   Discard",
			wantError: avm.NoError,
		},
		{
			name: "simple lock opening",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				0x11: {
					0: assembler.AssembleString("pushC64 0x2 invokeInternal lfLoadC16 2d0 jmpEqC16 2d10" +
						" pushC64 0x12 invokeDispatcher iAdd ret64"),
					2: assembler.AssembleString("enter pushC64 4 ret64"),
				},
				0x12: {
					0: assembler.AssembleString("pushC64 4 argC16 2d0 pushC64 0x11 invokeDispatcher ret64"),
				},
			}),
			calledApp:  0x11,
			arguments:  []byte{6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantOutput: []byte{14, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantError:  avm.NoError,
		},

		// Full programs:
		{
			name: "sum 1:1",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0: assembler.AssembleString("lfLoadC16 2d0 pushC64 -1 jmpEqC16 2d19 iAdd " +
						"argC16 2d0 pushC64 0 invokeInternal lfLoadC16 2d0 iAdd ret64 pushC64 0 ret64"),
				},
			}),
			calledApp:  17,
			arguments:  []byte{1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantOutput: []byte{1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantError:  avm.NoError,
		},
		{
			name: "sum 1:200",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0: assembler.AssembleString("lfLoadC16 2d0 pushC64 -1 jmpEqC16 2d19 iAdd " +
						"argC16 2d0 pushC64 0 invokeInternal lfLoadC16 2d0 iAdd ret64 pushC64 0 ret64"),
				},
			}),
			calledApp:  17,
			arguments:  []byte{200, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantOutput: []byte{0x84, 0x4e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantError:  avm.NoError,
		},
	}

	controller := avm.NewController()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.heap == nil {
				testCase.heap = memory.NewMocker(nil)
			}
			controller.SetupNewSession(
				testCase.calledApp,
				testCase.arguments,
				testCase.methodArea,
				testCase.heap,
			)
			gotOutput, gotError := controller.Emulate()
			assert.Equal(t, testCase.wantOutput, gotOutput, "invalid output")
			assert.Equal(t, testCase.wantError, gotError, "invalid error code")
			if testCase.wantHeapLog != "" {
				assert.Equal(t, testCase.wantHeapLog, testCase.heap.AccessLog(), "invalid heap access log")
			}
		})
	}
}

func BenchmarkFib(b *testing.B) {
	n := 8
	controller := avm.NewController()
	methodArea := memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
		17: {
			0: assembler.AssembleString(""),
		},
	})
	arguments := []byte{byte(n), 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		controller.SetupNewSession(17, arguments, methodArea, methodArea)
		controller.Emulate()
	}
}

func BenchmarkMemoryLoad(b *testing.B) {
	var mInterface tempInterface
	m := memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
		17: {
			0:  []byte{0x08, 0x02, 0x06, 0x08, 0x02, 0x06, 0x08, 0x02, 0x06, 0x08, 0x02, 0x06, 0x08, 0x02, 0x06},
			10: []byte{},
		},
	})
	m.LoadRoot(17).LoadChild(0)
	mInterface = m
	temp := make([]byte, 5*1024)

	b.Run("normal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < len(temp)-10; i++ {
				m.Load64(6, temp, int64(i))
			}
		}
	})

	b.Run("with interface", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < len(temp)-10; i++ {
				mInterface.Load64(6, temp, int64(i))
			}
		}
	})
}

type tempInterface interface {
	Load64(int64, []byte, int64)
}
