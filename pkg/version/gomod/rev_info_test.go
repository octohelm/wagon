package gomod

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	testingx "github.com/octohelm/x/testing"
	"golang.org/x/mod/module"
)

func TestRevInfo(t *testing.T) {
	vt, _ := module.PseudoVersionTime("v0.1.4-0.20221213104148-f2c1d8adfe96")

	testingx.Expect(t, (RevInfo{
		Name:  "v0.1.3",
		Short: "f2c1d8adfe96",
		Time:  vt,
	}).Version(), testingx.Equal("v0.1.3-20221213104148-f2c1d8adfe96"))

	testingx.Expect(t, (RevInfo{
		Name:   "v0.1.3",
		Offset: 8,
		Short:  "f2c1d8adfe96",
		Time:   vt,
	}).Version(), testingx.Equal("v0.1.4-0.20221213104148-f2c1d8adfe96"))
}

func TestLocalRevInfo(t *testing.T) {
	revInfo, _ := LocalRevInfo(".")
	spew.Dump(revInfo.Version())
}
