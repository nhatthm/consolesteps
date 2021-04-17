package consolesteps

import (
	"context"
	"sync"
	"time"

	"github.com/Netflix/go-expect"
	pseudotty "github.com/creack/pty"
	"github.com/cucumber/godog"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/require"
)

// Starter is a callback when console starts.
type Starter func(sc *godog.Scenario, console *expect.Console)

// Closer is a callback when console closes.
type Closer func(sc *godog.Scenario)

// Option configures Manager.
type Option func(m *Manager)

type session struct {
	console  *expect.Console
	terminal vt10x.Terminal
	output   *Buffer
}

// Manager manages console and its state.
type Manager struct {
	sessions map[string]*session
	current  string

	// Terminal size.
	termCols int
	termRows int

	starters []Starter
	closers  []Closer

	test TestingT

	mu sync.Mutex
}

type tHelper interface {
	Helper()
}

// RegisterContext register console Manager to test context.
func (m *Manager) RegisterContext(ctx *godog.ScenarioContext) {
	ctx.Before(func(_ context.Context, sc *godog.Scenario) (context.Context, error) {
		m.NewConsole(sc)

		return nil, nil
	})

	ctx.After(func(_ context.Context, sc *godog.Scenario, _ error) (context.Context, error) {
		m.CloseConsole(sc)

		return nil, nil
	})

	ctx.Step(`console output is:`, m.isConsoleOutput)
	ctx.Step(`console output matches:`, m.matchConsoleOutput)
}

func (m *Manager) session() *session {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.sessions[m.current]
}

// NewConsole creates a new console.
func (m *Manager) NewConsole(sc *godog.Scenario) (*expect.Console, vt10x.Terminal) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sess := &session{}

	if s, ok := m.sessions[sc.Id]; ok {
		return s.console, s.terminal
	}

	m.test.Logf("Console: %s (#%s)\n", sc.Name, sc.Id)

	sess.output = new(Buffer)
	sess.console, sess.terminal = newVT10XConsole(m.test, m.termCols, m.termCols, expect.WithStdout(sess.output))

	m.sessions[sc.Id] = sess
	m.current = sc.Id

	for _, fn := range m.starters {
		fn(sc, sess.console)
	}

	return sess.console, sess.terminal
}

// CloseConsole closes the current console.
func (m *Manager) CloseConsole(sc *godog.Scenario) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sess, ok := m.sessions[sc.Id]
	if !ok {
		return
	}

	m.flushSession(sess)

	for _, fn := range m.closers {
		fn(sc)
	}

	m.test.Logf("Raw output: %q\n", sess.output.String())
	// Dump the terminal's screen.
	m.test.Logf("State: \n%s\n", expect.StripTrailingEmptyLines(sess.terminal.String()))

	delete(m.sessions, sc.Id)
	m.current = ""
}

func (m *Manager) flushSession(s *session) {
	s.console.Expect(expect.EOF, expect.PTSClosed, expect.WithTimeout(10*time.Millisecond)) // nolint: errcheck, gosec
}

// Flush flushes console state.
func (m *Manager) Flush() {
	m.flushSession(m.session())
}

func (m *Manager) isConsoleOutput(expected *godog.DocString) error {
	m.Flush()

	t := teeError()
	AssertState(t, m.session().terminal, expected.Content)

	return t.LastError()
}

func (m *Manager) matchConsoleOutput(expected *godog.DocString) error {
	m.Flush()

	t := teeError()
	AssertStateRegex(t, m.session().terminal, expected.Content)

	return t.LastError()
}

// WithStarter adds a Starter to Manager.
func (m *Manager) WithStarter(s Starter) *Manager {
	m.starters = append(m.starters, s)

	return m
}

// WithCloser adds a Closer to Manager.
func (m *Manager) WithCloser(c Closer) *Manager {
	m.closers = append(m.closers, c)

	return m
}

// New initiates a new console Manager.
func New(t TestingT, options ...Option) *Manager {
	m := &Manager{
		sessions: make(map[string]*session),
		termCols: 80,
		termRows: 100,
		test:     t,
	}

	for _, o := range options {
		o(m)
	}

	return m
}

// WithStarter adds a Starter to Manager.
func WithStarter(s Starter) Option {
	return func(m *Manager) {
		m.WithStarter(s)
	}
}

// WithCloser adds a Closer to Manager.
func WithCloser(c Closer) Option {
	return func(m *Manager) {
		m.WithCloser(c)
	}
}

// WithTermSize sets terminal size cols x rows. Default is 80 x 100.
func WithTermSize(cols, rows int) Option {
	return func(m *Manager) {
		m.termCols = cols
		m.termRows = rows
	}
}

func newVT10XConsole(t TestingT, cols, rows int, opts ...expect.ConsoleOpt) (*expect.Console, vt10x.Terminal) {
	pty, tty, err := pseudotty.Open()
	require.NoError(t, err)

	term := vt10x.New(vt10x.WithWriter(tty))

	term.Resize(cols, rows)

	c, err := expect.NewConsole(append(opts, expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))...)
	require.NoError(t, err)

	return c, term
}
