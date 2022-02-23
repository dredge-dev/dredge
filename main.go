package main

import (
	"errors"
	"log"
	"os"

	"github.com/dredge-dev/dredge/internal/cmd"
	"github.com/dredge-dev/dredge/internal/config"
)

const defaultDredgeFile = "./Dredgefile"

func main() {
	var dredgeFile *config.DredgeFile
	if _, err := os.Stat(defaultDredgeFile); errors.Is(err, os.ErrNotExist) {
		dredgeFile = &config.DredgeFile{}
	} else {
		if dredgeFile, err = config.GetDredgeFile(defaultDredgeFile); err != nil {
			log.Fatalf("Error while reading Dredgefile: %s\n", err)
		}
	}

	err := cmd.Init(dredgeFile)
	if err != nil {
		log.Fatalf("Error during init: %v", err)
	}
	cmd.Execute()
}
