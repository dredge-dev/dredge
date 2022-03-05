package exec

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	osExec "os/exec"
	"path/filepath"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
)

const (
	DefaultDredgefileName  = "Dredgefile"
	DefaultDredgeRepo      = "https://github.com/dredge-dev/dredge-repo.git"
	LocalDredgeStorage     = ".dredge"
	LocalDredgeRepoStorage = LocalDredgeStorage + "/repo/"
)

func MergeSources(parent config.SourcePath, child config.SourcePath) config.SourcePath {
	if child == "" {
		return parent
	}
	if parent == "" {
		return child
	}
	c := string(child)
	if strings.HasPrefix(c, "./") {
		p := string(parent)
		parentPath := strings.Split(p, "/")
		parentDir := parentPath[:len(parentPath)-1]
		parts := append(parentDir, c[2:])
		return config.SourcePath(strings.Join(parts, "/"))
	}
	return child
}

func resolvePath(source config.SourcePath) (string, error) {
	s := string(source)
	if len(s) >= 2 && s[0] == '.' && os.IsPathSeparator(s[1]) {
		return s, nil
	}
	if !strings.Contains(s, ":") {
		s = DefaultDredgeRepo + ":" + s
	}
	split := strings.LastIndex(s, ":")
	return resolveRepo(s[:split], s[split+1:])
}

func resolveRepo(repo, path string) (string, error) {
	repoPath := resolveRepoPath(repo)
	if _, err := os.Stat(LocalDredgeRepoStorage); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(LocalDredgeRepoStorage, 0755)
		if err != nil {
			return "", err
		}
	}
	if _, err := os.Stat(repoPath); errors.Is(err, os.ErrNotExist) {
		cmd := osExec.Command("/bin/bash", "-c", fmt.Sprintf("git clone --depth 1 %s %s", repo, repoPath))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(repoPath, path), nil
}

func resolveRepoPath(repo string) string {
	hash := sha256.Sum256([]byte(repo))
	dir := hex.EncodeToString(hash[:])
	return LocalDredgeRepoStorage + dir
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
