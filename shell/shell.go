package shell

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	readline "gopkg.in/readline.v1"
)

const defaultPromp = "gli> "

type handlerFunc func(args ...string) (io.Reader, error)

type shellReader struct {
	rl *readline.Instance
}

type shell struct {
	funcs    map[string]handlerFunc
	commands map[string]map[string]handlerFunc
	reader   *shellReader
	writer   io.Writer
}

func New(w io.Writer) *shell {
	rl, err := readline.New(defaultPromp)
	if err != nil {
		return nil
	}

	s := &shell{
		funcs:    make(map[string]handlerFunc),
		commands: make(map[string]map[string]handlerFunc),
		reader: &shellReader{
			rl: rl,
		},
		writer: w,
	}
	return s
}

func (s *shell) Register(name string, f handlerFunc) error {
	if _, exists := s.funcs[name]; exists {
		return fmt.Errorf("%s already exists", name)
	}
	s.funcs[name] = f
	return nil
}

func (s *shell) Start() error {
	for name, handler := range s.funcs {
		names := strings.Split(name, "_")
		root, sub := names[0], names[1]

		if _, exists := s.commands[root]; !exists {
			s.commands[root] = make(map[string]handlerFunc)
		}
		if _, exists := s.commands[root][sub]; !exists {
			s.commands[root][sub] = handler
		}
	}

	var rootPcItems []readline.PrefixCompleterInterface
	for root, subHandler := range s.commands {
		var subPcItems []readline.PrefixCompleterInterface

		for sub, _ := range subHandler {
			subPcItems = append(subPcItems, readline.PcItem(sub))
		}
		rootPcItems = append(rootPcItems, readline.PcItem(root, subPcItems...))
	}
	completer := readline.NewPrefixCompleter(rootPcItems...)

	s.reader.rl.SetConfig(&readline.Config{
		Prompt:       defaultPromp,
		AutoComplete: completer,
	})
	s.reader.rl.Refresh()
	defer s.reader.rl.Close()

	for {
		line, err := s.reader.rl.Readline()
		if err != nil {
			return err
		}
		if err = s.handle(line); err != nil {
			log.Fatal(err)
			break
		}
	}

	return nil
}

func (s *shell) handle(line string) error {
	args := strings.Split(strings.TrimSpace(line), " ")

	_, exists := s.commands[args[0]]
	if !exists {
		return nil
	}
	handler, exists := s.commands[args[0]][args[1]]
	if !exists {
		return nil
	}

	result, err := handler(args[2:]...)
	if err != nil {
		return err
	}
	return s.printResult(result)
}

func (s *shell) printResult(r io.Reader) error {
	cmd := exec.Command("less", "-r")
	cmd.Stdin = r
	cmd.Stdout = s.writer

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
