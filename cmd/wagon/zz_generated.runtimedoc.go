/*
Package main GENERATED BY gengo:runtimedoc
DON'T EDIT THIS FILE
*/
package main

// nolint:deadcode,unused
func runtimeDoc(v any, names ...string) ([]string, bool) {
	if c, ok := v.(interface {
		RuntimeDoc(names ...string) ([]string, bool)
	}); ok {
		return c.RuntimeDoc(names...)
	}
	return nil, false
}

func (v Do) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Logger":
			return []string{}, true
		case "Pipeline":
			return []string{}, true

		}
		if doc, ok := runtimeDoc(v.Logger, names...); ok {
			return doc, ok
		}
		if doc, ok := runtimeDoc(v.Pipeline, names...); ok {
			return doc, ok
		}

		return nil, false
	}
	return []string{}, true
}

func (v Get) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Logger":
			return []string{}, true
		case "GetMod":
			return []string{}, true

		}
		if doc, ok := runtimeDoc(v.Logger, names...); ok {
			return doc, ok
		}
		if doc, ok := runtimeDoc(v.GetMod, names...); ok {
			return doc, ok
		}

		return nil, false
	}
	return []string{}, true
}

func (v GetMod) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Pkgs":
			return []string{}, true
		case "Update":
			return []string{
				"Update to latest",
			}, true
		case "Import":
			return []string{
				"declare language for generate. support values: go",
			}, true

		}

		return nil, false
	}
	return []string{}, true
}
