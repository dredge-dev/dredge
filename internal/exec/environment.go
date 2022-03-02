package exec

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/dredge-dev/dredge/internal/config"
)

type Env map[string]string

func NewEnv() Env {
	return Env{}
}

func (e Env) AddVariables(v config.Variables) {
	for key, value := range v {
		if _, ok := e[key]; !ok {
			e[key] = value
		}
	}
}

func (e Env) AddInput(name string, description string, input io.Reader) error {
	var value string
	value = os.Getenv(name)
	if value == "" {
		fmt.Printf("%s [%s]: ", description, name)
		scanner := bufio.NewScanner(input)
		if scanner.Scan() {
			value = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	e[name] = value
	return nil
}

func (e Env) Clone() Env {
	ret := NewEnv()
	for key, value := range e {
		ret[key] = value
	}
	return ret
}
