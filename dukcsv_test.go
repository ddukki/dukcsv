package dukcsv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//

func TestCSVReader(t *testing.T) {
	loc := "testdata/test.csv"
	cr, err := NewReader(loc, true)
	require.NoError(t, err)

	_, err = cr.read(4)
	assert.Error(t, err)

	nCols := len(cr.header)
	nRows := cr.Count()
	for r := int64(0); r < nRows; r++ {
		line, err := cr.read(r)
		require.NoError(t, err)

		assert.Len(t, line, nCols)

		for _, cell := range line {
			assert.NotEqual(t, '"', cell[0])
			assert.NotEqual(t, '"', cell[len(cell)-1])
		}
	}
}
