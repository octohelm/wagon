/*
Package strfmt GENERATED BY gengo:runtimedoc
DON'T EDIT THIS FILE
*/
package strfmt

// nolint:deadcode,unused
func runtimeDoc(v any, names ...string) ([]string, bool) {
	if c, ok := v.(interface {
		RuntimeDoc(names ...string) ([]string, bool)
	}); ok {
		return c.RuntimeDoc(names...)
	}
	return nil, false
}

func (v URL) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "URL":
			return []string{}, true

		}
		if doc, ok := runtimeDoc(v.URL, names...); ok {
			return doc, ok
		}

		return nil, false
	}
	return []string{}, true
}
