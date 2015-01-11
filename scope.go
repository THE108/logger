package logger

import (
	"fmt"
)

type Scope struct {
	out    func(v ...interface{})
	buffer []byte
}

func NewScope(out func(v ...interface{})) *Scope {

	return &Scope{
		out:    out,
		buffer: []byte{'\n'},
	}
}

func (s *Scope) String() string {
	return string(s.buffer)
}

func (s *Scope) Flush() {
	s.out(s.String())
}

func (s *Scope) Output(line string) {
	s.buffer = append(s.buffer, '\t', '\t', '|', ' ')
	s.buffer = append(s.buffer, line...)
	s.buffer = append(s.buffer, '\n')
}

func (s *Scope) Outputf(format string, v ...interface{}) {
	s.Output(fmt.Sprintf(format, v...))
}
