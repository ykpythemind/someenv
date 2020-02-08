package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

	git := path.Join(e.Dir, ".git")
	if isDirExist(git) {
		return e.pull()
	} else {
		return e.clone()
	}
}

func (e Env) clone() error {
	fmt.Println("clone " + e.Name)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s; git clone %s .", e.Dir, e.Source))
	stderr, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	slurp, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s: %s\n", e.Name, slurp)

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (e Env) pull() error {
	fmt.Println("pull " + e.Name)

	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s; git pull", e.Dir))
	stderr, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	slurp, _ := ioutil.ReadAll(stderr)
	if string(slurp) != "" {
		fmt.Printf("%s: %s", e.Name, slurp)
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
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
	}

	err = Run(envs)
	if err != nil {
		log.Fatalf("some errors: %s", err)
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
				log.Printf("err: %s\n", err)
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
