package main

import (
	"errors"
	"log"
	"os"

	"github.com/dredge-dev/dredge/internal/cmd"
	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
)

const DefaultDredgefilePath = "./" + exec.DefaultDredgefileName

func main() {
	source := DefaultDredgefilePath
	var de *exec.DredgeExec
	c := cmd.CliCallbacks{Reader: os.Stdin, Writer: os.Stdout}

	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		de = exec.EmptyExec(config.SourcePath(source), c)
	} else {
		de, err = exec.NewExec(config.SourcePath(source), c)
		if err != nil {
			log.Fatalf("Error while reading Dredgefile: %s\n", err)
		}
	}

	err := cmd.Init(de)
	if err != nil {
		log.Fatalf("Error during init: %v", err)
	}
	cmd.Execute()
}
