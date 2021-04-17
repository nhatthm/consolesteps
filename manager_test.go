package consolesteps_test

import (
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/cucumber/godog"

	"go.nhat.io/consolesteps"
)

func TestManager(t *testing.T) {
	t.Parallel()

	m := consolesteps.New(t,
		consolesteps.WithTermSize(80, 24),
		consolesteps.WithStarter(func(sc *godog.Scenario, console *expect.Console) {
			console.Write([]byte(`hello world`)) // nolint: errcheck, gosec
		}),
	)

	scenario := &godog.Scenario{}
	_, terminal := m.NewConsole(scenario)

	// New again does not affect the state.
	_, _ = m.NewConsole(scenario)

	m.Flush()

	expected := `hello world`

	consolesteps.AssertState(t, terminal, expected)

	m.CloseConsole(scenario)

	// Close again does not get error.
	m.CloseConsole(scenario)
}
