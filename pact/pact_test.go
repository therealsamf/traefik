// This is the main file that sets up pact (https://pact.io) tests using go-check.
package pact

import (
	"github.com/go-check/check"
	"testing"
)

func Test(t *testing.T) {
	check.Suite(&ForwardAuthPactSuite{})

	check.TestingT(t)
}
