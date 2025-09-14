package coordinator

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SSripilaipong/muon/common/chn"
)

func TestSubmit_guaranteeAtLeastHalfSuccess(t *testing.T) {
	t.Run("should return success for closed channel", func(t *testing.T) {
		ch := make(chan error)
		close(ch)
		ok, _ := guaranteeAtLeastHalfSuccess(context.Background(), 1, ch)
		assert.True(t, ok)
	})

	t.Run("should return success when error is not more than a half", func(t *testing.T) {
		ch := chn.FromSlice([]error{errors.New(""), errors.New("")})
		ok, _ := guaranteeAtLeastHalfSuccess(context.Background(), 5, ch)
		assert.True(t, ok)
	})

	t.Run("should return fail when error is more than a half", func(t *testing.T) {
		ch := chn.FromSlice([]error{errors.New(""), errors.New(""), errors.New("")})
		ok, _ := guaranteeAtLeastHalfSuccess(context.Background(), 5, ch)
		assert.False(t, ok)
	})

	t.Run("should return fail when error is a half", func(t *testing.T) {
		ch := chn.FromSlice([]error{errors.New(""), errors.New("")})
		ok, _ := guaranteeAtLeastHalfSuccess(context.Background(), 4, ch)
		assert.False(t, ok)
	})
}
