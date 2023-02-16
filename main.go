package main

import (
	"errors"
	"log"
	"os"

	"github.com/dredge-dev/dredge/internal/cmd"
	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/dredge-dev/dredge/internal/resource"
)

const DefaultDredgefilePath = "./" + exec.DefaultDredgefileName

func main() {
	var de *exec.DredgeExec

	source := DefaultDredgefilePath
	c := cmd.CliCallbacks{Reader: os.Stdin, Writer: os.Stdout}
	rd := resource.GetDefaultResourceDefinitions()

	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		de = exec.EmptyExec(config.SourcePath(source), rd, c)
	} else {
		de, err = exec.NewExec(config.SourcePath(source), rd, c)
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
