package cmd

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/dredge-dev/dredge/internal/callbacks"
	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
)

type CliCallbacks struct {
	reader io.Reader
	writer io.Writer
}

func (c CliCallbacks) Log(level callbacks.LogLevel, msg string) error {
	fmt.Fprintf(c.writer, "[%s] %s %s\n", time.Now().Format(time.RFC822), level, msg)
	return nil
}

func (c CliCallbacks) RequestInput(inputRequests []callbacks.InputRequest) (map[string]string, error) {
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

func (c CliCallbacks) readInput(ir callbacks.InputRequest) (string, error) {
	if ir.Type == callbacks.Text {
		fmt.Printf("%s [%s]: ", ir.Description, ir.Name)
		scanner := bufio.NewScanner(c.reader)
		if scanner.Scan() {
			value := scanner.Text()
			return value, nil
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	} else if ir.Type == callbacks.Select {
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
