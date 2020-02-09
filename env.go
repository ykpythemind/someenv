package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type Env struct {
	Name   string
	Source string
	Dir    string
	Err    error
}

func (e Env) CloneOrPull() error {
	if !isDirExist(e.Dir) {
		if err := os.Mkdir(e.Dir, 0777); err != nil {
			return err
		}
	}

	if isDirExist(path.Join(e.Dir, ".git")) {
		return e.pull()
	} else {
		return e.clone()
	}
}

func (e Env) clone() error {
	fmt.Println("clone " + e.Name)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s; git clone %s .", e.Dir, e.Source))
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	e.readAllAndPrint(stdout)
	e.readAllAndPrint(stderr)

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (e Env) pull() error {
	fmt.Println("pull " + e.Name)

	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s; git pull", e.Dir))
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	e.readAllAndPrint(stdout)
	e.readAllAndPrint(stderr)

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (e Env) readAllAndPrint(reader io.Reader) {
	slurp, _ := ioutil.ReadAll(reader)
	if string(slurp) != "" {
		fmt.Printf("%s: %s\n", e.Name, slurp)
	}
}
