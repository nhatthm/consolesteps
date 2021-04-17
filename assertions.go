package consolesteps

import (
	"fmt"
	"strings"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

type tError struct {
	err error
}

func (t *tError) Helper() {}

func (t *tError) Errorf(format string, args ...interface{}) {
	t.err = fmt.Errorf(format, args...) // nolint: goerr113
}

func (t *tError) LastError() error {
	return t.err
}

func teeError() *tError {
	return &tError{}
}

// AssertState asserts console state.
func AssertState(t assert.TestingT, terminal vt10x.Terminal, expected string) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	actual := trimTailingSpaces(expect.StripTrailingEmptyLines(terminal.String()))

	return assert.Equal(t, expected, actual)
}

// AssertStateRegex asserts console state.
func AssertStateRegex(t assert.TestingT, terminal vt10x.Terminal, expected string) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	actual := trimTailingSpaces(expect.StripTrailingEmptyLines(terminal.String()))

	return assert.Regexp(t, expected, actual)
}

func trimTailingSpaces(out string) string {
	lines := strings.Split(out, "\n")

	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}

	return strings.Join(lines, "\n")
}
