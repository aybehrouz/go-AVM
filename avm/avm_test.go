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
		name       string
		methodArea *memory.Module
		heap       *memory.Module
		calledApp  prefix.Identifier64
		arguments  []byte
		wantOutput []byte
		wantError  avm.ErrorCode
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
			name: "nonexistent App",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0: assembler.AssembleString("pushC64 443 invokeDispatcher ret0"),
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
				testCase.heap = testCase.methodArea
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
		})
	}
}

func BenchmarkLoad(b *testing.B) {
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

	b.Run("as interface", func(b *testing.B) {
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
