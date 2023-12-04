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

	line, err := cr.read(3)
	require.NoError(t, err)

	assert.Len(t, line, 3)
}
