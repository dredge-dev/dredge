package cmd

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
)

type CliCallbacks struct {
	Reader  io.Reader
	Writer  io.Writer
	Verbose *bool
}

func (c CliCallbacks) Log(level api.LogLevel, msg string, args ...interface{}) error {
	if *c.Verbose || (level != api.Debug && level != api.Trace) {
		fmt.Fprintf(c.Writer, "[%s] %s %s\n", time.Now().Format(time.RFC3339), level, fmt.Sprintf(msg, args...))
	}
	return nil
}

func (c CliCallbacks) RequestInput(inputRequests []api.InputRequest) (map[string]string, error) {
	inputs := map[string]string{}
	for _, inputRequest := range inputRequests {
		input, err := c.readInput(inputRequest)
		if err != nil {
			return nil, err
		}
		inputs[inputRequest.Name] = input
	}
	return inputs, nil
}

func (c CliCallbacks) readInput(ir api.InputRequest) (string, error) {
	if ir.Type == api.Text {
		fmt.Printf("%s [%s]: ", ir.Description, ir.Name)
		scanner := bufio.NewScanner(c.Reader)
		if scanner.Scan() {
			value := scanner.Text()
			return value, nil
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	} else if ir.Type == api.Select {
		prompt := promptui.Select{
			Label: fmt.Sprintf("%s [%s]", ir.Description, ir.Name),
			Items: ir.Values,
		}
		_, value, err := prompt.Run()
		if err != nil {
			return "", err
		}
		return value, nil
	}
	return "", fmt.Errorf("InputRequest.Type %d not implemented", ir.Type)
}

func (c CliCallbacks) OpenUrl(url string) error {
	return browser.OpenURL(url)
}

func (c CliCallbacks) Confirm(msg string, args ...interface{}) (bool, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf(msg, args...),
		Items: []string{"yes", "no"},
	}
	_, value, err := prompt.Run()
	if err != nil {
		return false, err
	}
	if value == "yes" {
		return true, nil
	}
	return false, nil
}
