package shell

import (
	"io"
	"os"
	"os/exec"
)

func Print(r io.Reader, w io.Writer) error {
	cmd := exec.Command("less", "-r", "-B")
	cmd.Stdin = r
	cmd.Stdout = w

	return cmd.Run()
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
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
