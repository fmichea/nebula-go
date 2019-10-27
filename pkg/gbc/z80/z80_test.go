package z80

import (
	"nebula-go/mocks/pkg/gbc/memorymocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewCPU(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMMU := memorymocks.NewMockMMU(mockCtrl)

	cpu := NewCPU(mockMMU)
	assert.Len(t, cpu.Opcodes, 0x100)
}
