package fsutil

import (
	"os"
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestRealPath(t *testing.T) {
	cwd, _ := os.Getwd()
	rootfs := NewOsFs()
	fs := NewBasePathFs(rootfs, cwd)

	p, _ := RealPath(fs)
	testingx.Expect(t, p, testingx.Be(cwd))

	proot, _ := RealPath(rootfs)
	testingx.Expect(t, proot, testingx.Be("/"))
}
