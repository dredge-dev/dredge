package main

import (
	"fmt"
	"os"

	"github.com/dredge-dev/dredge/internal/cmd"
	"github.com/dredge-dev/dredge/internal/config"
)

func main() {
	dredgeFile, err := config.GetDredgeFile("./Dredgefile")
	if err != nil {
		fmt.Printf("Error while reading Dredgefile: %s\n", err)
		os.Exit(1)
	}

	cmd.Init(dredgeFile)
	cmd.Execute()
}
