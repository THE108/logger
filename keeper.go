package logger

import (
	"os"
	"sync"
	"time"
)

const (
	backupTimeFormat = "2006-01-02T15-04-05.000"
	defaultMaxSize   = 100
)

type RotatedLogFile struct {
	filename string

	maxFileSize     int64
	maxFileDuration time.Duration

	maxUncompressedSize     int64
	maxUncompressedDuration time.Duration

	maxSize     int64
	maxDuration time.Duration

	fileSize  int64
	writeTime time.Time

	file *os.File
	mu   sync.Mutex
}

func (l *RotatedLogFile) Write(p []byte) (n int, err error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	writeLen := int64(len(p))
	now := time.Now()

	if l.file == nil {
		if err := l.openExistingOrNew(writeLen, &now); err != nil {
			return 0, err
		}
	}

	sizeOverflow := l.fileSize+writeLen >= l.maxFileSize
	timeOverflow := now.Sub(l.writeTime) >= l.maxFileDuration

	if sizeOverflow || timeOverflow {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = l.file.Write(p)

	l.fileSize += int64(n)
	l.writeTime = now

	return n, err
}

// Close implements io.Closer, and closes the current logfile.
func (l *RotatedLogFile) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
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

// openExistingOrNew opens the logfile if it exists and if the current write
// would not put it over MaxSize.  If there is no such file or the write would
// put it over the MaxSize, a new file is created.
func (l *RotatedLogFile) openExistingOrNew(writeLen int64, now *time.Time) error {

	info, err := os.Stat(l.filename)
	if os.IsNotExist(err) {
		return l.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	sizeOverflow := info.Size()+writeLen >= l.maxFileSize
	timeOverflow := now.Sub(info.ModTime()) >= l.maxFileDuration

	// the first file we find that matches our pattern will be the most
	// recently modified log file.
	if sizeOverflow || timeOverflow {
		file, err := os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			l.file = file
			l.fileSize = info.Size()
			l.writeTime = info.ModTime()
			return nil
		}
		// if we fail to open the old log file for some reason, just ignore
		// it and open a new log file.
	}
	return l.openNew()
}

// openNew opens a new log file for writing, moving any old log file out of the
// way.  This methods assumes the file has already been closed.
func (l *RotatedLogFile) openNew() error {

	err := os.MkdirAll(l.dir(), 0744)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := l.filename()
	mode := os.FileMode(0644)
	info, err := os.Stat(name)
	if err == nil {
		// Copy the mode off the old logfile.
		mode = info.Mode()
		// move the existing file
		newname := backupName(name, l.LocalTime)
		if err := os.Rename(name, newname); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}

		// this is a no-op anywhere but linux
		if err := chown(name, info); err != nil {
			return err
		}
	}

	// we use truncate here because this should only get called when we've moved
	// the file ourselves. if someone else creates the file in the meantime,
	// just wipe out the contents.
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	l.file = f
	l.size = 0
	return nil
}
