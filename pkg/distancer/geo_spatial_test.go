package distancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeoSpatialDistance(t *testing.T) {
	t.Run("between Munich and Stuttgart", func(t *testing.T) {
		munich := []float32{48.137154, 11.576124}
		stuttgart := []float32{48.783333, 9.183333}

		dist, err := NewGeoProvider().New(munich).Distance(stuttgart)
		require.Nil(t, err)
		assert.InDelta(t, 190000, dist, 1000)
	})
}
