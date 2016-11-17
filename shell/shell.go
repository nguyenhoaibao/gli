package shell

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	readline "gopkg.in/readline.v1"
)

const defaultPromp = "gli> "

type shellReader struct {
	rl *readline.Instance
}

type shell struct {
	funcs          map[string]HandlerFunc
	genericHandler HandlerFunc
	reader         *shellReader
	writer         io.Writer
}

func New(w io.Writer) *shell {
	rl, err := readline.New(defaultPromp)
	if err != nil {
		return nil
	}

	s := &shell{
		funcs: make(map[string]HandlerFunc),
		reader: &shellReader{
			rl: rl,
		},
		writer: w,
	}
	addDefaultHandlers(s)
	return s
}

func (s *shell) Register(name string, f HandlerFunc) error {
	if _, exists := s.funcs[name]; exists {
		return fmt.Errorf("%s already exists", name)
	}
	s.funcs[name] = f
	return nil
}

func (s *shell) RegisterGeneric(f HandlerFunc) error {
	if s.genericHandler != nil {
		return errors.New("Generic function already registered")
	}
	s.genericHandler = f
	return nil
}

func (s *shell) Start() error {
	var (
		sCmds []string
		mCmds = make(map[string]map[string]string)
	)

	for name := range s.funcs {
		if strings.Index(name, "_") == -1 {
			sCmds = append(sCmds, name)
			continue
		}

		names := strings.Split(name, "_")
		root, sub := names[0], names[1]
		if _, exists := mCmds[root]; !exists {
			mCmds[root] = make(map[string]string)
		}
		if _, exists := mCmds[root][sub]; !exists {
			mCmds[root][sub] = sub
		}
	}

	var pcItems []readline.PrefixCompleterInterface

	for root, subs := range mCmds {
		var subPcItems []readline.PrefixCompleterInterface

		for sub := range subs {
			subPcItems = append(subPcItems, readline.PcItem(sub))
		}
		pcItems = append(pcItems, readline.PcItem(root, subPcItems...))
	}

	for _, name := range sCmds {
		pcItems = append(pcItems, readline.PcItem(name))
	}

	completer := readline.NewPrefixCompleter(pcItems...)
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
	args := strings.Fields(line)

	handled, err := s.handleCommand(args)
	if err != nil {
		return err
	}
	if handled {
		return nil
	}

	return s.genericHandle(args)
}

func (s *shell) handleCommand(args []string) (bool, error) {
	line := strings.Join(args, "_")

	for name, handler := range s.funcs {
		if strings.Index(line, name) == -1 {
			continue
		}

		var (
			result io.Reader
			err    error
		)

		if strings.Index(name, "_") == -1 {
			result, err = handler.Handle(args[1:]...)
		} else {
			result, err = handler.Handle(args[2:]...)
		}

		if err != nil {
			return true, err
		}
		if result != nil {
			Print(result, s.writer)
		}
		return true, nil
	}
	return false, nil
}

func (s *shell) genericHandle(args []string) error {
	if s.genericHandler == nil {
		return nil
	}

	result, err := s.genericHandler.Handle(args...)
	if err != nil {
		return err
	}
	if result != nil {
		Print(result, s.writer)
	}
	return nil
}
