package main

import (
	"fmt"
	"os"
	"path"
	"sync"
)

// EnvList is *env list
var EnvList []Env

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	homeDir = path.Join(homeDir, "playground")

	EnvList = []Env{
		Env{Dir: path.Join(homeDir, ".pyenv"), Name: "pyenv", Source: "https://github.com/pyenv/pyenv"},
		Env{Dir: path.Join(homeDir, ".nodenv"), Name: "nodenv", Source: "https://github.com/nodenv/nodenv"},
		Env{Dir: path.Join(homeDir, ".goenv"), Name: "goenv", Source: "https://github.com/syndbg/goenv"},
		Env{Dir: path.Join(homeDir, ".rbenv"), Name: "rbenv", Source: "https://github.com/rbenv/rbenv"},
	}
}

func main() {
	err := Run(EnvList)
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
