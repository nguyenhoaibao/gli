package shell

import (
	"io"
	"os"
	"os/exec"
)

func Print(r io.Reader, w io.Writer) error {
	cmd := exec.Command("less", "-r")
	cmd.Stdin = r
	cmd.Stdout = w

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func ClearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

type HandlerFunc func(args ...string) (io.Reader, error)

func (f HandlerFunc) Handle(args ...string) (io.Reader, error) {
	return f(args...)
}

func addDefaultHandlers(s *shell) {
	s.Register("clear", clearFunc(s))
}

func clearFunc(s *shell) HandlerFunc {
	return func(args ...string) (io.Reader, error) {
		err := ClearScreen()
		return nil, err
	}
}
