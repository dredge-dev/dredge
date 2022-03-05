package exec

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestMergeSources(t *testing.T) {
	tests := map[string]struct {
		parent config.SourcePath
		child  config.SourcePath
		result config.SourcePath
	}{
		"parent without dir, child without dir": {
			parent: "./test.Dredgefile",
			child:  "./second.Dredgefile",
			result: "./second.Dredgefile",
		},
		"parent with dir, child without dir": {
			parent: "./parent/test.Dredgefile",
			child:  "./second.Dredgefile",
			result: "./parent/second.Dredgefile",
		},
		"parent without dir, child with dir": {
			parent: "./test.Dredgefile",
			child:  "./child/second.Dredgefile",
			result: "./child/second.Dredgefile",
		},
		"parent with dir, child with dir": {
			parent: "./parent/test.Dredgefile",
			child:  "./child/second.Dredgefile",
			result: "./parent/child/second.Dredgefile",
		},
		"empty parent, child": {
			parent: "",
			child:  "./child/Dredgefile",
			result: "./child/Dredgefile",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		assert.Equal(t, test.result, MergeSources(test.parent, test.child))
	}
}

func TestResolvePath(t *testing.T) {
	tests := map[string]struct {
		source config.SourcePath
		path   string
		errMsg string
	}{
		"starting with ./": {
			source: "./test",
			path:   "./test",
		},
		"not starting with ./": {
			source: "test",
			errMsg: "sources should start with ./",
		},
	}
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		result, err := resolvePath(test.source)
		if test.errMsg == "" {
			assert.Nil(t, err)
			assert.Equal(t, test.path, result)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
	}
}

func TestResolveDredgeFilePath(t *testing.T) {
	err := os.Mkdir("./drg-resolve-test-dir", 0755)
	assert.Nil(t, err)
	defer os.RemoveAll("./drg-resolve-test-dir")

	err = ioutil.WriteFile("./drg-resolve-test-dir/Dredgefile", []byte("workflows:"), 0644)
	assert.Nil(t, err)

	err = os.Mkdir("./drg-resolve-test-dir-empty", 0755)
	assert.Nil(t, err)
	defer os.RemoveAll("./drg-resolve-test-dir-empty")

	tests := map[string]struct {
		source     config.SourcePath
		fullSource config.SourcePath
		path       string
		errMsg     string
	}{
		"non existing file": {
			source: "./non-existing",
			errMsg: "stat ./non-existing: no such file or directory",
		},
		"full path to file": {
			source:     "./drg-resolve-test-dir/Dredgefile",
			fullSource: "./drg-resolve-test-dir/Dredgefile",
			path:       "./drg-resolve-test-dir/Dredgefile",
		},
		"path to directory without /": {
			source:     "./drg-resolve-test-dir",
			fullSource: "./drg-resolve-test-dir/Dredgefile",
			path:       "./drg-resolve-test-dir/Dredgefile",
		},
		"path to directory with /": {
			source:     "./drg-resolve-test-dir/",
			fullSource: "./drg-resolve-test-dir/Dredgefile",
			path:       "./drg-resolve-test-dir/Dredgefile",
		},
		"path to directory without Dredgefile": {
			source: "./drg-resolve-test-dir-empty",
			errMsg: "stat ./drg-resolve-test-dir-empty/Dredgefile: no such file or directory",
		},
	}
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		fullSource, path, err := resolveDredgeFilePath(test.source)
		if test.errMsg == "" {
			assert.Nil(t, err)
			assert.Equal(t, test.path, path)
			assert.Equal(t, test.fullSource, fullSource)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
	}
}
