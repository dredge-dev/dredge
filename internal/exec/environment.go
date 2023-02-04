package exec

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/manifoldco/promptui"
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

func (e Env) AddInputs(inputs map[string]string) {
	for key, value := range inputs {
		e[key] = value
	}
}

func (e Env) AddInput(input config.Input, reader io.Reader) error {
	var value string
	value = os.Getenv(input.Name)
	if input.Type == "" || input.Type == config.INPUT_TEXT {
		if value == "" {
			fmt.Printf("%s [%s]: ", input.Description, input.Name)
			scanner := bufio.NewScanner(reader)
			if scanner.Scan() {
				value = scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				return err
			}
		}
		e[input.Name] = value
		return nil
	} else if input.Type == config.INPUT_SELECT {
		if value == "" {
			var err error
			prompt := promptui.Select{
				Label: fmt.Sprintf("%s [%s]", input.Description, input.Name),
				Items: input.Values,
			}
			_, value, err = prompt.Run()
			if err != nil {
				return err
			}
		}
		if !input.HasValue(value) {
			return fmt.Errorf("Invalid value (%s) for Input %s", value, input.Name)
		}
		e[input.Name] = value
		return nil
	}
	return fmt.Errorf("Type %s not implemented", input.Type)
}

func (e Env) Clone() Env {
	ret := NewEnv()
	for key, value := range e {
		ret[key] = value
	}
	return ret
}
