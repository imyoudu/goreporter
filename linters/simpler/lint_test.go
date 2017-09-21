package simpler

import (
	"testing"

	"github.com/imyoudu/goreporter/linters/simpler/lint/testutil"
)

func TestAll(t *testing.T) {
	testutil.TestAll(t, NewChecker(), "")
}
