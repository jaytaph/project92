package screen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewport_GetScreenCoordinates(t *testing.T) {
	vp, _ := New(0, 0, 100, 100, nil)
	childVp, _ := New(5, 7, 40, 40, vp)
	child2Vp, _ := New(5, 46, 40, 10, vp)

	var (
		x, y int
		err  error
	)

	x, y, err = childVp.GetScreenCoordinates(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 5, x)
	assert.Equal(t, 7, y)

	x, y, err = childVp.GetScreenCoordinates(10, 10)
	assert.NoError(t, err)
	assert.Equal(t, 15, x)
	assert.Equal(t, 17, y)

	x, y, err = child2Vp.GetScreenCoordinates(10, 10)
	assert.NoError(t, err)
	assert.Equal(t, 15, x)
	assert.Equal(t, 56, y)

	x, y, err = vp.GetScreenCoordinates(15, 15)
	assert.NoError(t, err)
	assert.Equal(t, 15, x)
	assert.Equal(t, 15, y)
}
