package coordinator

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SSripilaipong/muon/common/chn"
)

func TestCommit_guaranteeQuorumSuccess(t *testing.T) {
	t.Run("should return success for closed channel", func(t *testing.T) {
		ch := make(chan error)
		close(ch)
		ok, _ := guaranteeQuorumSuccess(context.Background(), 1, ch)
		assert.True(t, ok)
	})

	t.Run("should return success when error is less than a quorum", func(t *testing.T) {
		ch := chn.FromSlice([]error{errors.New(""), errors.New("")})
		ok, _ := guaranteeQuorumSuccess(context.Background(), 5, ch)
		assert.True(t, ok)
	})

	t.Run("should return fail when error is at least a quorum", func(t *testing.T) {
		ch := chn.FromSlice([]error{errors.New(""), errors.New(""), errors.New("")})
		ok, _ := guaranteeQuorumSuccess(context.Background(), 5, ch)
		assert.False(t, ok)
	})
}
