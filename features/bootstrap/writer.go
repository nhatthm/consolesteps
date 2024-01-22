package bootstrap

import (
	"github.com/Netflix/go-expect"
	"github.com/cucumber/godog"
)

type writer struct {
	console *expect.Console
}

func (w *writer) RegisterSteps(s *godog.ScenarioContext) {
	s.Step(`write to console:`, func(s *godog.DocString) error {
		_, err := w.console.Write([]byte(s.Content))

		return err
	})
}

func (w *writer) start(_ *godog.Scenario, console *expect.Console) {
	w.console = console
}

func (w *writer) close(_ *godog.Scenario) {
	w.console = nil
}
