package tools

type Scanner struct {
	input    string
	position int
	line     int
	column   int
<<<<<<<< Updated upstream:htmlparser/lexer.go
========
	Ch       rune
>>>>>>>> Stashed changes:tools/scanner.go
}

func NewScanner(input string) *Scanner {
	return &Scanner{
		input:  input,
		line:   1,
		column: 1,
	}
<<<<<<<< Updated upstream:htmlparser/lexer.go
}

func (s *Scanner) Current() rune {
	if s.position >= len(s.input) {
		return 0
	}
	return rune(s.input[s.position])
========
	s.updateCh()
	return s
}

func (s *Scanner) updateCh() {
	if s.position >= len(s.input) {
		s.Ch = 0
	} else {
		s.Ch = rune(s.input[s.position])
	}
}

func (s *Scanner) Current() rune {
	return s.Ch
>>>>>>>> Stashed changes:tools/scanner.go
}

func (s *Scanner) PeekNext() rune {
	if s.position+1 >= len(s.input) {
		return 0
	}
	return rune(s.input[s.position+1])
}

func (s *Scanner) Previous() rune {
	if s.position <= 0 {
		return 0
	}
	return rune(s.input[s.position-1])
}

func (s *Scanner) Take() rune {
	if s.position >= len(s.input) {
		return 0
	}

<<<<<<<< Updated upstream:htmlparser/lexer.go
	current := s.Current()
========
	current := s.Ch
>>>>>>>> Stashed changes:tools/scanner.go
	s.position++

	if current == '\n' {
		s.line++
		s.column = 1
	} else {
		s.column++
	}

	return current
}

func (s *Scanner) Retreat() rune {
	if s.position <= 0 {
		return 0
	}

	s.position--
	current := s.Current()

<<<<<<<< Updated upstream:htmlparser/lexer.go
	if current == '\n' {
========
	if s.Ch == '\n' {
>>>>>>>> Stashed changes:tools/scanner.go
		s.line--
		s.column = s.columnAt(s.position)
	} else {
		s.column--
	}

<<<<<<<< Updated upstream:htmlparser/lexer.go
	return current
========
	return s.Ch
>>>>>>>> Stashed changes:tools/scanner.go
}

func (s *Scanner) columnAt(pos int) int {
	col := 1
	for i := pos - 1; i >= 0 && s.input[i] != '\n'; i-- {
		col++
	}
	return col
}

func (s *Scanner) Peek(offset int) rune {
	pos := s.position + offset
	if pos < 0 || pos >= len(s.input) {
		return 0
	}
	return rune(s.input[pos])
}

func (s *Scanner) Skip() {
	s.Take()
}

func (s *Scanner) SkipN(n int) {
	for i := 0; i < n && !s.EOF(); i++ {
		s.Take()
	}
}

func (s *Scanner) EOF() bool {
	return s.position >= len(s.input)
}

func (s *Scanner) Position() int {
	return s.position
}

func (s *Scanner) Line() int {
	return s.line
}

func (s *Scanner) Column() int {
	return s.column
}

func (s *Scanner) Mark() int {
	return s.position
}

func (s *Scanner) Reset(mark int) {
	if mark < 0 || mark > len(s.input) {
		return
	}

	s.position = mark
	s.line = 1
	s.column = 1

	for i := 0; i < mark; i++ {
		if s.input[i] == '\n' {
			s.line++
			s.column = 1
		} else {
			s.column++
		}
	}
}

func (s *Scanner) Slice(start, end int) string {
	if start < 0 || end > len(s.input) || start > end {
		return ""
	}
	return s.input[start:end]
}

func (s *Scanner) SliceFrom(start int) string {
	return s.Slice(start, s.position)
}

func (s *Scanner) Match(expected rune) bool {
<<<<<<<< Updated upstream:htmlparser/lexer.go
	if s.Current() == expected {
========
	if s.Ch == expected {
>>>>>>>> Stashed changes:tools/scanner.go
		s.Take()
		return true
	}
	return false
}

func (s *Scanner) MatchAny(chars ...rune) bool {
	current := s.Current()
	for _, char := range chars {
<<<<<<<< Updated upstream:htmlparser/lexer.go
		if current == char {
========
		if s.Ch == char {
>>>>>>>> Stashed changes:tools/scanner.go
			s.Take()
			return true
		}
	}
	return false
}

func (s *Scanner) MatchString(expected string) bool {
	if s.position+len(expected) > len(s.input) {
		return false
	}

	if s.input[s.position:s.position+len(expected)] == expected {
		for range expected {
			s.Take()
		}
		return true
	}
	return false
}

func (s *Scanner) ConsumeWhile(predicate func(rune) bool) string {
	start := s.position
<<<<<<<< Updated upstream:htmlparser/lexer.go
	for !s.EOF() && predicate(s.Current()) {
========
	for !s.EOF() && predicate(s.Ch) {
>>>>>>>> Stashed changes:tools/scanner.go
		s.Take()
	}
	return s.input[start:s.position]
}

func (s *Scanner) ConsumeUntil(predicate func(rune) bool) string {
	start := s.position
<<<<<<<< Updated upstream:htmlparser/lexer.go
	for !s.EOF() && !predicate(s.Current()) {
========
	for !s.EOF() && !predicate(s.Ch) {
>>>>>>>> Stashed changes:tools/scanner.go
		s.Take()
	}
	return s.input[start:s.position]
}

func (s *Scanner) ConsumeN(n int) string {
	start := s.position
	for i := 0; i < n && !s.EOF(); i++ {
		s.Take()
	}
	return s.input[start:s.position]
}

func (s *Scanner) Find(target rune) bool {
	for !s.EOF() {
<<<<<<<< Updated upstream:htmlparser/lexer.go
		if s.Current() == target {
========
		if s.Ch == target {
>>>>>>>> Stashed changes:tools/scanner.go
			return true
		}
		s.Take()
	}
	return false
}

func (s *Scanner) FindString(target string) bool {
	targetLen := len(target)
	for s.position+targetLen <= len(s.input) {
		if s.input[s.position:s.position+targetLen] == target {
			return true
		}
		s.Take()
	}
	return false
}

func (s *Scanner) Remaining() string {
	return s.input[s.position:]
}

func (s *Scanner) Len() int {
	return len(s.input)
}
<<<<<<<< Updated upstream:htmlparser/lexer.go
========

func (s *Scanner) Location() string {
	return fmt.Sprintf("%d:%d", s.line, s.column)
}

func (s *Scanner) SetLocation(line, column int) {
	s.line = line
	s.column = column
	s.updateCh()
}

func (s *Scanner) SkipWhitespace() {
	for !s.EOF() && (s.Ch == ' ' || s.Ch == '\t' || s.Ch == '\n' || s.Ch == '\r') {
		s.Take()
	}
}
>>>>>>>> Stashed changes:tools/scanner.go
