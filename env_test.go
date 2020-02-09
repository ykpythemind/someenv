package main

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestRun(t *testing.T) {
	tmpdir := "tmp"
	envs := []Env{
		Env{Dir: path.Join(tmpdir, ".pyenv"), Name: "pyenv", Source: "https://github.com/pyenv/pyenv"},
		Env{Dir: path.Join(tmpdir, ".nodenv"), Name: "nodenv", Source: "https://github.com/nodenv/nodenv"},
	}

	t.Log("clone")
	err := Run(envs)
	if err != nil {
		t.Error(err)
	}
	defer cleanupTmpdir(t, envs)

	t.Log("pull")
	err = Run(envs) // Pull.
	if err != nil {
		t.Error(err)
	}
}

func cleanupTmpdir(t *testing.T, envs []Env) {
	t.Helper()

	for _, e := range envs {
		err := os.RemoveAll(e.Dir)
		if err != nil {
			fmt.Println(err)
		}
	}
}
