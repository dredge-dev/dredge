package cmd

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGetDredgeFile(t *testing.T) {
	df := config.DredgeFile{
		Buckets: []config.Bucket{
			{
				Name: "b1",
			},
		},
	}

	tmpFile := fmt.Sprintf("./tmp-test-%d", rand.Intn(100000))
	defer os.Remove(tmpFile)

	data, err := yaml.Marshal(df)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(tmpFile, data, 0644)
	if err != nil {
		panic(err)
	}

	tests := map[string]struct {
		target     string
		dredgeFile *config.DredgeFile
		errMsg     string
	}{
		"local file": {
			target:     tmpFile,
			dredgeFile: &df,
			errMsg:     "",
		},
		"file not found": {
			target:     "./non-existing-file",
			dredgeFile: nil,
			errMsg:     "Error while parsing ./non-existing-file: open ./non-existing-file: no such file or directory",
		},
		"unsupported": {
			target:     "/hello",
			dredgeFile: nil,
			errMsg:     "Targets should start with ./",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		d, err := getDredgeFile(test.target)
		if test.dredgeFile == nil {
			assert.Nil(t, d)
		} else {
			assert.Equal(t, test.dredgeFile.Buckets[0].Name, d.Buckets[0].Name)
		}
		if test.errMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
	}
}
