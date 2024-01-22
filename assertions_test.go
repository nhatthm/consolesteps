package consolesteps

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestingT_LastError(t *testing.T) {
	t.Parallel()

	tee := teeError()
	tee.Errorf("error: %s", "unknown")

	require.EqualError(t, tee.LastError(), `error: unknown`)
}
