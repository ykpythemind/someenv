package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sync"
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

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	homeDir = path.Join(homeDir, "playground")

	envs := []Env{
		Env{Dir: path.Join(homeDir, ".pyenv"), Name: "pyenv", Source: "https://github.com/pyenv/pyenv"},
		Env{Dir: path.Join(homeDir, ".nodenv"), Name: "nodenv", Source: "https://github.com/nodenv/nodenv"},
		Env{Dir: path.Join(homeDir, ".goenv"), Name: "goenv", Source: "https://github.com/syndbg/goenv"},
	}

	err = Run(envs)
	if err != nil {
		fmt.Printf("some errors: %s\n", err)
	}

	fmt.Println("finished.")
}

func Run(envList []Env) error {
	wg := sync.WaitGroup{}

	for _, e := range envList {
		wg.Add(1)

		go func(e Env) {
			defer wg.Done()

			if err := e.CloneOrPull(); err != nil {
				e.Err = err
				fmt.Printf("err: %s\n", err)
			}
		}(e)
	}

	wg.Wait()

	for _, e := range envList {
		if e.Err != nil {
			return e.Err
		}
	}

	return nil
}

func isDirExist(pathname string) bool {
	info, err := os.Stat(pathname)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}
