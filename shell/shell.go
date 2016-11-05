package shell

import (
	"fmt"
	"io"
	"strings"

	readline "gopkg.in/readline.v1"
)

const defaultPromp = "gli> "

type handlerFunc func() (string, error)

type shell struct {
	funcs  map[string]handlerFunc
	writer io.Writer
}

func New(w io.Writer) *shell {
	s := &shell{
		funcs:  make(map[string]handlerFunc),
		writer: w,
	}
	return s
}

func (s *shell) Case(name string, f handlerFunc) error {
	if _, exists := s.funcs[name]; exists {
		return fmt.Errorf("%s already exists", name)
	}
	s.funcs[name] = f
	return nil
}

func (s *shell) Start() error {
	var pcItems []readline.PrefixCompleterInterface
	for k := range s.funcs {
		pcItems = append(pcItems, readline.PcItem(k))
	}

	completer := readline.NewPrefixCompleter(pcItems...)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       defaultPromp,
		AutoComplete: completer,
	})
	if err != nil {
		return err
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		if err = s.handle(line); err != nil {
			s.println(err)
			break
		}
	}

	return nil
}

func (s *shell) handle(line string) error {
	if _, exists := s.funcs[strings.TrimSpace(line)]; !exists {
		return nil
	}

	result, err := s.funcs[strings.TrimSpace(line)]()
	if err != nil {
		return err
	}
	s.println(result)

	return nil
}

func (s *shell) println(v ...interface{}) {
	fmt.Fprintln(s.writer, v...)
}
