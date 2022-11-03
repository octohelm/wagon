package fsutil

import (
	"os"
	"testing"

	"github.com/spf13/afero"

	testingx "github.com/octohelm/x/testing"
)

func TestFileLister(t *testing.T) {
	cwd, _ := os.Getwd()
	fs := afero.NewIOFS(NewBasePathFs(NewOsFs(), cwd))

	ls := WithFilter(NewFileLister(fs), nil, []string{"*_test.go"})

	count := 0

	_ = ls.Range("", func(filename string) error {
		count++
		return nil
	})

	testingx.Expect(t, count, testingx.Be(5))
}
