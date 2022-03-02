package main

import (
	"errors"
	"log"
	"os"

	"github.com/dredge-dev/dredge/internal/cmd"
	"github.com/dredge-dev/dredge/internal/exec"
)

const defaultDredgefilePath = "./Dredgefile"

func main() {
	source := defaultDredgefilePath
	var de *exec.DredgeExec

	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		de = exec.EmptyExec(source)
	} else {
		de, err = exec.NewExec(source)
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
