package logger

import (
	"os"
	"sync"
)

type RotatedLogFile struct {
	filename string
	file     *os.File
	mu       sync.Mutex
}

func NewRotatedLogFile(filename string) *RotatedLogFile {

	return &RotatedLogFile{
		filename: filename,
	}
}

func (l *RotatedLogFile) Write(p []byte) (n int, err error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		if err := l.open(); err != nil {
			return -1, err
		}
	}

	return l.file.Write(p)
}

func (l *RotatedLogFile) Close() error {

	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

func (l *RotatedLogFile) Rotate() error {

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.close(); err != nil {
		return err
	}

	return l.open()
}

// close closes the file if it is open.
func (l *RotatedLogFile) close() error {

	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

func (l *RotatedLogFile) open() error {

	file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.file = file
	return nil
}
