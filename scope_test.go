package logger

import (
	"bytes"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ScopeSuite struct{}

var _ = Suite(&ScopeSuite{})

func (s *ScopeSuite) TestScopeLogger(c *C) {

	var out bytes.Buffer

	scope := NewScope(&out, "test", DEBUG)

	scope.Info("INFOLN")

	time.Sleep(100 * time.Millisecond)

	scope.Debug("DEBUG")
	scope.Errorf("ERRORF -> %d", 1)

	time.Sleep(200 * time.Millisecond)

	scope.SetLevel(WARNING)

	scope.Debug("DEBUG")

	scope.SetFlags(LstdFlags | Lshortfile)

	scope.Errorf("ERRORF -> %d", 2)

	c.Assert(scope.Flush(), IsNil)

	o := out.String()

	c.Log(o)

	c.Assert(o[20:], Equals,
		`test
		[I] +0ms	INFOLN
		[D] +100ms	DEBUG
		[E] +100ms	ERRORF -> 1
		[E] +300ms	scope_test.go:39: ERRORF -> 2
`)
}
