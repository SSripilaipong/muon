package prl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SSripilaipong/muon/common/chn"
)

func TestCollect(t *testing.T) {
	t.Run("should return all values", func(t *testing.T) {
		assert.ElementsMatch(t, []int{1, 2, 3}, chn.All(Collect(context.Background(),
			func() int { return 2 },
			func() int { return 1 },
			func() int { return 3 },
		)))
	})
}
