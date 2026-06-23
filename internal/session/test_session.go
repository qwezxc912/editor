package session

import (
	"fmt"
	"slices"
)

var (
	nls    = 0
	rows   = 0
	cols   = 0
	ups    = 0
	downs  = 0
	lefts  = 0
	rights = 0
	chars  = 0
)

func (s *Session) testEdit(r rune) {
	switch r {
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
				s.testCursorUp()
			case 'B':
				s.testCursorDown()
			case 'C':
				s.testCursorRight()
			case 'D':
				s.testCursorLeft()
			}
		}
	case '\n':
		s.testNewLine()
	default:
		s.testChar(r)
	}
}

func (s *Session) testCursorUp() {
	fmt.Println("cursorUp")

	fmt.Printf("BEFORE LOGIC\n\n")
	s.DebugInfo()

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

	ups++

	fmt.Printf("AFTER LOGIC\n\n")
	s.DebugInfo()
}

func (s *Session) testCursorDown() {
	fmt.Println("cursorDown")

	fmt.Printf("BEFORE LOGIC\n\n")
	s.DebugInfo()

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

	if s.curs.row < len(s.data) {
		s.curs.row++
	}

	downs++

	fmt.Printf("AFTER LOGIC\n\n")
	s.DebugInfo()
}

func (s *Session) testCursorRight() {
	fmt.Println("cursorRight")

	fmt.Printf("BEFORE LOGIC\n\n")
	s.DebugInfo()

	if s.curs.col > len(s.data[s.curs.row])-1 {
		return
	}

	s.curs.maxIndent++

	s.curs.col++

	rights++

	fmt.Printf("AFTER LOGIC\n\n")
	s.DebugInfo()
}

func (s *Session) testCursorLeft() {
	fmt.Println("cursorLeft")

	fmt.Printf("BEFORE LOGIC\n\n")
	s.DebugInfo()

	if s.curs.col < 1 {
		return
	}

	s.curs.maxIndent--

	s.curs.col--

	lefts++

	fmt.Printf("AFTER LOGIC\n\n")
	s.DebugInfo()
}

func (s *Session) testChar(r rune) {
	fmt.Println("char")

	fmt.Printf("BEFORE LOGIC\n\n")
	s.DebugInfo()

	s.data[s.curs.row] = slices.Insert(
		s.data[s.curs.row],
		s.curs.col,
		r,
	)
	s.curs.col++

	chars++

	fmt.Printf("AFTER LOGIC\n\n")
	s.DebugInfo()
}

func (s *Session) testNewLine() {
	fmt.Println("newLine")

	fmt.Printf("BEFORE LOGIC\n\n")
	s.DebugInfo()

	s.data[s.curs.row] = slices.Insert(s.data[s.curs.row], s.curs.col, '\n')
	s.data = slices.Insert(s.data, s.curs.row, make([]rune, 1, 1))
	s.curs.row++
	s.curs.col = 0
	s.curs.maxIndent = 0

	nls++

	fmt.Printf("AFTER LOGIC\n\n")
	s.DebugInfo()
}

func (s *Session) DebugInfo() {
	fmt.Printf(`
			rows:     %d,
			cols:     %d,
			cursY:    %d,
			cursX:    %d,
			nls:      %d,
			ups:      %d,
			downs:    %d,
			lefts:    %d,
			rights:   %d,
			chars:    %d,
			rowsQuan: %d,
			rowLen:   %d,
		`,
		rows,
		cols,
		s.curs.row,
		s.curs.col,
		nls,
		ups,
		downs,
		lefts,
		rights,
		chars,
		len(s.data),
		len(s.data[s.curs.col]),
	)
}
