/*
Package plan GENERATED BY gengo:runtimedoc 
DON'T EDIT THIS FILE
*/
package plan

// nolint:deadcode,unused
func runtimeDoc(v any, names ...string) ([]string, bool) {
	if c, ok := v.(interface {
		RuntimeDoc(names ...string) ([]string, bool)
	}); ok {
		return c.RuntimeDoc(names...)
	}
	return nil, false
}

func (v Auth) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Address":
			return []string{}, true
		case "Username":
			return []string{}, true
		case "SecretID":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (WorkdirType) RuntimeDoc(names ...string) ([]string, bool) {
	return []string{}, true
}
