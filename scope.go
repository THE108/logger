package logger

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
)

type Scope struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix to write at beginning of each line
	flag   int        // properties
	level  int        // verbosity level
	out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write
}

func NewScope(out io.Writer, prefix string, level int) *Scope {

	return &Scope{
		prefix: prefix,
		flag:   LstdFlags,
		level:  level,
		out:    out,
	}
}

func (s *Scope) Flush() error {

	s.mu.Lock()
	defer s.mu.Unlock()

	var b []byte

	if s.flag&(Ldate|Ltime|Lmicroseconds) != 0 {

		t := time.Now() // get this early

		if s.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(&b, year, 4)
			b = append(b, '.')
			itoa(&b, int(month), 2)
			b = append(b, '.')
			itoa(&b, day, 2)
			b = append(b, ' ')
		}

		if s.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(&b, hour, 2)
			b = append(b, ':')
			itoa(&b, min, 2)
			b = append(b, ':')
			itoa(&b, sec, 2)
			if s.flag&Lmicroseconds != 0 {
				b = append(b, '.')
				itoa(&b, t.Nanosecond()/1e3, 6)
			}
			b = append(b, ' ')
		}
	}

	if len(s.prefix) > 0 {
		b = append(b, s.prefix...)
		b = append(b, ' ')
	}

	b[len(b)-1] = '\n'
	b = append(b, s.buf...)

	_, err := s.out.Write(b)
	return err
}

func (s *Scope) Output(calldepth, level int, message string) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.level > level {
		return
	}

	s.buf = append(s.buf, '\t', '\t', '|', ' ')

	if s.flag&(Lshortfile|Llongfile) != 0 {

		// var file string
		// var line int
		// var ok bool

		// release lock while getting caller info - it's expensive.
		s.mu.Unlock()
		_, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		s.mu.Lock()

		if s.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		s.buf = append(s.buf, file...)
		s.buf = append(s.buf, ':')
		itoa(&s.buf, line, -1)
		s.buf = append(s.buf, ": "...)
	}

	s.buf = append(s.buf, '[', levelChar[level], ']', ' ')
	s.buf = append(s.buf, message...)
	s.buf = append(s.buf, '\n')
}

// Debug calls s.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (s *Scope) Debug(v ...interface{}) {
	s.Output(2, DEBUG, fmt.Sprint(v...))
}

// Debugf calls s.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (s *Scope) Debugf(format string, v ...interface{}) {
	s.Output(2, DEBUG, fmt.Sprintf(format, v...))
}

// Info calls s.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (s *Scope) Info(v ...interface{}) {
	s.Output(2, INFO, fmt.Sprint(v...))
}

// Infof calls s.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (s *Scope) Infof(format string, v ...interface{}) {
	s.Output(2, INFO, fmt.Sprintf(format, v...))
}

// Warn calls s.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (s *Scope) Warn(v ...interface{}) {
	s.Output(2, WARNING, fmt.Sprint(v...))
}

// Warnf calls s.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (s *Scope) Warnf(format string, v ...interface{}) {
	s.Output(2, WARNING, fmt.Sprintf(format, v...))
}

// Error calls s.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (s *Scope) Error(v ...interface{}) {
	s.Output(2, ERROR, fmt.Sprint(v...))
}

// Errorf calls s.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (s *Scope) Errorf(format string, v ...interface{}) {
	s.Output(2, ERROR, fmt.Sprintf(format, v...))
}

// Flags returns the output flags for the logger.
func (s *Scope) Flags() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flag
}

// SetFlags sets the output flags for the logger.
func (s *Scope) SetFlags(flag int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flag = flag
}

// Prefix returns the output prefix for the logger.
func (s *Scope) Prefix() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.prefix
}

// SetPrefix sets the output prefix for the logger.
func (s *Scope) SetPrefix(prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefix = prefix
}

func (s *Scope) Level() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.level
}

func (s *Scope) SetLevel(level int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.level = level
}
