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
			name: "simple return",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0:  []byte{0x08},
					10: []byte{},
				},
			}),
			calledApp:  17,
			arguments:  []byte{},
			wantOutput: nil,
			wantError:  avm.NoError,
		},
		{
			name: "simple add",
			methodArea: memory.NewMocker(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0:  assembler.AssembleString("pushC64 2 pushC64 3 iAdd64 ret64"),
					10: []byte{},
				},
			}),
			calledApp:  17,
			arguments:  []byte{},
			wantOutput: []byte{0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
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
