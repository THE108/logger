package handler

import (
	"os"
	"sync"
)

type WatchedFile struct {
	filename string
	file     *os.File
	mu       sync.Mutex
	fi       os.FileInfo
}

func NewWatchedFile(filename string) *WatchedFile {

	return &WatchedFile{
		filename: filename,
	}
}

func (l *WatchedFile) Write(p []byte) (n int, err error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		if err := l.open(); err != nil {
			return -1, err
		}
	}

	if err := l.rotateIfNeeded(); err != nil {
		return -1, err
	}

	return l.file.Write(p)
}

func (l *WatchedFile) Close() error {

	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

func (l *WatchedFile) Rotate() error {

	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotate()
}

func (l *WatchedFile) rotateIfNeeded() error {

	fi, err := os.Stat(l.filename)

	isNotExists := os.IsNotExist(err)

	if err != nil && !isNotExists {
		return err
	}

	isSameFile := os.SameFile(fi, l.fi)

	if isNotExists || !isSameFile {
		if err := l.rotate(); err != nil {
			return err
		}
	}

	return nil
}

func (l *WatchedFile) rotate() error {

	if err := l.close(); err != nil {
		return err
	}

	return l.open()
}

// close closes the file if it is open.
func (l *WatchedFile) close() error {

	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	l.fi = nil
	return err
}

func (l *WatchedFile) open() error {

	file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	l.file = file
	l.fi = fi

	return nil
}
