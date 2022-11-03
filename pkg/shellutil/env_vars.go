package shellutil

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/octohelm/x/encoding"
	"mvdan.cc/sh/v3/expand"
)

func FromEnviron(environ []string) EnvVars {
	envVars := EnvVars{}

	for i := range environ {
		parts := strings.SplitN(environ[i], "=", 2)
		name := parts[0]
		value := ""
		if len(parts) > 1 {
			value = parts[1]
		}
		envVars[name] = value
	}

	return envVars
}

type EnvVars map[string]any

func (envVars EnvVars) Merge(envVars2 EnvVars) EnvVars {
	final := EnvVars{}

	for k := range envVars {
		final[k] = envVars[k]
	}

	for k := range envVars2 {
		final[k] = envVars2[k]
	}

	return final
}

func (envVars EnvVars) Each(f func(name string, vr expand.Variable) bool) {
	for name := range envVars {
		f(name, envVars.Get(name))
	}
}

func (envVars EnvVars) Get(name string) expand.Variable {
	if v, ok := envVars[name]; ok {
		str, err := encoding.MarshalText(v)
		if err == nil {
			return expand.Variable{
				Local: true,
				Kind:  expand.String,
				Str:   string(str),
			}
		}
	}
	return expand.Variable{
		Local: true,
		Kind:  expand.Unset,
	}
}

func (envVars EnvVars) String() string {
	buf := bytes.NewBuffer(nil)

	envKeys := make([]string, 0, len(envVars))
	for env := range envVars {
		envKeys = append(envKeys, env)
	}
	sort.Strings(envKeys)

	for _, envKey := range envKeys {
		buf.WriteString(envKey)
		buf.WriteString("=")
		v, _ := encoding.MarshalText(envVars[envKey])
		buf.WriteString(strconv.Quote(string(v)))
		buf.WriteString(" ")
	}

	return buf.String()
}
