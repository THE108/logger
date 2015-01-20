package handler

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type HandlerSuite struct {
	tmpDir string
}

var _ = Suite(&HandlerSuite{})

func (s *HandlerSuite) SetUpTest(c *C) {

	s.tmpDir = c.MkDir()
}

func (s *HandlerSuite) TestWatchedFileMoved(c *C) {

	content1 := []byte("TEST1\n")
	content2 := []byte("TEST2\n")

	filename := filepath.Join(s.tmpDir, "test.log")

	c.Log(filename)

	wf := NewWatchedFile(filename)

	bw, err := wf.Write(content1)

	c.Assert(bw, Equals, len(content1))
	c.Assert(err, IsNil)

	c.Assert(os.Rename(filename, filename+".1"), IsNil)

	bw, err = wf.Write(content2)

	c.Assert(bw, Equals, len(content2))
	c.Assert(err, IsNil)

	real_content1, err := ioutil.ReadFile(filename + ".1")

	c.Assert(string(real_content1), Equals, string(content1))
	c.Assert(err, IsNil)

	real_content2, err := ioutil.ReadFile(filename)

	c.Assert(string(real_content2), Equals, string(content2))
	c.Assert(err, IsNil)
}

func (s *HandlerSuite) TestWatchedFileReplaced(c *C) {

	content1 := []byte("TEST1\n")
	content2 := []byte("TEST2\n")
	replacedMarker := []byte("replaced\n")

	filename := filepath.Join(s.tmpDir, "test.log")

	c.Log(filename)

	wf := NewWatchedFile(filename)

	bw, err := wf.Write(content1)

	c.Assert(bw, Equals, len(content1))
	c.Assert(err, IsNil)

	c.Assert(os.Rename(filename, filename+".1"), IsNil)

	c.Assert(ioutil.WriteFile(filename, replacedMarker, 0777), IsNil)

	bw, err = wf.Write(content2)

	c.Assert(bw, Equals, len(content2))
	c.Assert(err, IsNil)

	real_content1, err := ioutil.ReadFile(filename + ".1")

	c.Assert(string(real_content1), Equals, string(content1))
	c.Assert(err, IsNil)

	real_content2, err := ioutil.ReadFile(filename)

	c.Assert(string(real_content2), Equals, string(append(replacedMarker, content2...)))
	c.Assert(err, IsNil)
}

func (s *HandlerSuite) BenchmarkWatchedFile(c *C) {

	filename := filepath.Join(s.tmpDir, "test.log")

	c.Log(filename)

	wf := NewWatchedFile(filename)

	content := []byte("TEST\n")

	for i := 0; i < c.N; i++ {
		bw, err := wf.Write(content)
		c.Assert(bw, Equals, len(content))
		c.Assert(err, IsNil)
	}
}
