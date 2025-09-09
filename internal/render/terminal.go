package render

import (
	"os"

	"golang.org/x/sys/unix"
)

type Terminal struct {
	fd    int
	state unix.Termios
}

func InitTerminal() (*Terminal, error) {
	ret := &Terminal{fd: int(os.Stdin.Fd())}

	termiosPtr, err := unix.IoctlGetTermios(ret.fd, unix.TIOCGETA)
	if err != nil {
		return nil, err
	}
	ret.state = *termiosPtr

	termios := *termiosPtr
	termios.Lflag &^= unix.ECHO | unix.ICANON
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

	if err = unix.IoctlSetTermios(ret.fd, unix.TIOCSETA, &termios); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *Terminal) Restore() error {
	return unix.IoctlSetTermios(t.fd, unix.TIOCSETA, &t.state)
}
