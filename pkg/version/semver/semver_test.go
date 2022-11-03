package semver

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestParse(t *testing.T) {
	sv := Parse("v0.1.4-0.20221213104148-f2c1d8adfe96")
	testingx.Expect(t, sv.Major, testingx.Be(0))
	testingx.Expect(t, sv.Minor, testingx.Be(1))
	testingx.Expect(t, sv.Patch, testingx.Be(3))
}
