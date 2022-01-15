package workflow

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Env map[string]string

func NewEnv() Env {
	return Env{}
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
