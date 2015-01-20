package logger

import (
	"bytes"
	"testing"

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
	scope.Debug("DEBUG")
	scope.Errorf("ERRORF -> %d", 1)

	scope.SetLevel(WARNING)

	scope.Debug("DEBUG")

	scope.SetFlags(LstdFlags | Lshortfile)

	scope.Errorf("ERRORF -> %d", 2)

	c.Assert(scope.Flush(), IsNil)

	o := out.String()

	c.Log(o)

	c.Assert(o, Matches,
		`[0-9][0-9][0-9][0-9].[0-9][0-9].[0-9][0-9] [0-9][0-9]:[0-9][0-9]:[0-9][0-9] test
		| [I] INFOLN
		| [D] DEBUG
		| [E] ERRORF -> 1
		| [D] DEBUG
		| scope_test.go:33: [E] ERRORF -> 2
`)
}
