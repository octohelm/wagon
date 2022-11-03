package spec

import (
	"archive/tar"
	"fmt"
)

type Ship struct {
	Package      Package           `json:"pkg"`
	Files        map[string]string `json:"files"`
	Scripts      map[string]string `json:"scripts"`
	Dependencies map[string]string `json:"dependencies"`
}

type Pkg struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Revision    string `json:"revision,omitempty"`
	Description string `json:"description,omitempty"`
}

func (pkg Pkg) String() string {
	return fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
}

type Package struct {
	Pkg
	Platform
	Files []tar.Header `json:"files,omitempty"`
}
