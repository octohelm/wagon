/*
Package task GENERATED BY gengo:runtimedoc
DON'T EDIT THIS FILE
*/
package task

// nolint:deadcode,unused
func runtimeDoc(v any, names ...string) ([]string, bool) {
	if c, ok := v.(interface {
		RuntimeDoc(names ...string) ([]string, bool)
	}); ok {
		return c.RuntimeDoc(names...)
	}
	return nil, false
}

func (v Client) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Env":
			return []string{}, true
		case "Filesystem":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v ClientFile) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Path":
			return []string{}, true
		case "Contents":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Copy) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Input":
			return []string{}, true
		case "Contents":
			return []string{}, true
		case "Source":
			return []string{}, true
		case "Dest":
			return []string{}, true
		case "Include":
			return []string{}, true
		case "Exclude":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Diff) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Upper":
			return []string{}, true
		case "Lower":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Dockerfile) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Source":
			return []string{}, true
		case "Dockerfile":
			return []string{}, true
		case "Target":
			return []string{}, true
		case "BuildArg":
			return []string{}, true
		case "Auth":
			return []string{}, true
		case "Platform":
			return []string{}, true
		case "Config":
			return []string{}, true
		case "Output":
			return []string{}, true
		case "Hosts":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v DockerfilePath) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Path":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Exec) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Input":
			return []string{}, true
		case "Mounts":
			return []string{}, true
		case "Env":
			return []string{}, true
		case "Workdir":
			return []string{}, true
		case "Args":
			return []string{}, true
		case "User":
			return []string{}, true
		case "Always":
			return []string{}, true
		case "Exit":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v HTTPFetch) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Source":
			return []string{}, true
		case "Dest":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Merge) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Inputs":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Mkdir) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Input":
			return []string{}, true
		case "Path":
			return []string{}, true
		case "Permissions":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Nop) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Input":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Pull) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Source":
			return []string{}, true
		case "Platform":
			return []string{}, true
		case "Auth":
			return []string{}, true
		case "Config":
			return []string{}, true
		case "Output":
			return []string{}, true
		case "ResolveMode":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Push) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Pusher":
			return []string{}, true

		}
		if doc, ok := runtimeDoc(v.Pusher, names...); ok {
			return doc, ok
		}

		return nil, false
	}
	return []string{}, true
}

func (v PushImage) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Dest":
			return []string{}, true
		case "Type":
			return []string{}, true
		case "Input":
			return []string{}, true
		case "Config":
			return []string{}, true
		case "Platform":
			return []string{}, true
		case "Auth":
			return []string{}, true
		case "Result":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v PushManifests) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Dest":
			return []string{}, true
		case "Type":
			return []string{}, true
		case "Inputs":
			return []string{}, true
		case "Auth":
			return []string{}, true
		case "Result":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v ReadFile) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Input":
			return []string{}, true
		case "Path":
			return []string{}, true
		case "Contents":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v RegistrySetting) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Auth":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Rm) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Input":
			return []string{}, true
		case "Path":
			return []string{}, true
		case "Output":
			return []string{}, true
		case "AllowWildcard":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Set) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Input":
			return []string{}, true
		case "Config":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Setting) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Registry":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Source) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Path":
			return []string{}, true
		case "Include":
			return []string{}, true
		case "Exclude":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v Version) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}

func (v WriteFile) RuntimeDoc(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
		case "Input":
			return []string{}, true
		case "Path":
			return []string{}, true
		case "Contents":
			return []string{}, true
		case "Permissions":
			return []string{}, true
		case "Output":
			return []string{}, true

		}

		return nil, false
	}
	return []string{}, true
}
