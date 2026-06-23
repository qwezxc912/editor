package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"

	"gihub.com/qweq1232/nedovim/internal/session"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("not enough args")
	}
	name := os.Args[len(os.Args)-1]

	file, err := openFile(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	setup(file)

	s := session.New(file)
	s.UploadFile()
	s.StartRead()

	if err = s.WriteFile(); err != nil {
		log.Fatal(err)
	}
	cleenup()
	s.PrintData()
}

func openFile(name string) (*os.File, error) {
	var file *os.File

	if _, err := os.Stat(name); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			file, err = os.Create(name)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	file, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func cleenup() {
	fmt.Print("\x1b[H")
	fmt.Print("\x1b[J")
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func setup(file *os.File) {
	fmt.Print("\x1b[H")
	fmt.Print("\x1b[J")

	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	file.WriteTo(os.Stdout)
}
