package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"slices"
)

type session struct {
	content [][]rune
	name    string
	cursor  *cursor
	reader  *bufio.Reader
	file    *os.File
}

type cursor struct {
	x int
	y int
}

func newSession(name string, file *os.File) *session {
	reader := bufio.NewReader(os.Stdin)

	cursor := newCursor()

	content := make([][]rune, 1, 1)

	return &session{
		content: content,
		name:    name,
		cursor:  cursor,
		reader:  reader,
		file:    file,
	}
}

func (s *session) read() error {
loop:
	for {
		r, _, err := s.reader.ReadRune()
		if err != nil {
			break
		}
		switch r {
		case '`':
			break loop
		case '\n':
			fmt.Print("\n")
			s.content = insert(s.content, s.cursor.x, s.cursor.y, r)
			s.cursor.y++
			s.cursor.x = 0
			s.content = append(s.content, []rune{})
		default:
			fmt.Printf("%c", r)
			s.content = insert(s.content, s.cursor.x, s.cursor.y, r)
			s.cursor.x++
		}
	}
	return nil
}

func newCursor() *cursor {
	return &cursor{x: 0, y: 0}
}

func setup(file *os.File) {
	fmt.Print("\x1b[H")
	fmt.Print("\x1b[J")
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	file.WriteTo(os.Stdout)
}

func cleenup() {
	fmt.Print("\x1b[H")
	fmt.Print("\x1b[J")
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

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
	defer cleenup()

	session := newSession(name, file)
	if err = session.read(); err != nil {
		log.Fatal(err)
	}

	if err = session.writeFile(); err != nil {
		log.Fatal(err)
	}
}

func openFile(name string) (*os.File, error) {
	var file *os.File

	info, err := os.Stat(name)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			file, err = os.Create(name)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if info.IsDir() {
		return nil, errors.New("not a file")
	}

	file, err = os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *session) writeFile() error {
	for _, line := range s.content {
		if _, err := s.file.Write([]byte(string(line))); err != nil {
			return err
		}
	}

	return nil
}

func insert(cont [][]rune, x int, y int, r rune) [][]rune {
	line := cont[y]

	if x == len(line) {
		line = append(line, r)

		cont[y] = line

		return cont
	}

	line = slices.Insert(line, x, r)

	cont[y] = line

	return cont
}
