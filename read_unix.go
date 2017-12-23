package spass

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

var (
	// state of terminal
	termState *unix.Termios
	// last state of terminal
	lTermState *unix.Termios
)

func read() (data []byte, err error) {
	if err := disableEchoing(); err == nil {
		defer enableEchoing()

		ch := captureSignals()
		defer stopSignals(ch)

		stop := false
		go func() {
			<-ch
			stop = true
		}()

		data = make([]byte, 0)
		bt := make([]byte, 1)
		for {
			if stop {
				break
			}

			_, err = os.Stdin.Read(bt)
			if err != nil {
				break
			}

			switch bt[0] {
			case '\n':
				return data, err
			default:
				data = append(data, bt...)
			}
		}
	}
	return data, err
}

func disableEchoing() error {
	var err error

	termState, err = unix.IoctlGetTermios(
		int(os.Stdin.Fd()), unix.TCGETS,
	)
	if err != nil {
		return err
	}

	lTermState = termState
	termState.Lflag &^= unix.ECHO

	return unix.IoctlSetTermios(
		int(os.Stdin.Fd()), unix.TCSETS, termState,
	)
}

func enableEchoing() error {
	return unix.IoctlSetTermios(
		int(os.Stdin.Fd()), unix.TCSETS, lTermState,
	)
}

func makeRaw() error {
	if termState == nil {
		var err error
		termState, err = unix.IoctlGetTermios(
			int(os.Stdin.Fd()), unix.TCGETS,
		)
		if err != nil {
			return err
		}
	}
	if lTermState == nil {
		lTermState = termState
	}

	termState.Iflag &^= (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	termState.Oflag &^= unix.OPOST
	termState.Lflag &^= (unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN)
	termState.Cflag &^= (unix.CSIZE | unix.PARENB)
	termState.Cflag |= unix.CS8
	termState.Cc[unix.VMIN] = 1
	termState.Cc[unix.VTIME] = 0
	return unix.IoctlSetTermios(
		int(os.Stdin.Fd()), unix.TCSETS, termState,
	)
}

func unmakeRaw() error {
	return unix.IoctlSetTermios(
		int(os.Stdin.Fd()), unix.TCSETS, lTermState,
	)
}

func captureSignals() chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	return ch
}

func stopSignals(ch chan os.Signal) {
	signal.Reset()
	signal.Stop(ch)
}
