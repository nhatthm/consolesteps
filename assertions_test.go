package consolesteps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestingT_LastError(t *testing.T) {
	t.Parallel()

	tee := teeError()
	tee.Errorf("error: %s", "unknown")

	assert.EqualError(t, tee.LastError(), `error: unknown`)
}
