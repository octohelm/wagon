package semver

import (
	"fmt"
	"strconv"

	"golang.org/x/mod/module"
)

type SemVer struct {
	Major int
	Minor int
	Patch int
}

func (v SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func Parse(v string) *SemVer {
	sv := &SemVer{}

	base, err := module.PseudoVersionBase(v)
	if err == nil {
		v = base
	}

	if len(v) > 0 && v[0] == 'v' {
		v = v[1:]
	}

	var ok bool
	sv.Major, v, ok = parseInt(v)
	if ok {
		if len(v) > 0 && v[0] == '.' {
			v = v[1:]
		}
		sv.Minor, v, ok = parseInt(v)

		if ok {
			if len(v) > 0 && v[0] == '.' {
				v = v[1:]
			}
			sv.Patch, _, _ = parseInt(v)
		}
	}

	return sv
}

func parseInt(v string) (t int, rest string, ok bool) {
	if v == "" {
		return
	}
	if v[0] < '0' || '9' < v[0] {
		return
	}
	i := 1
	for i < len(v) && '0' <= v[i] && v[i] <= '9' {
		i++
	}
	if v[0] == '0' && i != 1 {
		return
	}
	t, _ = strconv.Atoi(v[:i])
	return t, v[i:], true
}
