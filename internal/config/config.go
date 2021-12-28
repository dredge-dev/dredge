package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

type DredgeFile struct {
	Env Env
	Workflows []Workflow
}

type Env struct {
	Variables map[string]string
	Runtimes []Runtime
}

type Runtime struct {
	Name string
	Type string
	Image string
	Cache []string
}

type Workflow struct {
	Name string
	Description string
	Steps []Step
}

type Step struct {
	Exec string
	Runtime string
}

func ReadDredgeFile(filename string) (*DredgeFile, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	dredgeFile := &DredgeFile{}
	err = yaml.Unmarshal(buf, dredgeFile)
	if err != nil {
		return nil, err
	}

	return dredgeFile, nil
}
