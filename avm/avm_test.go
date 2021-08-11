package avm_test

import (
	"avm"
	"avm/memory"
	"avm/prefix"
	"github.com/stretchr/testify/assert"
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
			methodArea: memory.NewModule(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
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
			methodArea: memory.NewModule(map[prefix.Identifier64]map[prefix.Identifier64][]byte{
				17: {
					0:  []byte{0x10, 0x01, 0x10, 0x02, 0x11, 0x09},
					10: []byte{},
				},
			}),
			calledApp:  17,
			arguments:  []byte{},
			wantOutput: []byte{0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
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
