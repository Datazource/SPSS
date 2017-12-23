package spass

import (
	"github.com/Azure/go-ansiterm/winterm"
)

var (
	termState  uint32
	lTermState uint32
)

func read() ([]byte, error) {

}

// https://docs.microsoft.com/en-us/windows/console/high-level-console-modes
func disableEcho() error {
	var err error
	termState, err = winterm.GetConsole(os.Stdin.Fd())
	if err != nil {
		return err
	}
	lTermState = termState

	termState &^= winterm.ENABLE_ECHO_INPUT
	termState |= winterm.ENABLE_PROCESSED_INPUT | winterm.ENABLE_LINE_INPUT

	return winterm.SetConsoleMode(
		os.Stdin.Fd(), termState,
	)
}

func enableEcho() error {
	return winterm.SetConsoleMode(
		os.Stdin.Fd(), lTermState,
	)
}
