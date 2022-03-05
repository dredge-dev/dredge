package exec

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
)

const DefaultDredgefileName = "Dredgefile"

func MergeSources(parent config.SourcePath, child config.SourcePath) config.SourcePath {
	c := string(child)
	p := string(parent)
	if c == "" {
		return parent
	} else if strings.HasPrefix(c, "./") {
		if strings.HasPrefix(p, "./") {
			parentPath := strings.Split(p, "/")
			parentDir := parentPath[:len(parentPath)-1]
			parts := append(parentDir, c[2:])
			return config.SourcePath(strings.Join(parts, "/"))
		}
	}
	return child
}

func resolvePath(source config.SourcePath) (string, error) {
	s := string(source)
	if !strings.HasPrefix(s, "./") {
		return "", fmt.Errorf("sources should start with ./")
	}
	return s, nil
}

func resolveDredgeFilePath(source config.SourcePath) (config.SourcePath, string, error) {
	path, err := resolvePath(source)
	if err != nil {
		return "", "", err
	}
	stat, err := os.Stat(path)
	if err != nil {
		return "", "", err
	}
	if stat.IsDir() {
		fullSource := string(source)
		if !os.IsPathSeparator(path[len(path)-1]) {
			fullSource += string(os.PathSeparator)
		}
		fullSource += DefaultDredgefileName
		return resolveDredgeFilePath(config.SourcePath(fullSource))
	}
	return source, path, nil
}

func readSource(source config.SourcePath) ([]byte, error) {
	path, err := resolvePath(source)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(path)
}

func ReadDredgeFile(source config.SourcePath) (config.SourcePath, *config.DredgeFile, error) {
	fullSource, path, err := resolveDredgeFilePath(source)
	if err != nil {
		return "", nil, err
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	df, err := config.NewDredgeFile(content)
	if err != nil {
		return "", nil, err
	}
	return fullSource, df, nil
}
