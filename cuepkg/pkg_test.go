package cuepkg

import (
	"fmt"
	"github.com/octohelm/wagon/pkg/fsutil"
	"testing"
)

func TestCuePkg(t *testing.T) {
	cuepkgs, err := createWagonModule(daggerPortalModules)
	if err != nil {
		t.Fatal(err)
	}

	_ = fsutil.RangeFile(cuepkgs, "", func(filename string) error {
		fmt.Println(filename)
		return nil
	})
}
