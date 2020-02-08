package main

import (
	"path"
	"testing"
)

func TestRun(t *testing.T) {
	tmpdir := "tmp"
	envs := []Env{
		Env{Dir: path.Join(tmpdir, ".pyenv"), Name: "pyenv", Source: "https://github.com/pyenv/pyenv"},
		Env{Dir: path.Join(tmpdir, ".nodenv"), Name: "nodenv", Source: "https://github.com/nodenv/nodenv"},
	}

	err := Run(envs)
	if err != nil {
		t.Error(err)
	}
}
