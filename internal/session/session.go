package session

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

const (
	edit   = "edit"
	cmd    = "cmd"
	visual = "visual"
)

type Session struct {
	curs   *curs
	reader *bufio.Reader
	data   [][]rune
	mode   string
	file   *os.File
}

type curs struct {
	maxIndent int
	col       int
	row       int
}

func New(file *os.File) *Session {
	reader := bufio.NewReader(os.Stdin)

	cursor := &curs{col: 0, row: 1, maxIndent: 0}

	d := make([][]rune, 1, 1)

	return &Session{
		curs:   cursor,
		reader: reader,
		data:   d,
		mode:   edit,
		file:   file,
	}
}

func (s *Session) StartRead() {
	for {
		r, _, err := s.reader.ReadRune()
		if err != nil {
			continue
		}
		if r == '`' {
			break
		}
		switch s.mode {
		case edit:
			s.edit(r)
		}
	}
}

func (s *Session) UploadFile() {
	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		s.data = append(s.data, []rune(scanner.Text()+"\n"))
	}

	s.curs.row = len(s.data) - 1
}

func (s *Session) WriteFile() error {
	for _, line := range s.data {
		if _, err := s.file.Write([]byte(string(line))); err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) PrintData() {
	for _, line := range s.data {
		fmt.Println(line)
	}
}

func (s *Session) edit(r rune) {
	switch r {
	case 0:
		break
	case '\x1b':
		r, _, err := s.reader.ReadRune()
		if err != nil {
			break
		}
		if r == '[' {
			r, _, err := s.reader.ReadRune()
			if err != nil {
				break
			}
			switch r {
			case 'A':
				s.cursorUp()
			case 'B':
				s.cursorDown()
			case 'C':
				s.cursorRight()
			case 'D':
				s.cursorLeft()
			}
		}
	case '\n':
		s.newLine()
	default:
		s.char(r)
	}
}

func (s *Session) newLine() {
	lf := s.data[s.curs.row][s.curs.col:]

	s.data[s.curs.row] = slices.Insert(s.data[s.curs.row], s.curs.col, '\n')

	s.data = slices.Insert(s.data, s.curs.row+1, lf)

	s.curs.row++
	s.curs.col = len(lf)
	s.curs.maxIndent = len(lf)

	fmt.Printf("\x1b[0K\n")
	fmt.Printf("%s", string(lf))
}

func (s *Session) char(r rune) {
	s.data[s.curs.row] = slices.Insert(
		s.data[s.curs.row],
		s.curs.col,
		r,
	)
	fmt.Printf("\x1b[1@%c", r)
	s.curs.col++
}

func (s *Session) cursorUp() {
	if s.curs.row <= 1 {
		return
	}

	currIndent := s.curs.col
	if currIndent > s.curs.maxIndent {
		s.curs.maxIndent = currIndent
	}

	if s.curs.col > len(s.data[s.curs.row-1]) {
		s.curs.col = len(s.data[s.curs.row-1])
	}
	if len(s.data[s.curs.row-1]) >= s.curs.maxIndent {
		s.curs.col = s.curs.maxIndent
	} else {
		s.curs.col = len(s.data[s.curs.row-1])
	}

	s.curs.row--
	fmt.Printf("\x1b[%d;%df", s.curs.row, s.curs.col)
	fmt.Printf("%d", s.curs.row)
}

func (s *Session) cursorDown() {
	if s.curs.row >= len(s.data)-1 {
		return
	}

	currIndent := s.curs.col
	if currIndent > s.curs.maxIndent {
		s.curs.maxIndent = currIndent
	}

	if s.curs.col > len(s.data[s.curs.row+1]) {
		s.curs.col = len(s.data[s.curs.row+1])
	}
	if len(s.data[s.curs.row+1]) >= s.curs.maxIndent {
		s.curs.col = s.curs.maxIndent
	} else {
		s.curs.col = len(s.data[s.curs.row+1])
	}

	s.curs.row++
	fmt.Printf("\x1b[%d;%df", s.curs.row, s.curs.col)
	fmt.Printf("%d", s.curs.row)
}

func (s *Session) cursorRight() {
	if s.curs.col > len(s.data[s.curs.row])-1 {
		return
	}

	s.curs.maxIndent++

	s.curs.col++
	fmt.Printf("\x1b[1C")
}

func (s *Session) cursorLeft() {
	if s.curs.col < 1 {
		return
	}

	s.curs.maxIndent--

	s.curs.col--
	fmt.Printf("\x1b[1D")
}

func (s *Session) backspace() {

}
