package fsutil

import (
	"fmt"
	"testing"

	testingx "github.com/octohelm/x/testing"
	"github.com/spf13/afero"
)

func TestIOFS(t *testing.T) {
	mfs := afero.NewMemMapFs()

	for i := 0; i < 10; i++ {
		_ = afero.WriteFile(mfs, fmt.Sprintf("%d.txt", i), []byte(fmt.Sprintf("%d.txt", i)), 0644)
	}

	miofs := afero.NewIOFS(mfs)

	t.Run("#RangeFile", func(t *testing.T) {
		files := make([]string, 0)

		_ = RangeFile(miofs, "", func(filename string) error {
			files = append(files, filename)
			return nil
		})

		testingx.Expect(t, len(files), testingx.Be(10))
	})
}
