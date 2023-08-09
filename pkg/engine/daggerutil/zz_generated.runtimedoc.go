/*
Package daggerutil GENERATED BY gengo:runtimedoc 
DON'T EDIT THIS FILE
*/
package daggerutil

// nolint:deadcode,unused
func runtimeDoc(v any, names ...string) ([]string, bool) {
	if c, ok := v.(interface {
		RuntimeDoc(names ...string) ([]string, bool)
	}); ok {
		return c.RuntimeDoc(names...)
	}
	return nil, false
}

func (v Container) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "ID":
			return []string{}, true
		case "Platform":
			return []string{}, true
		case "Entrypoint":
			return []string{}, true
		case "DefaultArgs":
			return []string{}, true
		case "Workdir":
			return []string{}, true
		case "User":
			return []string{}, true
		case "EnvVariables":
			return []string{}, true
		case "Labels":
			return []string{}, true
		case "RootFS":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (DirectConn) RuntimeDoc(names ...string) ([]string, bool) {
	return []string{}, true
}
func (v Directory) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "ID":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v EnvVariable) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Name":
			return []string{}, true
		case "Value":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Label) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Name":
			return []string{}, true
		case "Value":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}
